package module

import (
	"cosmossdk.io/core/appmodule"
	"cosmossdk.io/core/store"
	"cosmossdk.io/depinject"
	upgradekeeper "cosmossdk.io/x/upgrade/keeper"
	"github.com/cosmos/cosmos-sdk/codec"
)

// var _ appmodule.AppModule = AppModule{}

func init() {
	// appmodule.Register(
	// 	// new(modulev1.Module),
	// 	// appmodule.Provide(ProvideModule),
	// )
}

type Inputs struct {
	depinject.In

	Cdc          codec.Codec
	StoreService store.KVStoreService

	// StakingKeeper stakingkeeper.Keeper
	UpgradeKeeper upgradekeeper.Keeper
}

type Outputs struct {
	depinject.Out

	Module appmodule.AppModule
	// Keeper *keeper.Keeper
}
