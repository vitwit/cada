package types

import sdk "github.com/cosmos/cosmos-sdk/types"

const (
	TypeMsgUpdateBlobStatus = "update_blob_Status"
)

var _ sdk.Msg = (*MsgUpdateBlobStatusRequest)(nil)

func NewMsgUpdateBlobStatus(valAddr string, blockRange Range, availHeight uint64, isSuccess bool) *MsgUpdateBlobStatusRequest {
	return &MsgUpdateBlobStatusRequest{
		ValidatorAddress: valAddr,
		BlocksRange:      &blockRange,
		AvailHeight:      availHeight,
		IsSuccess:        isSuccess,
	}
}
