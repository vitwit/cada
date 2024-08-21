package relayer

import (
	"sync"
	"time"

	"cosmossdk.io/log"
	"github.com/cosmos/cosmos-sdk/client"
	servertypes "github.com/cosmos/cosmos-sdk/server/types"

	"github.com/vitwit/avail-da-module/relayer/local"

	"github.com/cosmos/cosmos-sdk/codec"
)

const (
	DefaultMaxFlushSize = int(20)
	MaxMaxFlushSize     = int(100)
)

// Relayer is responsible for posting new blocks to Avail
type Relayer struct {
	logger log.Logger

	provenHeights      chan int64
	latestProvenHeight int64

	commitHeights      chan int64
	latestCommitHeight int64

	pollInterval         time.Duration
	blockProofCacheLimit int

	mu sync.Mutex

	submittedBlocksCache map[int64]bool

	rpcClient     *AvailClient
	localProvider *local.CosmosProvider
	clientCtx     client.Context

	availChainID string

	availPublishBlockInterval int
	availLastQueriedHeight    int64
	availAppID                int
}

// NewRelayer creates a new Relayer instance
func NewRelayer(
	logger log.Logger,
	cdc codec.BinaryCodec,
	appOpts servertypes.AppOptions,
	keyDir string,
) (*Relayer, error) {
	cfg := AvailConfigFromAppOpts(appOpts)

	client, err := NewAvailClient(cfg)
	if err != nil {
		// panic(fmt.Sprintf("cannot create client:%v", err))
	}

	// local sdk-based chain provider
	localProvider, err := local.NewProvider(cdc)
	if err != nil {
		return nil, err
	}

	return &Relayer{
		logger: logger,

		pollInterval: cfg.ProofQueryInterval,

		provenHeights: make(chan int64, 10000),
		commitHeights: make(chan int64, 10000),

		rpcClient:                 client,
		localProvider:             localProvider,
		availChainID:              cfg.ChainID,
		availLastQueriedHeight:    1, // Defaults to 1, but init genesis can set this based on client state's latest height
		submittedBlocksCache:      make(map[int64]bool),
		availAppID:                cfg.AppID,
		availPublishBlockInterval: 5,
	}, nil
}

func (r *Relayer) SetClientContext(clientCtx client.Context) {
	r.clientCtx = clientCtx
}
