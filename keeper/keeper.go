package keeper

import (
	"encoding/binary"
	"errors"

	"cosmossdk.io/collections"
	storetypes2 "cosmossdk.io/store/types"
	upgradekeeper "cosmossdk.io/x/upgrade/keeper"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/crypto/keyring"
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
	keyring                 keyring.Keyring

	storeKey storetypes2.StoreKey

	cdc codec.BinaryCodec

	publishToAvailBlockInterval int
	PublishToAvailBlockInterval uint64
	MaxBlocksForBlob            uint
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

		unprovenBlocks:              make(map[int64][]byte),
		MaxBlocksForBlob:            10, //Todo: call this from app.go, later change to params
		PublishToAvailBlockInterval: 5,  //Todo: call this from app.go, later change to params
	}
}

func (k *Keeper) SetRelayer(r *relayer.Relayer) {
	k.relayer = r
}

func (k *Keeper) SetBlobStatusPending(ctx sdk.Context, provenHeight, endHeight uint64) error {

	store := ctx.KVStore(k.storeKey)

	if IsStateReady(store) { //TOodo: we should check for expiration too
		return errors.New("a block range with same start height is already being processed")
	}

	UpdateBlobStatus(ctx, store, PENDING_STATE)
	UpdateEndHeight(ctx, store, endHeight)
	return nil
}

func (k *Keeper) GetProvenHeightFromStore(ctx sdk.Context) uint64 {
	store := ctx.KVStore(k.storeKey)
	heightBytes := store.Get(availblob1.ProvenHeightKey)
	if heightBytes == nil || len(heightBytes) == 0 {
		return 0
	}

	provenHeight := binary.BigEndian.Uint64(heightBytes)
	return provenHeight
}

// Todo: remove this method later
func (k *Keeper) SubmitBlob(ctx sdk.Context, req *types.MsgSubmitBlobRequest) (*types.MsgSubmitBlobResponse, error) {

	return &types.MsgSubmitBlobResponse{}, nil
}

func (k *Keeper) UpdateBlobStatus(ctx sdk.Context, req *types.MsgUpdateBlobStatusRequest) (*types.MsgUpdateBlobStatusResponse, error) {
	//Todo: status should be changed to Voting or Ready, depending on the request
	return &types.MsgUpdateBlobStatusResponse{}, nil
}
