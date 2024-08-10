package relayer

import (
	"time"

	servertypes "github.com/cosmos/cosmos-sdk/server/types"
	"github.com/spf13/cast"
)

const (
	FlagAppRpcURL           = "avail.app-rpc-url"
	FlagAppRpcTimeout       = "avail.app-rpc-timeout"
	FlagChainID             = "avail.chain-id"
	FlagGasPrices           = "avail.gas-prices"
	FlagGasAdjustment       = "avail.gas-adjustment"
	FlagNodeRpcURL          = "avail.node-rpc-url"
	FlagNodeAuthToken       = "avail.node-auth-token"
	FlagOverrideAppID       = "avail.override-app-id"
	FlagOverridePubInterval = "avail.override-pub-interval"
	FlagQueryInterval       = "avail.proof-query-interval"
	FlagMaxFlushSize        = "avail.max-flush-size"
	FlagSeed                = "avail-seed"
	FlagLightClientURL      = "light_client_url"

	DefaultConfigTemplate = `

	[celestia]
	# RPC URL of celestia-app node for posting block data, querying proofs & light blocks
	app-rpc-url = "https://rpc-mocha.pops.one:443"

	# RPC Timeout for transaction broadcasts and queries to celestia-app node
	app-rpc-timeout = "30s"

	# Celestia chain id
	chain-id = "celestia-1"

	# Gas price to pay for celestia transactions
	gas-prices = "0.01utia"

	# Gas adjustment for celestia transactions
	gas-adjustment = 1.0

	# RPC URL of celestia-node for querying blobs
	node-rpc-url = "http://127.0.0.1:26658"

	# Auth token for celestia-node RPC, n/a if --rpc.skip-auth is used on start
	node-auth-token = "auth-token"

	# Overrides the expected chain's app-id, test-only
	override-app-id = ""

	# Overrides the expected chain's publish-to-celestia block interval, test-only
	override-pub-interval = 0
	
	# Query celestia for new block proofs this often
	proof-query-interval = "12s"

	# Only flush at most this many block proofs in an injected tx per block proposal
	# Must be greater than 0 and less than 100, proofs are roughly 1KB each
	# tiablob will try to aggregate multiple blobs published at the same height w/ a single proof
	max-flush-size = 32

	`
)

var DefaultAvailConfig = AvailConfig{
	AppRpcURL:          "http://127.0.0.1:3000", // TODO remove hardcoded URL
	AppRpcTimeout:      30 * time.Second,
	ChainID:            "avail-1",
	GasPrice:           "0.01utia",
	GasAdjustment:      1.0,
	NodeRpcURL:         "ws://127.0.0.1:8000",
	NodeAuthToken:      "auth-token",
	ProofQueryInterval: 12 * time.Second,
	MaxFlushSize:       32,
	Seed:               "bottom drive obey lake curtain smoke basket hold race lonely fit walk//Alice",
	AppID:              0,
}

// AvailConfig defines the configuration for the in-process Avail relayer.
type AvailConfig struct {
	// RPC URL of celestia-app
	AppRpcURL string `mapstructure:"app-rpc-url"`

	// RPC Timeout for celestia-app
	AppRpcTimeout time.Duration `mapstructure:"app-rpc-timeout"`

	// Celestia chain ID
	ChainID string `mapstructure:"chain-id"`

	// Gas price to pay for Celestia transactions
	GasPrice string `mapstructure:"gas-prices"`

	// Gas adjustment for Celestia transactions
	GasAdjustment float64 `mapstructure:"gas-adjustment"`

	// RPC URL of celestia-node
	NodeRpcURL string `mapstructure:"node-rpc-url"`

	// RPC Timeout for celestia-node
	NodeAuthToken string `mapstructure:"node-auth-token"`

	// Overrides built-in app-id used
	AppID int `mapstructure:"app-id"`

	// Overrides built-in publish-to-celestia block interval
	OverridePubInterval int `mapstructure:"override-pub-interval"`

	// Query Celestia for new block proofs this often
	ProofQueryInterval time.Duration `mapstructure:"proof-query-interval"`

	// Only flush at most this many block proofs in an injected tx per block proposal
	MaxFlushSize int `mapstructure:"max-flush-size"`

	// avail config
	Seed string `json:"seed"`

	LightClientURL string `json:"light_client_url"`
}

func AvailConfigFromAppOpts(appOpts servertypes.AppOptions) AvailConfig {

	return AvailConfig{
		AppRpcURL:           cast.ToString(appOpts.Get(FlagAppRpcURL)),
		AppRpcTimeout:       cast.ToDuration(appOpts.Get(FlagAppRpcTimeout)),
		ChainID:             cast.ToString(appOpts.Get(FlagChainID)),
		GasPrice:            cast.ToString(appOpts.Get(FlagGasPrices)),
		GasAdjustment:       cast.ToFloat64(appOpts.Get(FlagGasAdjustment)),
		NodeRpcURL:          cast.ToString(appOpts.Get(FlagNodeRpcURL)),
		NodeAuthToken:       cast.ToString(appOpts.Get(FlagNodeAuthToken)),
		AppID:               cast.ToInt(appOpts.Get(FlagOverrideAppID)),
		OverridePubInterval: cast.ToInt(appOpts.Get(FlagOverridePubInterval)),
		ProofQueryInterval:  cast.ToDuration(appOpts.Get(FlagQueryInterval)),
		MaxFlushSize:        cast.ToInt(appOpts.Get(FlagMaxFlushSize)),
		Seed:                cast.ToString(appOpts.Get(FlagSeed)),
		LightClientURL:      cast.ToString(appOpts.Get(FlagLightClientURL)),
	}
}
