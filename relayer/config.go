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
	FlagMaxFlushSize        = "avail.max-flush-size"
	FlagSeed                = "avail-seed"
	FlagLightClientURL      = "avail.light-client-url"

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

	`
)

var DefaultAvailConfig = AvailConfig{
	ChainID:            "avail-1",
	ProofQueryInterval: 12 * time.Second,
	MaxFlushSize:       32,
	Seed:               "bottom drive obey lake curtain smoke basket hold race lonely fit walk//Alice",
	AppID:              0,
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

	// Only flush at most this many block proofs in an injected tx per block proposal
	MaxFlushSize int `mapstructure:"max-flush-size"`

	// avail config
	Seed string `json:"seed"`
}

func AvailConfigFromAppOpts(appOpts servertypes.AppOptions) AvailConfig {

	return AvailConfig{
		ChainID:             cast.ToString(appOpts.Get(FlagChainID)),
		AppID:               cast.ToInt(appOpts.Get(FlagOverrideAppID)),
		OverridePubInterval: cast.ToInt(appOpts.Get(FlagOverridePubInterval)),
		ProofQueryInterval:  cast.ToDuration(appOpts.Get(FlagQueryInterval)),
		MaxFlushSize:        cast.ToInt(appOpts.Get(FlagMaxFlushSize)),
		Seed:                cast.ToString(appOpts.Get(FlagSeed)),
		LightClientURL:      cast.ToString(appOpts.Get(FlagLightClientURL)),
	}
}
