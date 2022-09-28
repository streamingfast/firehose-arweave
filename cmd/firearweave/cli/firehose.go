package cli

import (
	"fmt"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	dauthAuthenticator "github.com/streamingfast/dauth/authenticator"
	"github.com/streamingfast/dlauncher/launcher"
	"github.com/streamingfast/dmetering"
	"github.com/streamingfast/dmetrics"
	firehoseApp "github.com/streamingfast/firehose/app/firehose"
	"github.com/streamingfast/logging"
)

var metricset = dmetrics.NewSet()
var headBlockNumMetric = metricset.NewHeadBlockNumber("firehose")
var headTimeDriftmetric = metricset.NewHeadTimeDrift("firehose")

func init() {
	appLogger, _ := logging.PackageLogger("firehose", "github.com/streamingfast/firehose-arweave/firehose")

	launcher.RegisterApp(rootLog, &launcher.AppDef{
		ID:          "firehose",
		Title:       "Block Firehose",
		Description: "Provides on-demand filtered blocks, depends on mergd blocks and live source",
		RegisterFlags: func(cmd *cobra.Command) error {
			cmd.Flags().String("firehose-grpc-listen-addr", FirehoseGRPCServingAddr, "Address on which the firehose will listen, appending * to the end of the listen address will start the server over an insecure TLS connection. By default Firehose will start in plain-text mode.")
			return nil
		},

		FactoryFunc: func(runtime *launcher.Runtime) (launcher.App, error) {

			// FIXME: That should be a shared dependencies across `Ethereum on StreamingFast`
			authenticator, err := dauthAuthenticator.New(viper.GetString("common-auth-plugin"))
			if err != nil {
				return nil, fmt.Errorf("unable to initialize dauth: %w", err)
			}

			// FIXME: That should be a shared dependencies across `Ethereum on StreamingFast`, it will avoid the need to call `dmetering.SetDefaultMeter`
			metering, err := dmetering.New(viper.GetString("common-metering-plugin"))
			if err != nil {
				return nil, fmt.Errorf("unable to initialize dmetering: %w", err)
			}
			dmetering.SetDefaultMeter(metering)

			var possibleIndexSizes []uint64
			for _, size := range viper.GetIntSlice("firehose-block-index-sizes") {
				if size < 0 {
					return nil, fmt.Errorf("invalid negative size for firehose-block-index-sizes: %d", size)
				}
				possibleIndexSizes = append(possibleIndexSizes, uint64(size))
			}

			sfDataDir := runtime.AbsDataDir
			return firehoseApp.New(appLogger, &firehoseApp.Config{
				MergedBlocksStoreURL:    MustReplaceDataDir(sfDataDir, viper.GetString("common-merged-blocks-store-url")),
				OneBlocksStoreURL:       MustReplaceDataDir(sfDataDir, viper.GetString("common-one-blocks-store-url")),
				ForkedBlocksStoreURL:    MustReplaceDataDir(sfDataDir, viper.GetString("common-forked-blocks-store-url")),
				BlockStreamAddr:         viper.GetString("common-live-source-addr"),
				GRPCListenAddr:          viper.GetString("firehose-grpc-listen-addr"),
				GRPCShutdownGracePeriod: time.Second,
			}, &firehoseApp.Modules{
				Authenticator:         authenticator,
				HeadTimeDriftMetric:   headTimeDriftmetric,
				HeadBlockNumberMetric: headBlockNumMetric,
			}), nil

		},
	})
}
