package keeper

import (
	"errors"
	"fmt"

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
	ProvenHeight            collections.Item[uint64]
	PendingBlocksToTimeouts collections.Map[int64, int64]
	TimeoutsToPendingBlocks collections.Map[int64, types.PendingBlocks]
	keyring                 keyring.Keyring

	storeKey storetypes2.StoreKey

	cdc codec.BinaryCodec

	PublishToAvailBlockInterval uint64
	MaxBlocksForBlob            uint
	VotingInterval              uint64
	appID                       int

	unprovenBlocks map[int64][]byte

	proposerAddress []byte
	ClientCmd       *cobra.Command
}

func NewKeeper(
	cdc codec.BinaryCodec,
	appOpts servertypes.AppOptions,
	storeService storetypes.KVStoreService,
	uk *upgradekeeper.Keeper,
	key storetypes2.StoreKey,
	publishToAvailBlockInterval int,
	appId int,

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

		appID: appId,

		unprovenBlocks:              make(map[int64][]byte),
		MaxBlocksForBlob:            20, //Todo: call this from app.go, later change to params
		PublishToAvailBlockInterval: 5,  //Todo: call this from app.go, later change to params
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

func (k *Keeper) UpdateBlobStatus(ctx sdk.Context, req *types.MsgUpdateBlobStatusRequest) (*types.MsgUpdateBlobStatusResponse, error) {
	// status should be changed to Voting or Ready, depending on the request
	store := ctx.KVStore(k.storeKey)
	provenHeight := k.GetProvenHeightFromStore(ctx)
	endHeight := k.GetEndHeightFromStore(ctx)
	status := GetStatusFromStore(store)

	if req.BlocksRange.From != provenHeight+1 || req.BlocksRange.To != endHeight {
		return nil, fmt.Errorf("invalid blocks range request: expected range [%d -> %d], got [%d -> %d]",
			provenHeight+1, endHeight, req.BlocksRange.From, req.BlocksRange.To)
	}

	if status != PENDING_STATE {
		return nil, errors.New("can update the status if it is not pending")
	}

	newStatus := IN_VOTING_STATE
	if !req.IsSuccess {
		newStatus = FAILURE_STATE
	} else {
		currentHeight := ctx.BlockHeight()
		UpdateAvailHeight(ctx, store, req.AvailHeight)                            // updates avail height at which the blocks got submitted to DA
		UpdateVotingEndHeight(ctx, store, uint64(currentHeight)+k.VotingInterval) // TODO: Now voting interval is 5, so check whether we can process votes at next block after tx exec.
	}

	UpdateBlobStatus(ctx, store, newStatus) // updates blob status after based on tx exec

	return &types.MsgUpdateBlobStatusResponse{}, nil
}

func (k *Keeper) SubmitBlobStatus(ctx sdk.Context, _ *types.QuerySubmitBlobStatusRequest) (*types.QuerySubmitBlobStatusResponse, error) {
	// Todo: implement query
	store := ctx.KVStore(k.storeKey)
	startHeight := k.GetStartHeightFromStore(ctx)
	endHeight := k.GetEndHeightFromStore(ctx)
	status := GetStatusFromStore(store)
	statusString := ParseStatus(status)
	provenHeight := k.GetProvenHeightFromStore(ctx)
	votingEndHeight := k.GetVotingEndHeightFromStore(ctx)

	return &types.QuerySubmitBlobStatusResponse{
		Range:                &types.Range{From: startHeight, To: endHeight},
		Status:               statusString,
		ProvenHeight:         provenHeight,
		LastBlobVotingEndsAt: votingEndHeight,
	}, nil
}
