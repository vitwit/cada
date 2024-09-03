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

	// Set proven height to genesis height, we do not init any pending block on a genesis init/restart
	// if err := k.SetProvenHeight(ctx, ctx.HeaderInfo().Height); err != nil {
	// 	return err
	// }

	k.relayer.NotifyProvenHeight(ctx.HeaderInfo().Height)

	// TODO: client state
	k.SetAvailGenesisState(ctx, data)

	return nil
}

// ExportGenesis exports the module's state to a genesis state.
func (k *Keeper) ExportGenesis(ctx sdk.Context) *types.GenesisState {
	vals, err := k.GetAllValidators(ctx)
	if err != nil {
		panic(err)
	}

	provenHeight, err := k.GetProvenHeight(ctx)
	if err != nil {
		panic(err)
	}

	pendingBlocks, err := k.GetPendingBlocksWithExpiration(ctx)
	if err != nil {
		panic(err)
	}

	return &types.GenesisState{
		Validators:   vals.Validators,
		ProvenHeight: provenHeight,

		PendingBlocks: pendingBlocks,
	}
}

// SetAvailGenesisState imports avail light client's full state
func (k Keeper) SetAvailGenesisState(ctx sdk.Context, gs *types.GenesisState) {
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
