module github.com/streamingfast/firehose-arweave

go 1.16

require (
	github.com/ShinyTrinkets/overseer v0.3.0
	github.com/dvsekhvalnov/jose2go v1.5.0
	github.com/golang/protobuf v1.5.2
	github.com/spf13/cobra v1.4.0
	github.com/spf13/viper v1.8.1
	github.com/streamingfast/bstream v0.0.2-0.20220419181641-fdf5ab55791d
	github.com/streamingfast/cli v0.0.4-0.20220113202443-f7bcefa38f7e
	github.com/streamingfast/dauth v0.0.0-20220404140613-a40f4cd81626
	github.com/streamingfast/derr v0.0.0-20220307162255-f277e08753fa
	github.com/streamingfast/dgrpc v0.0.0-20220307180102-b2d417ac8da7
	github.com/streamingfast/dlauncher v0.0.0-20220307153121-5674e1b64d40
	github.com/streamingfast/dmetering v0.0.0-20220307162406-37261b4b3de9
	github.com/streamingfast/dmetrics v0.0.0-20220307162521-2389094ab4a1
	github.com/streamingfast/dstore v0.1.1-0.20220315134935-980696943a79
	github.com/streamingfast/firehose v0.1.1-0.20220524194622-31eac090a1a5
	github.com/streamingfast/firehose-arweave/types v0.0.0-20220509041238-3d3270820c99
	github.com/streamingfast/logging v0.0.0-20220304214715-bc750a74b424
	github.com/streamingfast/merger v0.0.3-0.20220510150626-2e0bad630abf
	github.com/streamingfast/node-manager v0.0.2-0.20220506192940-9e2e278b6008
	github.com/streamingfast/pbgo v0.0.6-0.20220428192744-f80aee7d4688
	github.com/streamingfast/relayer v0.0.2-0.20220307182103-5f4178c54fde
	github.com/streamingfast/sf-tools v0.0.0-20220510152242-8343cb8e91aa
	github.com/stretchr/testify v1.7.1-0.20210427113832-6241f9ab9942
	go.uber.org/zap v1.21.0
	google.golang.org/grpc v1.44.0
	google.golang.org/protobuf v1.27.1
)

replace github.com/ShinyTrinkets/overseer => github.com/streamingfast/overseer v0.2.1-0.20210326144022-ee491780e3ef

replace github.com/streamingfast/dauth => github.com/eosnationftw/dauth v0.0.0-20220525063810-eca84ecd2ac5
