package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/vitwit/avail-da-module/types"
)

// InitGenesis initializes the module's state from a genesis state.
func (k *Keeper) InitGenesis(ctx sdk.Context, data *types.GenesisState) error {

	return nil
}

// ExportGenesis exports the module's state to a genesis state.
func (k *Keeper) ExportGenesis(ctx sdk.Context) *types.GenesisState {

	return &types.GenesisState{}
}

// SetAvailGenesisState imports avail light client's full state
func (k Keeper) SetAvailGenesisState(_ sdk.Context, _ *types.GenesisState) {
	// if gs != nil {
	// 	store := ctx.KVStore(k.storeKey)
	// 	for _, metadata := range gs.Metadata {
	// 		store.Set(metadata.Key, metadata.Value)
	// 	}
	// 	avail.SetClientState(store, k.cdc, &gs.ClientState)
	// 	for _, consensusStateWithHeight := range gs.ConsensusStates {
	// 		avail.SetConsensusState(store, k.cdc, &consensusStateWithHeight.ConsensusState, consensusStateWithHeight.Height)
	// 	}
	// }
}
