package keeper

import (
	"context"
	"fmt"
	"time"

	"cosmossdk.io/collections"
	"github.com/vitwit/avail-da-module/types"
)

const (
	// Window for a transaction to be committed
	ResubmissionTime = 75 * time.Second

	// Buffer for relayer polling logic to retrieve a proof
	RelayerPollingBuffer = 15 * time.Second
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

func (k *Keeper) SetProvenHeight(ctx context.Context, height uint64) error {
	return k.ProvenHeight.Set(ctx, height)
}

func (k *Keeper) GetProvenHeight(ctx context.Context) (uint64, error) {
	return k.ProvenHeight.Get(ctx)
}

func (k *Keeper) SetClientID(ctx context.Context, clientID string) error {
	return k.ClientID.Set(ctx, clientID)
}

func (k *Keeper) GetClientID(ctx context.Context) (string, error) {
	return k.ClientID.Get(ctx)
}

// IsBlockPending return true if a block height is already pending
func (k Keeper) IsBlockPending(ctx context.Context, blockHeight int64) bool {
	found, err := k.PendingBlocksToTimeouts.Has(ctx, blockHeight)
	if err != nil {
		return false
	}
	return found
}

// IsBlockExpired will return true if a block is pending and expired, otherwise it returns false
func (k *Keeper) IsBlockExpired(ctx context.Context, currentBlockTime time.Time, blockHeight int64) bool {
	currentBlockTimeNs := currentBlockTime.UnixNano()
	found, err := k.PendingBlocksToTimeouts.Has(ctx, blockHeight)
	if err != nil {
		return false
	}
	if found {
		expiration, err := k.PendingBlocksToTimeouts.Get(ctx, blockHeight)
		if err != nil {
			return false
		}
		if currentBlockTimeNs >= expiration {
			return true
		}
	}
	return false
}

// AddUpdatePendingBlock will add a new pending block or update an existing pending block
func (k *Keeper) AddUpdatePendingBlock(ctx context.Context, pendingBlock int64, currentBlockTime time.Time) error {
	found, err := k.PendingBlocksToTimeouts.Has(ctx, pendingBlock)
	if err != nil {
		return fmt.Errorf("remove pending blocks, block %d error", pendingBlock)
	}
	if found {
		if err = k.RemovePendingBlock(ctx, pendingBlock); err != nil {
			return err
		}
	}
	expiration := currentBlockTime.Add(ResubmissionTime + RelayerPollingBuffer).UnixNano()
	if err = k.PendingBlocksToTimeouts.Set(ctx, pendingBlock, expiration); err != nil {
		return fmt.Errorf("add/update pending block, set pending block (%d) to timeout (%d)", pendingBlock, expiration)
	}
	if err = k.AddPendingBlockToTimeoutsMap(ctx, pendingBlock, expiration); err != nil {
		return fmt.Errorf("add/update pending block, add pending block to timeouts map, %v", err)
	}
	return nil
}

func (k *Keeper) AddPendingBlockToTimeoutsMap(ctx context.Context, height int64, expiration int64) error {
	found, err := k.TimeoutsToPendingBlocks.Has(ctx, expiration)
	if err != nil {
		return err
	}
	var pendingBlocks types.PendingBlocks
	if found {
		pendingBlocks, err = k.TimeoutsToPendingBlocks.Get(ctx, expiration)
		if err != nil {
			return err
		}
	}
	pendingBlocks.BlockHeights = append(pendingBlocks.BlockHeights, height)
	if err = k.TimeoutsToPendingBlocks.Set(ctx, expiration, pendingBlocks); err != nil {
		return err
	}
	return nil
}

// // RemovePendingBlock removes proven block from pending state
// This function will remove the proven block from the PendingBlocksToTimeouts map and TimeoutsToPendingBlocks map
func (k *Keeper) RemovePendingBlock(ctx context.Context, provenBlock int64) error {
	found, err := k.PendingBlocksToTimeouts.Has(ctx, provenBlock)
	if err != nil {
		return fmt.Errorf("remove pending blocks, block %d error", provenBlock)
	}
	if found {
		expiration, err := k.PendingBlocksToTimeouts.Get(ctx, provenBlock)
		if err != nil {
			return fmt.Errorf("remove pending blocks, getting pending block %d", provenBlock)
		}
		if err = k.PendingBlocksToTimeouts.Remove(ctx, provenBlock); err != nil {
			return fmt.Errorf("remove pending blocks, removing block %d", provenBlock)
		}
		pendingBlocks, err := k.TimeoutsToPendingBlocks.Get(ctx, expiration)
		if err != nil {
			return fmt.Errorf("remove pending blocks, getting expiration %d", expiration)
		}
		var newPendingBlocks []int64
		for _, blockHeight := range pendingBlocks.BlockHeights {
			if blockHeight != provenBlock {
				newPendingBlocks = append(newPendingBlocks, blockHeight)
			}
		}
		if len(newPendingBlocks) > 0 {
			pendingBlocks.BlockHeights = newPendingBlocks
			if err = k.TimeoutsToPendingBlocks.Set(ctx, expiration, pendingBlocks); err != nil {
				return fmt.Errorf("remove pending block, set new pending blocks")
			}
		} else {
			if err = k.TimeoutsToPendingBlocks.Remove(ctx, expiration); err != nil {
				return fmt.Errorf("remove pending blocks, removing timeout set %d", expiration)
			}
		}
	}
	return nil
}

// GetExpiredBlocks returns all expired blocks, proposer will propose publishing based on this set
func (k Keeper) GetExpiredBlocks(ctx context.Context, currentBlockTime time.Time) []int64 {
	currentBlockTimeNs := currentBlockTime.UnixNano()
	iterator, err := k.TimeoutsToPendingBlocks.
		Iterate(ctx, (&collections.Range[int64]{}).StartInclusive(0).EndInclusive(currentBlockTimeNs))
	if err != nil {
		return nil
	}
	defer iterator.Close()

	var expiredBlocks []int64
	for ; iterator.Valid(); iterator.Next() {
		pendingBlocks, err := iterator.Value()
		if err != nil {
			return nil
		}
		expiredBlocks = append(expiredBlocks, pendingBlocks.BlockHeights...)
	}
	return expiredBlocks
}
