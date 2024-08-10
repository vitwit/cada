package relayer

import (
	"fmt"
	"sync"
	"time"

	"cosmossdk.io/log"
	"github.com/cosmos/cosmos-sdk/client"
	servertypes "github.com/cosmos/cosmos-sdk/server/types"

	// appns "github.com/rollchains/tiablob/celestia/namespace"
	// "github.com/rollchains/tiablob/lightclients/celestia"
	// celestiaprovider "github.com/rollchains/tiablob/relayer/celestia"
	"github.com/vitwit/avail-da-module/relayer/local"

	"github.com/cosmos/cosmos-sdk/codec"
)

const (
	DefaultMaxFlushSize = int(20)
	MaxMaxFlushSize     = int(100)
)

// Relayer is responsible for posting new blocks to Celestia and relaying block proofs from Celestia via the current proposer
type Relayer struct {
	logger log.Logger

	provenHeights      chan int64
	latestProvenHeight int64

	commitHeights      chan int64
	latestCommitHeight int64

	pollInterval         time.Duration
	blockProofCacheLimit int
	nodeRpcUrl           string
	nodeAuthToken        string

	// These items are shared state, must be access with mutex
	// blockProofCache     map[int64]*celestia.BlobProof
	// celestiaHeaderCache map[int64]*celestia.Header
	// updateClient        *celestia.Header
	mu sync.Mutex

	rpcClient     *AvailClient
	localProvider *local.CosmosProvider
	clientCtx     client.Context

	availChainID string
	// celestiaNamespace            appns.Namespace
	availGasPrice             string
	availGasAdjustment        float64
	availPublishBlockInterval int
	availLastQueriedHeight    int64
	availAppID                int
}

// NewRelayer creates a new Relayer instance
func NewRelayer(
	logger log.Logger,
	cdc codec.BinaryCodec,
	appOpts servertypes.AppOptions,
	// celestiaNamespace appns.Namespace,
	keyDir string,
	// celestiaPublishBlockInterval int,
) (*Relayer, error) {
	cfg := AvailConfigFromAppOpts(appOpts)

	fmt.Println("config above..........", cfg)

	cfg.AppID = 1
	cfg.AppRpcTimeout = 100 * time.Second
	cfg.AppRpcURL = "ws://127.0.0.1:9944"
	cfg.ChainID = "avail-1"
	cfg.GasPrice = "0.0utia"
	cfg.GasAdjustment = 1.0
	cfg.NodeRpcURL = "http://127.0.0.1:26657"
	cfg.ProofQueryInterval = 12 * time.Second
	cfg.MaxFlushSize = 32
	cfg.Seed = "bottom drive obey lake curtain smoke basket hold race lonely fit walk//Alice"

	if cfg.MaxFlushSize < 1 || cfg.MaxFlushSize > MaxMaxFlushSize {
		cfg.MaxFlushSize = DefaultMaxFlushSize
		//panic(fmt.Sprintf("invalid relayer max flush size: %d", cfg.MaxFlushSize))
	}

	client, err := NewAvailClient(cfg)
	if err != nil {
		panic(fmt.Sprintf("cannot create client:%v", err))
	}

	// local sdk-based chain provider
	localProvider, err := local.NewProvider(cdc)
	if err != nil {
		return nil, err
	}

	// if cfg.OverrideNamespace != "" {
	// 	celestiaNamespace = appns.MustNewV0([]byte(cfg.OverrideNamespace))
	// }

	// if cfg.OverridePubInterval > 0 {
	// 	celestiaPublishBlockInterval = cfg.OverridePubInterval
	// }

	return &Relayer{
		logger: logger,

		pollInterval: cfg.ProofQueryInterval,

		provenHeights: make(chan int64, 10000),
		commitHeights: make(chan int64, 10000),

		rpcClient:          client,
		localProvider:      localProvider,
		availChainID:       cfg.ChainID,
		availGasPrice:      cfg.GasPrice,
		availGasAdjustment: cfg.GasAdjustment,
		// celestiaPublishBlockInterval: celestiaPublishBlockInterval,
		availLastQueriedHeight: 1, // Defaults to 1, but init genesis can set this based on client state's latest height

		nodeRpcUrl:    cfg.NodeRpcURL,
		nodeAuthToken: cfg.NodeAuthToken,

		// blockProofCache:      make(map[int64]*celestia.BlobProof),
		blockProofCacheLimit: cfg.MaxFlushSize,
		// celestiaHeaderCache:  make(map[int64]*celestia.Header),
		availAppID:                cfg.AppID,
		availPublishBlockInterval: 10,
	}, nil
}

func (r *Relayer) SetClientContext(clientCtx client.Context) {
	r.clientCtx = clientCtx
}
