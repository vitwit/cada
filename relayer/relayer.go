package relayer

import (
	"time"

	"cosmossdk.io/log"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/vitwit/avail-da-module/relayer/avail"
	"github.com/vitwit/avail-da-module/relayer/local"
	"github.com/vitwit/avail-da-module/types"
)

// Relayer is responsible for posting new blocks to Avail
type Relayer struct {
	logger log.Logger

	provenHeights      chan int64
	latestProvenHeight int64

	commitHeights      chan int64
	latestCommitHeight int64

	pollInterval time.Duration

	submittedBlocksCache map[int64]bool

	localProvider *local.CosmosProvider
	availDAClient avail.DA
	clientCtx     client.Context

	availChainID string
	AvailConfig  types.AvailConfiguration
	NodeDir      string
}

// NewRelayer creates a new Relayer instance
func NewRelayer(
	logger log.Logger,
	cdc codec.BinaryCodec,
	cfg types.AvailConfiguration,
	nodeDir string,
	daClient avail.DA,
) (*Relayer, error) {
	// local sdk-based chain provider
	localProvider, err := local.NewProvider(cdc, cfg.CosmosNodeRPC)
	if err != nil {
		return nil, err
	}

	return &Relayer{
		logger: logger,

		provenHeights: make(chan int64, 10000),
		commitHeights: make(chan int64, 10000),

		localProvider:        localProvider,
		availChainID:         cfg.ChainID,
		submittedBlocksCache: make(map[int64]bool),
		NodeDir:              nodeDir,
		AvailConfig:          cfg,
		availDAClient:        daClient,
	}, nil
}

// SetClientContext sets the provided client context for the Relayer.
func (r *Relayer) SetClientContext(clientCtx client.Context) {
	r.clientCtx = clientCtx
}
