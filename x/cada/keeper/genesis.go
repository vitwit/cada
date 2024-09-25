package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/vitwit/avail-da-module/x/cada/types"
)

// InitGenesis initializes the module's state from a genesis state.
func (k *Keeper) InitGenesis(_ sdk.Context, _ *types.GenesisState) error {
	return nil
}

// ExportGenesis exports the module's state to a genesis state.
func (k *Keeper) ExportGenesis(_ sdk.Context) *types.GenesisState {
	return &types.GenesisState{}
}
