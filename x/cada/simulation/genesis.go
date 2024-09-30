package simulation

import (
	"github.com/cosmos/cosmos-sdk/types/module"
	types "github.com/vitwit/avail-da-module/x/cada/types"
)

// RandomizedGenState creates a randomized GenesisState for testing.
func RandomizedGenState(simState *module.SimulationState) {
	// Since your GenesisState is empty, there's not much to randomize.
	// We'll just set the GenesisState to its empty struct.
	genesis := types.GenesisState{}

	// Here we use simState to set the default genesis
	simState.GenState[types.ModuleName] = simState.Cdc.MustMarshalJSON(&genesis)
}
