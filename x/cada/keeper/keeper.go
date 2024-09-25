package keeper

import (
	"cosmossdk.io/collections"
	storetypes "cosmossdk.io/core/store"
	"cosmossdk.io/log"
	storetypes2 "cosmossdk.io/store/types"
	upgradekeeper "cosmossdk.io/x/upgrade/keeper"
	"github.com/cosmos/cosmos-sdk/codec"
	servertypes "github.com/cosmos/cosmos-sdk/server/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/spf13/cobra"
	"github.com/vitwit/avail-da-module/relayer"
	types "github.com/vitwit/avail-da-module/x/cada/types"
)

type Keeper struct {
	upgradeKeeper *upgradekeeper.Keeper

	relayer *relayer.Relayer

	Validators collections.Map[string, string]

	storeKey storetypes2.StoreKey

	cdc codec.BinaryCodec

	unprovenBlocks map[int64][]byte

	ProposerAddress []byte

	ClientCmd *cobra.Command
}

func NewKeeper(
	cdc codec.BinaryCodec,
	storeService storetypes.KVStoreService,
	uk *upgradekeeper.Keeper,
	key storetypes2.StoreKey,
	_ servertypes.AppOptions,
	_ log.Logger,
	relayer *relayer.Relayer,
) *Keeper {
	sb := collections.NewSchemaBuilder(storeService)

	return &Keeper{
		upgradeKeeper: uk,

		Validators: collections.NewMap(sb, types.ValidatorsKey, "validators", collections.StringKey, collections.StringValue),

		storeKey: key,

		cdc:            cdc,
		unprovenBlocks: make(map[int64][]byte),
		relayer:        relayer,
	}
}

// SetRelayer sets the relayer instance for the Keeper.
func (k *Keeper) SetRelayer(r *relayer.Relayer) {
	k.relayer = r
}

// GetBlobStatus retrieves the current status of the blob from the store.
func (k *Keeper) GetBlobStatus(ctx sdk.Context) uint32 {
	store := ctx.KVStore(k.storeKey)
	return GetStatusFromStore(store)
}
