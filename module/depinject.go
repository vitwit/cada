package module

import (
	upgradekeeper "cosmossdk.io/x/upgrade/keeper"
	"github.com/cosmos/cosmos-sdk/codec"
	servertypes "github.com/cosmos/cosmos-sdk/server/types"

	"cosmossdk.io/core/appmodule"
	"cosmossdk.io/core/store"
	"cosmossdk.io/depinject"
	storetypes "cosmossdk.io/store/types"
)

// var _ appmodule.AppModule = AppModule{}

func init() {
	// appmodule.Register(
	// 	// new(modulev1.Module),
	// 	// appmodule.Provide(ProvideModule),
	// )
}

type ModuleInputs struct {
	depinject.In

	Cdc          codec.Codec
	appOpts      servertypes.AppOptions
	StoreService store.KVStoreService

	// StakingKeeper stakingkeeper.Keeper
	UpgradeKeeper upgradekeeper.Keeper

	storeKey storetypes.StoreKey

	publishToAvailBlockInterval int
}

type ModuleOutputs struct {
	depinject.Out

	Module appmodule.AppModule
	// Keeper *keeper.Keeper
}

// func ProvideModule(in ModuleInputs) ModuleOutputs {
// 	k := keeper.NewKeeper(in.Cdc, in.appOpts, in.StoreService, &in.StakingKeeper, &in.UpgradeKeeper, in.storeKey, in.publishToAvailBlockInterval)
// 	m := NewAppModule(in.Cdc, k)

// 	return ModuleOutputs{Module: m, Keeper: k}
// }
