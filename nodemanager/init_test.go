package nodemanager

import "github.com/streamingfast/logging"

var zlog, _ = logging.PackageLogger("nodemanager", "github.com/ChainSafe/firehose-arweave/nodemanager_tests")

func init() {
	logging.InstantiateLoggers()
}
