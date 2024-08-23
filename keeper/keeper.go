package keeper

import (
	"errors"
	"fmt"

	"cosmossdk.io/collections"
	storetypes2 "cosmossdk.io/store/types"
	upgradekeeper "cosmossdk.io/x/upgrade/keeper"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/spf13/cobra"
	availblob1 "github.com/vitwit/avail-da-module"
	"github.com/vitwit/avail-da-module/relayer"
	"github.com/vitwit/avail-da-module/types"

	stakingkeeper "github.com/cosmos/cosmos-sdk/x/staking/keeper"

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
	app_id                      int

	unprovenBlocks map[int64][]byte

	proposerAddress []byte
	ClientCmd       *cobra.Command
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
	sb := collections.NewSchemaBuilder(storeService)

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
		app_id:                      appId,

		unprovenBlocks: make(map[int64][]byte),
	}
}

func (k *Keeper) SetRelayer(r *relayer.Relayer) {
	k.relayer = r
}

func (k *Keeper) SubmitBlob(ctx sdk.Context, req *types.MsgSubmitBlobRequest) (*types.MsgSubmitBlobResponse, error) {
	store := ctx.KVStore(k.storeKey)
	if IsAlreadyExist(ctx, store, *req.BlocksRange) {
		return &types.MsgSubmitBlobResponse{}, errors.New("the range is already processed")
	}

	err := updateBlobStatus(ctx, store, *req.BlocksRange, PENDING)
	fmt.Println("errr.........", err)
	return &types.MsgSubmitBlobResponse{}, err
}

func (k *Keeper) UpdateBlobStatus(ctx sdk.Context, req *types.MsgUpdateBlobStatusRequest) (*types.MsgUpdateBlobStatusResponse, error) {
	store := ctx.KVStore(k.storeKey)
	if !IsAlreadyExist(ctx, store, *req.BlocksRange) {
		return &types.MsgUpdateBlobStatusResponse{}, errors.New("the range does not exist")
	}

	status := FAILURE
	if req.IsSuccess {
		status = SUCCESS
	}
	err := updateBlobStatus(ctx, store, *req.BlocksRange, status)
	return &types.MsgUpdateBlobStatusResponse{}, err
}
