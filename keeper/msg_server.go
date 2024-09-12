package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/vitwit/avail-da-module/types"
)

type msgServer struct {
	k *Keeper
}

var _ types.MsgServer = msgServer{}

// NewMsgServerImpl returns an implementation of the module MsgServer interface.
func NewMsgServerImpl(keeper *Keeper) types.MsgServer {
	return &msgServer{k: keeper}
}

func (s msgServer) UpdateBlobStatus(ctx context.Context, req *types.MsgUpdateBlobStatusRequest) (*types.MsgUpdateBlobStatusResponse, error) {
	sdkCtx := sdk.UnwrapSDKContext(ctx)

	return s.k.UpdateBlobStatus(sdkCtx, req)
}
