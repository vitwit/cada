package keeper

import (
	"cosmossdk.io/collections"
	storetypes "cosmossdk.io/core/store"
	storetypes2 "cosmossdk.io/store/types"
	upgradekeeper "cosmossdk.io/x/upgrade/keeper"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/spf13/cobra"
	availblob1 "github.com/vitwit/avail-da-module"
	"github.com/vitwit/avail-da-module/relayer"
	"github.com/vitwit/avail-da-module/types"
)

type Keeper struct {
	// stakingKeeper *stakingkeeper.Keeper
	upgradeKeeper *upgradekeeper.Keeper
	relayer       *relayer.Relayer

	Validators              collections.Map[string, string]
	ClientID                collections.Item[string]
	ProvenHeight            collections.Item[uint64]
	PendingBlocksToTimeouts collections.Map[int64, int64]
	TimeoutsToPendingBlocks collections.Map[int64, types.PendingBlocks]
	// keyring                 keyring.Keyring

	storeKey storetypes2.StoreKey

	cdc codec.BinaryCodec

	PublishToAvailBlockInterval uint64
	MaxBlocksForBlob            uint
	VotingInterval              uint64
	appID                       int

	unprovenBlocks map[int64][]byte

	ProposerAddress []byte
	ClientCmd       *cobra.Command
}

func NewKeeper(
	cdc codec.BinaryCodec,
	storeService storetypes.KVStoreService,
	uk *upgradekeeper.Keeper,
	key storetypes2.StoreKey,
	appID int,
) *Keeper {
	sb := collections.NewSchemaBuilder(storeService)

	return &Keeper{
		upgradeKeeper: uk,

		Validators:              collections.NewMap(sb, availblob1.ValidatorsKey, "validators", collections.StringKey, collections.StringValue),
		ClientID:                collections.NewItem(sb, availblob1.ClientIDKey, "client_id", collections.StringValue),
		ProvenHeight:            collections.NewItem(sb, availblob1.ProvenHeightKey, "proven_height", collections.Uint64Value),
		PendingBlocksToTimeouts: collections.NewMap(sb, availblob1.PendingBlocksToTimeouts, "pending_blocks_to_timeouts", collections.Int64Key, collections.Int64Value),
		TimeoutsToPendingBlocks: collections.NewMap(sb, availblob1.TimeoutsToPendingBlocks, "timeouts_to_pending_blocks", collections.Int64Key, codec.CollValue[types.PendingBlocks](cdc)),

		storeKey: key,

		cdc: cdc,

		appID: appID,

		unprovenBlocks:              make(map[int64][]byte),
		MaxBlocksForBlob:            20, // Todo: call this from app.go, later change to params
		PublishToAvailBlockInterval: 5,  // Todo: call this from app.go, later change to params
		VotingInterval:              5,
	}
}

func (k *Keeper) SetRelayer(r *relayer.Relayer) {
	k.relayer = r
}

func (k *Keeper) GetBlobStatus(ctx sdk.Context) uint32 {
	store := ctx.KVStore(k.storeKey)
	return GetStatusFromStore(store)
}
