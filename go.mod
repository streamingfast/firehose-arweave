module github.com/streamingfast/firehose-arweave

go 1.16

require (
	cloud.google.com/go/iam v0.3.0 // indirect
	github.com/ShinyTrinkets/overseer v0.3.0
	github.com/abourget/llerrgroup v0.2.0 // indirect
	github.com/dvsekhvalnov/jose2go v1.5.0
	github.com/golang/groupcache v0.0.0-20210331224755-41bb18bfe9da // indirect
	github.com/golang/protobuf v1.5.2
	github.com/kr/text v0.2.0 // indirect
	github.com/mattn/go-colorable v0.1.12 // indirect
	github.com/spf13/cobra v1.4.0
	github.com/spf13/viper v1.8.1
	github.com/streamingfast/bstream v0.0.2-0.20220809161028-014a633a8d8e
	github.com/streamingfast/cli v0.0.4-0.20220113202443-f7bcefa38f7e
	github.com/streamingfast/dauth v0.0.0-20220404140613-a40f4cd81626
	github.com/streamingfast/derr v0.0.0-20220526184630-695c21740145
	github.com/streamingfast/dgrpc v0.0.0-20220307180102-b2d417ac8da7
	github.com/streamingfast/dlauncher v0.0.0-20220307153121-5674e1b64d40
	github.com/streamingfast/dmetering v0.0.0-20220307162406-37261b4b3de9
	github.com/streamingfast/dmetrics v0.0.0-20220307162521-2389094ab4a1
	github.com/streamingfast/dstore v0.1.1-0.20220607202639-35118aeaf648
	github.com/streamingfast/firehose v0.1.1-0.20220804184723-a790c529fe15
	github.com/streamingfast/firehose-arweave/types v0.0.0-20220509041238-3d3270820c99
	github.com/streamingfast/logging v0.0.0-20220511154537-ce373d264338
	github.com/streamingfast/merger v0.0.3-0.20220803202246-1277c51d3487
	github.com/streamingfast/node-manager v0.0.2-0.20220804015313-01ef0ea2678c
	github.com/streamingfast/pbgo v0.0.6-0.20220801202203-c32e42ac42a8
	github.com/streamingfast/relayer v0.0.2-0.20220802193804-8c63614023a9
	github.com/streamingfast/sf-tools v0.0.0-20220808190933-7bf947e2cc81
	github.com/stretchr/testify v1.7.1
	go.uber.org/zap v1.21.0
	google.golang.org/grpc v1.44.0
	google.golang.org/protobuf v1.27.1
)

replace github.com/ShinyTrinkets/overseer => github.com/streamingfast/overseer v0.2.1-0.20210326144022-ee491780e3ef
