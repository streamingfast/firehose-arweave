package cli

const (
	// Common ports
	MetricsListenAddr string = ":9102"

	// Firehose chain specific port
	//
	// The initial 16XXX prefix is different for every chain supported by Firehose.
	// The current prefix is the one you should use for your chain. Once you have copied
	// this whole repository, you should open a PR on firehose-arweave to bump it again
	// so the next team supporting Firehose will use 16XXXX and so forth.
	ReaderGRPCAddr           string = ":16010"
	NodeManagerAPIAddr       string = ":16041"
	ReaderNodeManagerAPIAddr string = ":16009"
	MergerServingAddr        string = ":16012"
	RelayerServingAddr       string = ":16011"
	FirehoseGRPCServingAddr  string = ":16042"

	// Data storage default locations
	ATMDirectory         string = "file://{data-dir}/atm"
	MergedBlocksStoreURL string = "file://{data-dir}/storage/merged-blocks"
	OneBlockStoreURL     string = "file://{data-dir}/storage/one-blocks"

	// Tweak this for your chain
	FirstStreamableBlock int = 0

	// Native node instance port definitions, adjust those for your chain
	// usually all chains have a P2P and RPC port available

	ReaderNodeP2PPort string = "30305"
	ReaderNodeRPCPort string = "8547"

	NodeP2PPort string = "30303"
	NodeRPCPort string = "3030"
	NodeRPCAddr string = "http://localhost:3030"
)
