package avail

import (
	"sync"
	"time"

	"github.com/99designs/keyring"
	gsrpc "github.com/centrifuge/go-substrate-rpc-client/v4"
	// provtypes "github.com/tendermint/tendermint/light/provider"
	// prov "github.com/tendermint/tendermint/light/provider/http"
)

type CosmosProvider struct {
	cdc Codec
	//lightProvider provtypes.Provider
	rpcClient *gsrpc.SubstrateAPI
	keybase   keyring.Keyring

	keyDir string

	walletStateMap map[string]*WalletState
	walletStateMu  sync.Mutex
}

type WalletState struct {
	NextAccountSequence uint64
	Mu                  sync.Mutex
}

func (ws *WalletState) updateNextAccountSequence(seq uint64) {
	if seq > ws.NextAccountSequence {
		ws.NextAccountSequence = seq
	}
}

// NewProvider validates the CosmosProviderConfig, instantiates a ChainClient and then instantiates a CosmosProvider
func NewProvider(rpcURL string, keyDir string, timeout time.Duration, chainID string) (*CosmosProvider, error) {
	// lightProvider, err := prov.New(chainID, rpcURL)
	// if err != nil {
	// 	return nil, err
	// }

	cp := &CosmosProvider{
		cdc: makeCodec(ModuleBasics),
		// lightProvider: lightProvider,
		// rpcClient:      rpcClient,
		keyDir:         keyDir,
		walletStateMap: make(map[string]*WalletState),
	}

	return cp, nil
}
