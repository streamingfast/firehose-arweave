// Copyright 2021 dfuse Platform Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cli

import (
	"context"
	"time"

	"github.com/ChainSafe/firehose-arweave/nodemanager/codec"
	pbarweave "github.com/ChainSafe/firehose-arweave/types/pb/sf/arweave/type/v1"
	"github.com/spf13/cobra"
	"github.com/streamingfast/bstream"
	"github.com/streamingfast/bstream/blockstream"
	"github.com/streamingfast/logging"
	nodeManager "github.com/streamingfast/node-manager"
	"github.com/streamingfast/node-manager/mindreader"
	"go.uber.org/zap"
)

func init() {
	registerNode("mindreader", registerMindreaderNodeFlags, MindreaderNodeManagerAPIAddr)
}

func registerMindreaderNodeFlags(cmd *cobra.Command) error {
	cmd.Flags().String("mindreader-node-grpc-listen-addr", MindreaderGRPCAddr, "gRPC listening address to use for serving real-time blocks")
	cmd.Flags().Bool("mindreader-node-discard-after-stop-num", false, "Ignore remaining blocks being processed after stop num (only useful if we discard the mindreader data after reprocessing a chunk of blocks)")
	cmd.Flags().String("mindreader-node-working-dir", "{data-dir}/mindreader/work", "Path where mindreader will stores its files")
	cmd.Flags().Uint("mindreader-node-start-block-num", 0, "Blocks that were produced with smaller block number then the given block num are skipped")
	cmd.Flags().Uint("mindreader-node-stop-block-num", 0, "Shutdown mindreader when we the following 'stop-block-num' has been reached, inclusively.")
	cmd.Flags().Int("mindreader-node-blocks-chan-capacity", 100, "Capacity of the channel holding blocks read by the mindreader. Process will shutdown superviser/geth if the channel gets over 90% of that capacity to prevent horrible consequences. Raise this number when processing tiny blocks very quickly")
	cmd.Flags().String("mindreader-node-one-block-suffix", "default", FlagDescription(`
		Unique identifier for Mindreader, so that it can produce 'oneblock files' in the same store as another instance without competing
		for writes. You should set this flag if you have multiple mindreader running, each one should get a unique identifier, the
		hostname value is a good value to use.
	`))
	cmd.Flags().Duration("mindreader-node-wait-upload-complete-on-shutdown", 30*time.Second, "When the mindreader is shutting down, it will wait up to that amount of time for the archiver to finish uploading the blocks before leaving anyway")
	cmd.Flags().String("mindreader-node-merge-threshold-block-age", "24h", "When processing blocks with a blocktime older than this threshold, they will be automatically merged")

	return nil
}

func getMindreaderLogPlugin(
	blockStreamServer *blockstream.Server,
	oneBlockStoreURL string,
	mergedBlockStoreURL string,
	mergeThresholdBlockAge string,
	workingDir string,
	batchStartBlockNum uint64,
	batchStopBlockNum uint64,
	blocksChanCapacity int,
	waitTimeForUploadOnShutdown time.Duration,
	oneBlockFileSuffix string,
	operatorShutdownFunc func(error),
	metricsAndReadinessManager *nodeManager.MetricsAndReadinessManager,
	tracker *bstream.Tracker,
	appLogger *zap.Logger,
	appTracer logging.Tracer,
) (*mindreader.MindReaderPlugin, error) {
	tracker.AddGetter(bstream.NetworkLIBTarget, func(ctx context.Context) (bstream.BlockRef, error) {
		// FIXME: Need to re-enable the tracker through blockmeta later on (see commented code below), might need to tweak some stuff to make mindreader work...
		return bstream.BlockRefEmpty, nil
	})

	consoleReaderFactory := func(lines chan string) (mindreader.ConsolerReader, error) {
		return codec.NewConsoleReader(lines, func(block *pbarweave.Block) {
			metricsAndReadinessManager.UpdateHeadBlock(block.Height, string(block.IndepHash), block.Time())
		})
	}

	return mindreader.NewMindReaderPlugin(
		oneBlockStoreURL,
		mergedBlockStoreURL,
		mergeThresholdBlockAge,
		workingDir,
		consoleReaderFactory,
		batchStartBlockNum,
		batchStopBlockNum,
		blocksChanCapacity,
		metricsAndReadinessManager.UpdateHeadBlock,
		func(error) {
			operatorShutdownFunc(nil)
		},
		waitTimeForUploadOnShutdown,
		oneBlockFileSuffix,
		blockStreamServer,
		appLogger,
		appTracer,
	)
}
