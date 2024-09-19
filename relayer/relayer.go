package relayer

import (
	"time"

	"cosmossdk.io/log"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/codec"
	servertypes "github.com/cosmos/cosmos-sdk/server/types"
	"github.com/vitwit/avail-da-module/relayer/avail"
	http_client "github.com/vitwit/avail-da-module/relayer/http"
	"github.com/vitwit/avail-da-module/relayer/local"
	"github.com/vitwit/avail-da-module/types"
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
	appOpts servertypes.AppOptions,
	nodeDir string,
) (*Relayer, error) {
	cfg := types.AvailConfigFromAppOpts(appOpts)

	// local sdk-based chain provider
	localProvider, err := local.NewProvider(cdc, cfg.CosmosNodeRPC)
	if err != nil {
		return nil, err
	}

	httpClient := http_client.NewHTTPClientHandler()

	// Avail-DA client
	availDAClient := avail.NewLightClient(cfg.LightClientURL, httpClient)

	return &Relayer{
		logger: logger,

		provenHeights: make(chan int64, 10000),
		commitHeights: make(chan int64, 10000),

		localProvider:        localProvider,
		availChainID:         cfg.ChainID,
		submittedBlocksCache: make(map[int64]bool),
		NodeDir:              nodeDir,
		AvailConfig:          cfg,
		availDAClient:        availDAClient,
	}, nil
}

func (r *Relayer) SetClientContext(clientCtx client.Context) {
	r.clientCtx = clientCtx
}
