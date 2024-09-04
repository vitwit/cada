package keeper

import (
	"encoding/binary"
	"errors"

	"cosmossdk.io/log"
	cometabci "github.com/cometbft/cometbft/abci/types"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type VoteExtensionHandler struct {
	logger             log.Logger
	voteExtensionCodec *codec.Codec
}

// NewVoteExtensionHandler creates a new VoteExtensionHandler.
func NewVoteExtensionHandler(logger log.Logger, codec *codec.Codec) *VoteExtensionHandler {
	return &VoteExtensionHandler{
		logger:             logger,
		voteExtensionCodec: codec,
	}
}

// ExtendVoteHandler returns a handler that includes the block number in the vote extension.
func (h *VoteExtensionHandler) ExtendVoteHandler() sdk.ExtendVoteHandler {
	return func(ctx sdk.Context, req *cometabci.RequestExtendVote) (resp *cometabci.ResponseExtendVote, err error) {
		// Get the block number (block height)
		blockNumber := uint64(ctx.BlockHeight())

		// Encode the block number into the vote extension
		extension := h.GenerateVoteExtension(blockNumber)

		// return cometabci.ResponseExtendVote{
		// 	VoteExtension: extension,
		// }
		return &cometabci.ResponseExtendVote{VoteExtension: extension}, err
	}
}

// GenerateVoteExtension encodes the block number into a byte slice.
func (h *VoteExtensionHandler) GenerateVoteExtension(blockNumber uint64) []byte {
	extension := make([]byte, 8)
	binary.BigEndian.PutUint64(extension, blockNumber)
	return extension
}

// VerifyVoteExtensionHandler returns a handler to verify the vote extension.
func (h *VoteExtensionHandler) VerifyVoteExtensionHandler() sdk.VerifyVoteExtensionHandler {
	return func(ctx sdk.Context, req *cometabci.RequestVerifyVoteExtension) (_ *cometabci.ResponseVerifyVoteExtension, err error) {

		if req == nil {
			err := errors.New("nil request")
			h.logger.Error("VerifyVoteExtensionHandler received a nil request")
			return &cometabci.ResponseVerifyVoteExtension{Status: cometabci.ResponseVerifyVoteExtension_REJECT}, err
		}

		// By default, accept empty vote extensions
		if len(req.VoteExtension) == 0 {
			h.logger.Info(
				"empty vote extension",
				"height", req.Height,
			)
			return &cometabci.ResponseVerifyVoteExtension{Status: cometabci.ResponseVerifyVoteExtension_ACCEPT}, nil
		}

		// Expected block number (block height)
		expectedBlockNumber := uint64(req.Height)

		// Verify the vote extension contains the correct block number
		err = h.VerifyVoteExtension(req.VoteExtension, expectedBlockNumber)
		if err != nil {
			return nil, err
		}

		return &cometabci.ResponseVerifyVoteExtension{Status: cometabci.ResponseVerifyVoteExtension_ACCEPT}, nil
	}
}

// VerifyVoteExtension checks if the extension matches the expected block number.
func (h *VoteExtensionHandler) VerifyVoteExtension(extension []byte, expectedBlockNumber uint64) error {
	if len(extension) != 8 {
		return errors.New("invalid extension length")
	}

	blockNumber := binary.BigEndian.Uint64(extension)
	if blockNumber != expectedBlockNumber {
		return errors.New("block number mismatch in vote extension")
	}

	return nil
}

// func (h *VoteExtensionHandler) ExtendVoteHandler() sdk.ExtendVoteHandler {
//     return func(ctx sdk.Context, req *types.RequestExtendVote) (*types.ResponseExtendVote, error) {
//         availDAHeight := h.getAvailDAHeight(ctx)

//         voteExtensionBytes, err := h.voteExtensionCodec.Encode(availDAHeight)
//         if err != nil {
//             h.logger.Printf("Failed to encode vote extension: %v", err)
//             return nil, err
//         }

//         return &types.ResponseExtendVote{
//             VoteExtension: voteExtensionBytes,
//         }, nil
//     }
// }

// // Mock function to retrieve the Avail DA height.
// func (h *VoteExtensionHandler) getAvailDAHeight(ctx sdk.Context) uint64 {
//     // This function should be implemented to retrieve the actual Avail DA height.
//     return 12345 // Example height
// }
