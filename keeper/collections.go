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
