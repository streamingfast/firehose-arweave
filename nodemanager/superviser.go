package nodemanager

import (
	"fmt"
	"net"
	"strconv"
	"strings"
	"sync"

	"github.com/ShinyTrinkets/overseer"
	nodeManager "github.com/streamingfast/node-manager"
	logplugin "github.com/streamingfast/node-manager/log_plugin"
	"github.com/streamingfast/node-manager/metrics"
	"github.com/streamingfast/node-manager/superviser"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Superviser struct {
	*superviser.Superviser

	//backupMutex         sync.Mutex
	infoMutex           sync.Mutex
	binary              string
	arguments           []string
	dataDir             string
	lastBlockSeen       uint64
	serverId            string
	headBlockUpdateFunc nodeManager.HeadBlockUpdater
	Logger              *zap.Logger
}

func (s *Superviser) GetName() string {
	return "arweave"
}

func NewSuperviser(
	binary string,
	arguments []string,
	dataDir string,
	headBlockUpdateFunc nodeManager.HeadBlockUpdater,
	debugFirehoseLogs bool,
	logToZap bool,
	appLogger *zap.Logger,
	nodelogger *zap.Logger,
) *Superviser {
	// Ensure process manager line buffer is large enough (50 MiB) for our Firehose instrumentation outputting lot's of text.
	overseer.DEFAULT_LINE_BUFFER_SIZE = 50 * 1024 * 1024

	supervisor := &Superviser{
		Superviser:          superviser.New(appLogger, binary, arguments),
		Logger:              appLogger,
		binary:              binary,
		arguments:           arguments,
		dataDir:             dataDir,
		headBlockUpdateFunc: headBlockUpdateFunc,
	}

	supervisor.RegisterLogPlugin(logplugin.LogPluginFunc(supervisor.lastBlockSeenLogPlugin))

	if logToZap {
		supervisor.RegisterLogPlugin(newToZapLogPlugin(debugFirehoseLogs, nodelogger))
	} else {
		supervisor.RegisterLogPlugin(logplugin.NewToConsoleLogPlugin(debugFirehoseLogs))
	}

	appLogger.Info("created arweave superviser", zap.Object("superviser", supervisor))
	return supervisor
}

func (s *Superviser) setServerId(serverId string) error {
	ipAddr := getIPAddress()
	if ipAddr == "" {
		return fmt.Errorf("cannot find local IP address")
	}

	s.infoMutex.Lock()
	defer s.infoMutex.Unlock()
	s.serverId = fmt.Sprintf(`${1}@%s:30303`, ipAddr)
	return nil
}

func (s *Superviser) GetCommand() string {
	return s.binary + " " + strings.Join(s.arguments, " ")
}

func (s *Superviser) IsRunning() bool {
	isRunning := s.Superviser.IsRunning()
	isRunningMetricsValue := float64(0)
	if isRunning {
		isRunningMetricsValue = float64(1)
	}

	metrics.NodeosCurrentStatus.SetFloat64(isRunningMetricsValue)

	return isRunning
}

func (s *Superviser) LastSeenBlockNum() uint64 {
	return s.lastBlockSeen
}

func (s *Superviser) ServerID() (string, error) {
	return s.serverId, nil
}

func (s *Superviser) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	enc.AddString("binary", s.binary)
	enc.AddArray("arguments", stringArray(s.arguments))
	enc.AddString("data_dir", s.dataDir)
	enc.AddUint64("last_block_seen", s.lastBlockSeen)
	enc.AddString("server_id", s.serverId)

	return nil
}

func (s *Superviser) lastBlockSeenLogPlugin(line string) {
	// FIRE BLOCK <HEIGHT> ...
	if !strings.HasPrefix(line, "FIRE BLOCK") {
		return
	}

	blockNumStr := line[11:]
	nextSpace := strings.Index(blockNumStr, " ")
	if nextSpace < 0 {
		s.Logger.Error("unable to extract last block num, missing space", zap.String("line", line))
		return
	}

	blockNumStr = blockNumStr[0:nextSpace]

	blockNum, err := strconv.ParseUint(blockNumStr, 10, 64)
	if err != nil {
		s.Logger.Error("unable to extract last block num",
			zap.String("line", line),
			zap.String("block_num_str", blockNumStr),
			zap.Error(err),
		)
		return
	}

	// FIXME: Instrumentation needs to always have a way to easily decode height,
	// hash and timestamp. Right now in Arweave, we have only the height.
	//
	// It's not important for now because only readers are running and those
	// updates are carried on directly by the console reader.
	// s.headBlockUpdateFunc(s.lastBlockSeen,

	s.lastBlockSeen = blockNum
}

func getIPAddress() string {
	ifaces, err := net.Interfaces()
	if err != nil {
		return ""
	}
	for _, i := range ifaces {
		addrs, err := i.Addrs()
		if err != nil {
			continue
		}
		for _, addr := range addrs {
			var ip net.IP
			switch v := addr.(type) {
			case *net.IPNet:
				ip = v.IP
			case *net.IPAddr:
				ip = v.IP
			}
			if ip.IsGlobalUnicast() {
				return ip.String()
			}
		}
	}
	return ""
}
