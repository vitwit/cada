package types

import (
	servertypes "github.com/cosmos/cosmos-sdk/server/types"
	"github.com/spf13/cast"
)

// AvailConfig defines the configuration for the in-process Avail relayer.
type AvailConfiguration struct {
	// avail light node url
	LightClientURL string `mapstructure:"light-client-url"`

	// Overrides built-in app-id used
	AppID int `mapstructure:"app-id"`

	// avail config
	Seed string `json:"seed"`

	// RPC of the cosmos node to fetch the block data
	CosmosNodeRPC string `json:"cosmos-node-rpc"`

	MaxBlobBlocks uint64 `json:"max-blob-blocks"`

	PublishBlobInterval uint64 `json:"publish-blob-interval"`

	VoteInterval uint64 `json:"vote-interval"`

	ExpirationInterval uint64 `json:"expiration-interval"`
}

const (
	FlagChainID             = "avail.chain-id"
	FlagOverrideAppID       = "avail.override-app-id"
	FlagOverridePubInterval = "avail.override-pub-interval"
	FlagQueryInterval       = "avail.proof-query-interval"
	FlagSeed                = "avail-seed"
	FlagLightClientURL      = "avail.light-client-url"
	FlagCosmosNodeRPC       = "avail.cosmos-node-rpc"
	FlagMaxBlobBlocks       = "avail.max-blob-blocks"
	FlagPublishBlobInterval = "avail.publish-blob-interval"
	FlagVoteInterval        = "avail.vote-interval"
	FlagValidatorKey        = "avail.validator-key"
	FlagExpirationInterval  = "avail.expiration-interval"
	FlagCosmosNodeDir       = "avail.cosmos-node-dir"
	FlagKeyringBackendType  = "avail.keyring-backend-type"

	DefaultConfigTemplate = `

	[avail]

	# Avail light client node url for posting data
	light-client-url = "http://127.0.0.1:8000"

	# Overrides the expected  avail app-id, test-only
	override-app-id = "1"

	# Seed for avail
	seed = "bottom drive obey lake curtain smoke basket hold race lonely fit walk//Alice"

	# RPC of cosmos node to get the block data
	cosmos-node-rpc = "http://127.0.0.1:26657"

	# Maximum number of blocks over which blobs can be processed
	max-blob-blocks = 20

	# The frequency at which block data is posted to the Avail Network
	publish-blob-interval = 5

	# It is the period before validators verify whether data is truly included in
	# Avail and confirm it with the network using vote extension
	vote-interval = 5

	# If the previous blocks status remains pending beyond the expiration interval, it should be marked as expired
	expiration-interval = 30
	`
)

var DefaultAvailConfig = AvailConfiguration{
	Seed:                "bottom drive obey lake curtain smoke basket hold race lonely fit walk//Alice",
	AppID:               1,
	CosmosNodeRPC:       "http://127.0.0.1:26657",
	MaxBlobBlocks:       20,
	PublishBlobInterval: 5,
	VoteInterval:        5,
	ExpirationInterval:  30,
	LightClientURL:      "http://127.0.0.1:8000",
}

func AvailConfigFromAppOpts(appOpts servertypes.AppOptions) AvailConfiguration {
	return AvailConfiguration{
		// ChainID:             cast.ToString(appOpts.Get(FlagChainID)),
		AppID:               cast.ToInt(appOpts.Get(FlagOverrideAppID)),
		Seed:                cast.ToString(appOpts.Get(FlagSeed)),
		LightClientURL:      cast.ToString(appOpts.Get(FlagLightClientURL)),
		CosmosNodeRPC:       cast.ToString(appOpts.Get(FlagCosmosNodeRPC)),
		MaxBlobBlocks:       cast.ToUint64(appOpts.Get(FlagMaxBlobBlocks)),
		PublishBlobInterval: cast.ToUint64(appOpts.Get(FlagPublishBlobInterval)),
		VoteInterval:        cast.ToUint64(appOpts.Get(FlagVoteInterval)),
		ExpirationInterval:  cast.ToUint64(appOpts.Get(FlagExpirationInterval)),
	}
}
