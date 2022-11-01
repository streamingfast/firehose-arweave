package cli

import (
	"fmt"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/streamingfast/bstream/blockstream"
	"github.com/streamingfast/dlauncher/launcher"
	"github.com/streamingfast/firehose-arweave/nodemanager"
	"github.com/streamingfast/logging"
	nodeManager "github.com/streamingfast/node-manager"
	nodeManagerApp "github.com/streamingfast/node-manager/app/node_manager2"
	"github.com/streamingfast/node-manager/metrics"
	"github.com/streamingfast/node-manager/operator"
	pbbstream "github.com/streamingfast/pbgo/sf/bstream/v1"
	pbheadinfo "github.com/streamingfast/pbgo/sf/headinfo/v1"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

var nodeLogger, nodeTracer = logging.PackageLogger("node", "github.com/streamingfast/firehose-arweave/node")
var nodeArweaveLogger, nodeArweaveTracer = logging.PackageLogger("node.arweave", "github.com/streamingfast/firehose-arweave/node/dummy-chain", DefaultLevelInfo)

var readerLogger, readerTracer = logging.PackageLogger("reader", "github.com/streamingfast/firehose-arweave/reader")
var readerArweaveLogger, readerArweaveTracer = logging.PackageLogger("reader.arweave", "github.com/streamingfast/firehose-arweave/reader/dummy-chain", DefaultLevelInfo)

func registerCommonNodeFlags(cmd *cobra.Command, flagPrefix string, managerAPIAddr string) {
	cmd.Flags().String(flagPrefix+"path", "thegarii", FlagDescription(`
		Process that will be invoked reader component, can be a full path or just the binary's name, in which case the binary is
		searched for paths listed by the PATH environment variable (following operating system rules around PATH handling).
	`))
	cmd.Flags().String(flagPrefix+"data-dir", "{data-dir}/{node-role}/data", "Directory for node data ({node-role} is either reader, peering or dev-miner)")
	cmd.Flags().Bool(flagPrefix+"debug-firehose-logs", false, "[DEV] Prints Firehose instrumentation logs to standard output, should be use for debugging purposes only")
	cmd.Flags().Bool(flagPrefix+"log-to-zap", true, FlagDescription(`
		When sets to 'true', all standard error output emitted by the invoked process defined via '%s'
		is intercepted, split line by line and each line is then transformed and logged through the Firehose stack
		logging system. The transformation extracts the level and remove the timestamps creating a 'sanitized' version
		of the logs emitted by the blockchain's managed client process. If this is not desirable, disabled the flag
		and all the invoked process standard error will be redirect to 'firearweave' standard's output.
	`, flagPrefix+"path"))
	cmd.Flags().String(flagPrefix+"manager-api-addr", managerAPIAddr, "Arweave node manager API address")
	cmd.Flags().Duration(flagPrefix+"readiness-max-latency", 10*time.Minute, "Determine the maximum head block latency at which the instance will be determined healthy. Some chains have more regular block production than others.")
	cmd.Flags().String(flagPrefix+"arguments", "", "If not empty, overrides the list of default node arguments (computed from node type and role). Start with '+' to append to default args instead of replacing. ")
}

func registerNode(kind string, extraFlagRegistration func(cmd *cobra.Command) error, managerAPIaddr string) {
	if kind != "reader" {
		panic(fmt.Errorf("invalid kind value, must be either 'reader', got %q", kind))
	}

	app := fmt.Sprintf("%s-node", kind)
	flagPrefix := fmt.Sprintf("%s-", app)

	launcher.RegisterApp(rootLog, &launcher.AppDef{
		ID:          app,
		Title:       fmt.Sprintf("Arweave Node (%s)", kind),
		Description: fmt.Sprintf("Arweave %s node with built-in operational manager", kind),
		RegisterFlags: func(cmd *cobra.Command) error {
			registerCommonNodeFlags(cmd, flagPrefix, managerAPIaddr)
			extraFlagRegistration(cmd)
			return nil
		},
		InitFunc: func(runtime *launcher.Runtime) error {
			return nil
		},
		FactoryFunc: nodeFactoryFunc(flagPrefix, kind),
	})
}

func nodeFactoryFunc(flagPrefix, kind string) func(*launcher.Runtime) (launcher.App, error) {
	return func(runtime *launcher.Runtime) (launcher.App, error) {
		var appLogger *zap.Logger
		var appTracer logging.Tracer
		var supervisedProcessLogger *zap.Logger

		switch kind {
		case "node":
			appLogger = nodeLogger
			appTracer = nodeTracer
			supervisedProcessLogger = nodeArweaveLogger
		case "reader":
			appLogger = readerLogger
			appTracer = readerTracer

			supervisedProcessLogger = readerArweaveLogger
		default:
			panic(fmt.Errorf("unknown node kind %q", kind))
		}

		sfDataDir := runtime.AbsDataDir

		nodePath := viper.GetString(flagPrefix + "path")
		nodeDataDir := replaceNodeRole(kind, mustReplaceDataDir(sfDataDir, viper.GetString(flagPrefix+"data-dir")))

		readinessMaxLatency := viper.GetDuration(flagPrefix + "readiness-max-latency")
		debugFirehoseLogs := viper.GetBool(flagPrefix + "debug-firehose-logs")
		logToZap := viper.GetBool(flagPrefix + "log-to-zap")
		shutdownDelay := viper.GetDuration("common-system-shutdown-signal-delay") // we reuse this global value
		httpAddr := viper.GetString(flagPrefix + "manager-api-addr")
		batchStartBlockNum := viper.GetUint64("reader-node-start-block-num")
		batchStopBlockNum := viper.GetUint64("reader-node-stop-block-num")
		endpoints := viper.GetStringSlice("reader-node-endpoints")

		arguments := viper.GetString(flagPrefix + "arguments")
		nodeArguments, err := buildNodeArguments(
			nodeDataDir,
			kind,
			endpoints,
			batchStartBlockNum,
			batchStopBlockNum,
			arguments,
		)
		if err != nil {
			return nil, fmt.Errorf("cannot build node bootstrap arguments")
		}
		metricsAndReadinessManager := buildMetricsAndReadinessManager(strings.TrimSuffix(flagPrefix, "-"), readinessMaxLatency)

		superviser := nodemanager.NewSuperviser(
			nodePath,
			nodeArguments,
			nodeDataDir,
			metricsAndReadinessManager.UpdateHeadBlock,
			debugFirehoseLogs,
			logToZap,
			appLogger,
			supervisedProcessLogger,
		)

		bootstrapper := &bootstrapper{
			nodeDataDir: nodeDataDir,
		}

		chainOperator, err := operator.New(
			appLogger,
			superviser,
			metricsAndReadinessManager,
			&operator.Options{
				ShutdownDelay:              shutdownDelay,
				EnableSupervisorMonitoring: true,
				Bootstrapper:               bootstrapper,
			})
		if err != nil {
			return nil, fmt.Errorf("unable to create chain operator: %w", err)
		}

		if kind != "reader" {
			return nodeManagerApp.New(&nodeManagerApp.Config{
				HTTPAddr: httpAddr,
			}, &nodeManagerApp.Modules{
				Operator:                   chainOperator,
				MetricsAndReadinessManager: metricsAndReadinessManager,
			}, appLogger), nil
		}

		blockStreamServer := blockstream.NewUnmanagedServer(blockstream.ServerOptionWithLogger(appLogger))

		mergedBlocksStoreURL, oneBlockStoreURL, err := getCommonStoresURLs(runtime.AbsDataDir)
		workingDir := mustReplaceDataDir(sfDataDir, viper.GetString("reader-node-working-dir"))
		gprcListenAdrr := viper.GetString("reader-node-grpc-listen-addr")
		waitTimeForUploadOnShutdown := viper.GetDuration("reader-node-wait-upload-complete-on-shutdown")
		oneBlockFileSuffix := viper.GetString("reader-node-one-block-suffix")
		blocksChanCapacity := viper.GetInt("reader-node-blocks-chan-capacity")

		readerPlugin, err := getReaderLogPlugin(
			blockStreamServer,
			oneBlockStoreURL,
			mergedBlocksStoreURL,
			workingDir,
			batchStartBlockNum,
			batchStopBlockNum,
			blocksChanCapacity,
			waitTimeForUploadOnShutdown,
			oneBlockFileSuffix,
			chainOperator.Shutdown,
			metricsAndReadinessManager,
			appLogger,
			appTracer,
		)
		if err != nil {
			return nil, fmt.Errorf("new reader plugin: %w", err)
		}

		superviser.RegisterLogPlugin(readerPlugin)

		return nodeManagerApp.New(&nodeManagerApp.Config{
			HTTPAddr: httpAddr,
			GRPCAddr: gprcListenAdrr,
		}, &nodeManagerApp.Modules{
			Operator:                   chainOperator,
			MindreaderPlugin:           readerPlugin,
			MetricsAndReadinessManager: metricsAndReadinessManager,
			RegisterGRPCService: func(server grpc.ServiceRegistrar) error {
				pbheadinfo.RegisterHeadInfoServer(server, blockStreamServer)
				pbbstream.RegisterBlockStreamServer(server, blockStreamServer)

				return nil
			},
		}, appLogger), nil
	}
}

type bootstrapper struct {
	nodeDataDir string
}

func (b *bootstrapper) Bootstrap() error {
	// You can copy coniguration files here into your working data dir to run the node off of
	return nil
}

type nodeArgsByRole map[string]string

func buildNodeArguments(nodeDataDir, nodeRole string, endpoints []string, start, stop uint64, args string) ([]string, error) {
	thegariiArgs := []string{"-d", "-B", "20", "-c", "20", "console", "-f", "--data-directory", filepath.Join(nodeDataDir, "thegarii")}
	if len(endpoints) > 0 {
		setEndpoints := append([]string{"-e"}, endpoints...)
		thegariiArgs = append(setEndpoints, thegariiArgs...)
	} else {
		endpoints = []string{"-e", "https://arweave.net"}
		thegariiArgs = append(endpoints, thegariiArgs...)
	}

	if start != 0 {
		thegariiArgs = append(thegariiArgs, []string{"-s", strconv.FormatUint(start, 10)}...)
	}

	if stop != 0 {
		thegariiArgs = append(thegariiArgs, []string{"-e", strconv.FormatUint(stop, 10)}...)
	}

	typeRoles := nodeArgsByRole{
		"reader": strings.Join(thegariiArgs, " "),
	}

	argsString, ok := typeRoles[nodeRole]
	if !ok {
		return nil, fmt.Errorf("invalid node role: %s", nodeRole)
	}

	if strings.HasPrefix(args, "+") {
		argsString = strings.Replace(argsString, "{extra-arg}", args[1:], -1)
	} else if args == "" {
		argsString = strings.Replace(argsString, "{extra-arg}", "", -1)
	} else {
		argsString = args
	}

	argsString = strings.Replace(argsString, "{node-data-dir}", nodeDataDir, -1)
	argsSlice := strings.Fields(argsString)
	return argsSlice, nil
}

func buildMetricsAndReadinessManager(name string, maxLatency time.Duration) *nodeManager.MetricsAndReadinessManager {
	headBlockTimeDrift := metrics.NewHeadBlockTimeDrift(name)
	headBlockNumber := metrics.NewHeadBlockNumber(name)
	appReadiness := metrics.NewAppReadiness(name)

	metricsAndReadinessManager := nodeManager.NewMetricsAndReadinessManager(
		headBlockTimeDrift,
		headBlockNumber,
		appReadiness,
		maxLatency,
	)
	return metricsAndReadinessManager
}

func replaceNodeRole(nodeRole, in string) string {
	return strings.Replace(in, "{node-role}", nodeRole, -1)
}
