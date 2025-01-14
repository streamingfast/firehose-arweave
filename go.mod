module github.com/streamingfast/firehose-arweave

go 1.16

require (
	github.com/ShinyTrinkets/overseer v0.3.0
	github.com/dvsekhvalnov/jose2go v1.5.0
	github.com/golang/protobuf v1.5.2
	github.com/spf13/cobra v1.4.0
	github.com/spf13/viper v1.8.1
	github.com/streamingfast/bstream v0.0.2-0.20230228213106-2b6a3160e01e
	github.com/streamingfast/cli v0.0.4-0.20220630165922-bc58c6666fc8
	github.com/streamingfast/dauth v0.0.0-20221027185237-b209f25fa3ff
	github.com/streamingfast/dbin v0.9.1-0.20220513054835-1abebbb944ad // indirect
	github.com/streamingfast/derr v0.0.0-20221125175206-82e01d420d45
	github.com/streamingfast/dgrpc v0.0.0-20230113212008-1898f17e0ac7 // indirect
	github.com/streamingfast/dlauncher v0.0.0-20220909121534-7a9aa91dbb32
	github.com/streamingfast/dmetering v0.0.0-20220307162406-37261b4b3de9
	github.com/streamingfast/dmetrics v0.0.0-20221129121022-a1733eca1981
	github.com/streamingfast/dstore v0.1.1-0.20230202164314-93694544e2ca
	github.com/streamingfast/firehose v0.1.1-0.20221101130227-3a0b1980aa0b
	github.com/streamingfast/firehose-arweave/types v0.0.0-20220509041238-3d3270820c99
	github.com/streamingfast/logging v0.0.0-20220813175024-b4fbb0e893df
	github.com/streamingfast/merger v0.0.3-0.20221123202507-445dfd357868
	github.com/streamingfast/node-manager v0.0.2-0.20221115101723-d9823ffd7ad5
	github.com/streamingfast/pbgo v0.0.6-0.20221020131607-255008258d28
	github.com/streamingfast/relayer v0.0.2-0.20220909122435-e67fbc964fd9
	github.com/streamingfast/sf-tools v0.0.0-20221129171534-a0708b599ce5
	github.com/stretchr/testify v1.8.0
	go.uber.org/zap v1.21.0
	google.golang.org/grpc v1.50.1
	google.golang.org/protobuf v1.28.1
)

replace github.com/ShinyTrinkets/overseer => github.com/streamingfast/overseer v0.2.1-0.20210326144022-ee491780e3ef
