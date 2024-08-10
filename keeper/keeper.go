package keeper

import (
	"cosmossdk.io/collections"
	storetypes2 "cosmossdk.io/store/types"
	upgradekeeper "cosmossdk.io/x/upgrade/keeper"
	"github.com/cosmos/cosmos-sdk/codec"
	availblob1 "github.com/vitwit/avail-da-module"
	"github.com/vitwit/avail-da-module/relayer"
	"github.com/vitwit/avail-da-module/types"

	stakingkeeper "github.com/cosmos/cosmos-sdk/x/staking/keeper"
	// stakingkeeper "cosmossdk.io/x/staking/keeper"

	storetypes "cosmossdk.io/core/store"

	servertypes "github.com/cosmos/cosmos-sdk/server/types"
)

type Keeper struct {
	stakingKeeper *stakingkeeper.Keeper
	upgradeKeeper *upgradekeeper.Keeper
	relayer       *relayer.Relayer

	Validators              collections.Map[string, string]
	ClientID                collections.Item[string]
	ProvenHeight            collections.Item[int64]
	PendingBlocksToTimeouts collections.Map[int64, int64]
	TimeoutsToPendingBlocks collections.Map[int64, types.PendingBlocks]

	storeKey storetypes2.StoreKey

	cdc codec.BinaryCodec

	publishToAvailBlockInterval int
	injectedProofsLimit         int
	appId                       int

	unprovenBlocks map[int64][]byte

	proposerAddress []byte
}

func NewKeeper(
	cdc codec.BinaryCodec,
	appOpts servertypes.AppOptions,
	storeService storetypes.KVStoreService,
	// sk *stakingkeeper.Keeper,
	uk *upgradekeeper.Keeper,
	key storetypes2.StoreKey,
	publishToAvailBlockInterval int,
	appId int,

) *Keeper {
	// cfg := availblob1.AvailConfigFromAppOpts(appOpts)
	sb := collections.NewSchemaBuilder(storeService)

	// if cfg.OverridePubInterval > 0 {
	// 	publishToAvailBlockInterval = cfg.OverridePubInterval
	// }

	return &Keeper{
		// stakingKeeper: sk,
		upgradeKeeper: uk,

		Validators:              collections.NewMap(sb, availblob1.ValidatorsKey, "validators", collections.StringKey, collections.StringValue),
		ClientID:                collections.NewItem(sb, availblob1.ClientIDKey, "client_id", collections.StringValue),
		ProvenHeight:            collections.NewItem(sb, availblob1.ProvenHeightKey, "proven_height", collections.Int64Value),
		PendingBlocksToTimeouts: collections.NewMap(sb, availblob1.PendingBlocksToTimeouts, "pending_blocks_to_timeouts", collections.Int64Key, collections.Int64Value),
		TimeoutsToPendingBlocks: collections.NewMap(sb, availblob1.TimeoutsToPendingBlocks, "timeouts_to_pending_blocks", collections.Int64Key, codec.CollValue[types.PendingBlocks](cdc)),

		storeKey: key,

		cdc: cdc,

		publishToAvailBlockInterval: publishToAvailBlockInterval,
		// injectedProofsLimit:            cfg.MaxFlushSize,
		appId: appId,

		unprovenBlocks: make(map[int64][]byte),
	}
}

func (k *Keeper) SetRelayer(r *relayer.Relayer) {
	k.relayer = r
}
