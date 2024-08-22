package relayer

import (
	"time"

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

	DefaultConfigTemplate = `

	[avail]

	# Avail light client node url for posting data
	light-client-url = "http://127.0.0.1:8000"

	# Avail chain id
	chain-id = "avail-1"

	# Overrides the expected  avail app-id, test-only
	override-app-id = "1"

	# Overrides the expected chain's publish-to-avail block interval, test-only
	override-pub-interval = 5

	# Seed for avail
	seed = "bottom drive obey lake curtain smoke basket hold race lonely fit walk//Alice"

	# RPC of cosmos node to get the block data
	cosmos-node-rpc = "http://127.0.0.1:26657"
	`
)

var DefaultAvailConfig = AvailConfig{
	ChainID:            "avail-1",
	ProofQueryInterval: 12 * time.Second,
	Seed:               "bottom drive obey lake curtain smoke basket hold race lonely fit walk//Alice",
	AppID:              0,
	CosmosNodeRPC:      "http://127.0.0.1:26657",
}

// AvailConfig defines the configuration for the in-process Avail relayer.
type AvailConfig struct {

	// avail light node url
	LightClientURL string `mapstructure:"light-client-url"`

	// avail chain ID
	ChainID string `mapstructure:"chain-id"`

	// Overrides built-in app-id used
	AppID int `mapstructure:"app-id"`

	// Overrides built-in publish-to-avail block interval
	OverridePubInterval int `mapstructure:"override-pub-interval"`

	// Query avail for new block proofs this often
	ProofQueryInterval time.Duration `mapstructure:"proof-query-interval"`

	// avail config
	Seed string `json:"seed"`

	// RPC of the cosmos node to fetch the block data
	CosmosNodeRPC string `json:"cosmos-node-rpc"`
}

func AvailConfigFromAppOpts(appOpts servertypes.AppOptions) AvailConfig {

	return AvailConfig{
		ChainID:             cast.ToString(appOpts.Get(FlagChainID)),
		AppID:               cast.ToInt(appOpts.Get(FlagOverrideAppID)),
		OverridePubInterval: cast.ToInt(appOpts.Get(FlagOverridePubInterval)),
		ProofQueryInterval:  cast.ToDuration(appOpts.Get(FlagQueryInterval)),
		Seed:                cast.ToString(appOpts.Get(FlagSeed)),
		LightClientURL:      cast.ToString(appOpts.Get(FlagLightClientURL)),
		CosmosNodeRPC:       cast.ToString(appOpts.Get(FlagCosmosNodeRPC)),
	}
}
