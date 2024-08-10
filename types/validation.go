package types

import (
	"cosmossdk.io/core/address"
)

func (msg *MsgSetAvailAddress) Validate(ac address.Codec) ([]byte, error) {
	return ac.StringToBytes(msg.ValidatorAddress)
}
