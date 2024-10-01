package module

import (
	"cosmossdk.io/core/appmodule"
	"cosmossdk.io/core/store"
	"cosmossdk.io/depinject"
	upgradekeeper "cosmossdk.io/x/upgrade/keeper"
	"github.com/cosmos/cosmos-sdk/codec"
)

type Inputs struct {
	depinject.In

	Cdc          codec.Codec
	StoreService store.KVStoreService

	UpgradeKeeper upgradekeeper.Keeper
}

type Outputs struct {
	depinject.Out

	Module appmodule.AppModule
}
