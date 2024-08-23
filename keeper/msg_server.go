package keeper

import (
	"context"
	"fmt"

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

func (s msgServer) SetAvailAddress(ctx context.Context, msg *types.MsgSetAvailAddress) (*types.MsgSetAvailAddressResponse, error) {
	valAddr, err := msg.Validate(s.k.stakingKeeper.ValidatorAddressCodec())
	if err != nil {
		return nil, err
	}

	// verify that the validator exists
	if _, err := s.k.stakingKeeper.GetValidator(ctx, valAddr); err != nil {
		return nil, err
	}

	if err = s.k.SetValidatorAvailAddress(ctx, types.Validator{
		ValidatorAddress: msg.ValidatorAddress,
		AvailAddress:     msg.AvailAddress,
	}); err != nil {
		return nil, err
	}

	return new(types.MsgSetAvailAddressResponse), nil
}

func (s msgServer) SubmitBlob(ctx context.Context, req *types.MsgSubmitBlobRequest) (*types.MsgSubmitBlobResponse, error) {

	fmt.Println("msgServer sub,itblob.............", req)
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	return s.k.SubmitBlob(sdkCtx, req)
}

func (s msgServer) UpdateBlobStatus(ctx context.Context, req *types.MsgUpdateBlobStatusRequest) (*types.MsgUpdateBlobStatusResponse, error) {
	sdkCtx := sdk.UnwrapSDKContext(ctx)

	//TODO: query the light client
	return s.k.UpdateBlobStatus(sdkCtx, req)
}

/*
rpc SubmitBlob(MsgSubmitBlobRequest) returns (MsgSubmitBlobResponse);

  // UpdateBlobStatus
  rpc UpdateBlobStatus(MsgUpdateBlobStatusRequest) returns (MsgUpdateBlobStatusResponse);
*/
