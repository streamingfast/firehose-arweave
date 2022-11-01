module github.com/streamingfast/firehose-arweave

go 1.16

require (
	github.com/ShinyTrinkets/overseer v0.3.0
	github.com/dvsekhvalnov/jose2go v1.5.0
	github.com/golang/protobuf v1.5.2
	github.com/spf13/cobra v1.4.0
	github.com/spf13/viper v1.8.1
	github.com/streamingfast/bstream v0.0.2-0.20221017131819-2a7e38be1047
	github.com/streamingfast/cli v0.0.4-0.20220630165922-bc58c6666fc8
	github.com/streamingfast/dauth v0.0.0-20220404140613-a40f4cd81626
	github.com/streamingfast/derr v0.0.0-20220526184630-695c21740145
	github.com/streamingfast/dlauncher v0.0.0-20220909121534-7a9aa91dbb32
	github.com/streamingfast/dmetering v0.0.0-20220307162406-37261b4b3de9
	github.com/streamingfast/dmetrics v0.0.0-20220811180000-3e513057d17c
	github.com/streamingfast/dstore v0.1.1-0.20221021155138-4baa2d406146
	github.com/streamingfast/firehose v0.1.1-0.20221101130227-3a0b1980aa0b
	github.com/streamingfast/firehose-arweave/types v0.0.0-20220509041238-3d3270820c99
	github.com/streamingfast/logging v0.0.0-20220511154537-ce373d264338
	github.com/streamingfast/merger v0.0.3-0.20221101144843-b39ece2e2ebc
	github.com/streamingfast/node-manager v0.0.2-0.20220912235129-6c08463b0c01
	github.com/streamingfast/pbgo v0.0.6-0.20221020131607-255008258d28
	github.com/streamingfast/relayer v0.0.2-0.20220909122435-e67fbc964fd9
	github.com/streamingfast/sf-tools v0.0.0-20221020185155-d5fe94d7578e
	github.com/stretchr/testify v1.8.0
	go.uber.org/zap v1.21.0
	google.golang.org/grpc v1.49.0
	google.golang.org/protobuf v1.28.0
)

replace github.com/ShinyTrinkets/overseer => github.com/streamingfast/overseer v0.2.1-0.20210326144022-ee491780e3ef
