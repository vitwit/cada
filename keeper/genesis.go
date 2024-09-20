package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/vitwit/avail-da-module/types"
)

// InitGenesis initializes the module's state from a genesis state.
func (k *Keeper) InitGenesis(ctx sdk.Context, data *types.GenesisState) error {
	for _, v := range data.Validators {
		if err := k.SetValidatorAvailAddress(ctx, v); err != nil {
			return err
		}
	}

	k.SetAvailGenesisState(ctx, data)

	return nil
}

// ExportGenesis exports the module's state to a genesis state.
func (k *Keeper) ExportGenesis(ctx sdk.Context) *types.GenesisState {
	vals, err := k.GetAllValidators(ctx)
	if err != nil {
		panic(err)
	}

	// provenHeight, err := k.GetProvenHeight(ctx)
	// if err != nil {
	// 	panic(err)
	// }

	return &types.GenesisState{
		Validators: vals.Validators,
		// ProvenHeight: provenHeight,
	}
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
