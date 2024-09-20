package keeper

import (
	"context"

	"github.com/vitwit/avail-da-module/types"
)

func (k *Keeper) SetValidatorAvailAddress(ctx context.Context, validator types.Validator) error {
	return k.Validators.Set(ctx, validator.ValidatorAddress, validator.AvailAddress)
}

func (k *Keeper) GetValidatorAvailAddress(ctx context.Context, validatorAddress string) (string, error) {
	return k.Validators.Get(ctx, validatorAddress)
}

func (k *Keeper) GetAllValidators(ctx context.Context) (types.Validators, error) {
	var validators types.Validators
	it, err := k.Validators.Iterate(ctx, nil)
	if err != nil {
		return validators, err
	}

	defer it.Close()

	for ; it.Valid(); it.Next() {
		var validator types.Validator
		validator.ValidatorAddress, err = it.Key()
		if err != nil {
			return validators, err
		}
		validator.AvailAddress, err = it.Value()
		if err != nil {
			return validators, err
		}
		validators.Validators = append(validators.Validators, validator)
	}

	return validators, nil
}

// func (k *Keeper) SetProvenHeight(ctx context.Context, height uint64) error {
// 	return k.ProvenHeight.Set(ctx, height)
// }

// func (k *Keeper) GetProvenHeight(ctx context.Context) (uint64, error) {
// 	return k.ProvenHeight.Get(ctx)
// }

// func (k *Keeper) SetClientID(ctx context.Context, clientID string) error {
// 	return k.ClientID.Set(ctx, clientID)
// }

// func (k *Keeper) GetClientID(ctx context.Context) (string, error) {
// 	return k.ClientID.Get(ctx)
// }

// IsBlockPending return true if a block height is already pending
// func (k Keeper) IsBlockPending(ctx context.Context, blockHeight int64) bool {
// 	found, err := k.PendingBlocksToTimeouts.Has(ctx, blockHeight)
// 	if err != nil {
// 		return false
// 	}
// 	return found
// }

// IsBlockExpired will return true if a block is pending and expired, otherwise it returns false
// func (k *Keeper) IsBlockExpired(ctx context.Context, currentBlockTime time.Time, blockHeight int64) bool {
// 	currentBlockTimeNs := currentBlockTime.UnixNano()
// 	found, err := k.PendingBlocksToTimeouts.Has(ctx, blockHeight)
// 	if err != nil {
// 		return false
// 	}
// 	if found {
// 		expiration, err := k.PendingBlocksToTimeouts.Get(ctx, blockHeight)
// 		if err != nil {
// 			return false
// 		}
// 		if currentBlockTimeNs >= expiration {
// 			return true
// 		}
// 	}
// 	return false
// }
