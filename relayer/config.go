package relayer

import (
	servertypes "github.com/cosmos/cosmos-sdk/server/types"
	"github.com/spf13/cast"
)

const (
	FlagChainID             = "avail.chain-id"
	FlagOverrideAppID       = "avail.override-app-id"
	FlagOverridePubInterval = "avail.override-pub-interval"
	FlagQueryInterval       = "avail.proof-query-interval"
	FlagSeed                = "avail-seed"
	FlagLightClientURL      = "avail.light-client-url"
	FlagCosmosNodeRPC       = "avail.cosmos-node-rpc"
	FlagMaxBlobBlocks       = "avail.max-blob-blocks"
	FlagBlobInterval        = "avail.blob-interval"
	FlagVoteInterval        = "avail.vote-interval"

	DefaultConfigTemplate = `

	[avail]

	# Avail light client node url for posting data
	light-client-url = "http://127.0.0.1:8000"

	# Avail chain id
	chain-id = "avail-1"

	# Overrides the expected  avail app-id, test-only
	override-app-id = "1"

	# Seed for avail
	seed = "bottom drive obey lake curtain smoke basket hold race lonely fit walk//Alice"

	# RPC of cosmos node to get the block data
	cosmos-node-rpc = "http://127.0.0.1:26657"

	# Maximum number of blocks over which blobs can be processed
	max-blob-blocks = 10

	# The frequency at which block data is posted to the Avail Network
	blob-interval = 5

	# It is the period before validators verify whether data is truly included in
	# Avail and confirm it with the network using vote extension
	vote-interval = 5
	`
)

var DefaultAvailConfig = AvailConfig{
	ChainID:        "avail-1",
	Seed:           "bottom drive obey lake curtain smoke basket hold race lonely fit walk//Alice",
	AppID:          0,
	CosmosNodeRPC:  "http://127.0.0.1:26657",
	MaxBlobBlocks:  20,
	BlobInterval:   10,
	VoteInterval:   5,
	LightClientURL: "http://127.0.0.1:8000",
}

// AvailConfig defines the configuration for the in-process Avail relayer.
type AvailConfig struct {
	// avail light node url
	LightClientURL string `mapstructure:"light-client-url"`

	// avail chain ID
	ChainID string `mapstructure:"chain-id"`

	// Overrides built-in app-id used
	AppID int `mapstructure:"app-id"`

	// avail config
	Seed string `json:"seed"`

	// RPC of the cosmos node to fetch the block data
	CosmosNodeRPC string `json:"cosmos-node-rpc"`

	MaxBlobBlocks uint64 `json:"max-blob-blocks"`

	BlobInterval uint64 `json:"blob-interval"`

	VoteInterval uint64 `json:"vote-interval"`
}

func AvailConfigFromAppOpts(appOpts servertypes.AppOptions) AvailConfig {
	return AvailConfig{
		ChainID:        cast.ToString(appOpts.Get(FlagChainID)),
		AppID:          cast.ToInt(appOpts.Get(FlagOverrideAppID)),
		Seed:           cast.ToString(appOpts.Get(FlagSeed)),
		LightClientURL: cast.ToString(appOpts.Get(FlagLightClientURL)),
		CosmosNodeRPC:  cast.ToString(appOpts.Get(FlagCosmosNodeRPC)),
		MaxBlobBlocks:  cast.ToUint64(appOpts.Get(FlagMaxBlobBlocks)),
		BlobInterval:   cast.ToUint64(appOpts.Get(FlagBlobInterval)),
		VoteInterval:   cast.ToUint64(appOpts.Get(FlagVoteInterval)),
	}
}
