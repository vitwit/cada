package relayer_test

import (
	"testing"

	"cosmossdk.io/log"
	tmproto "github.com/cometbft/cometbft/proto/tendermint/types"
	coretypes "github.com/cometbft/cometbft/rpc/core/types"
	tmtypes "github.com/cometbft/cometbft/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	relayer "github.com/vitwit/avail-da-module/relayer"
	localmocks "github.com/vitwit/avail-da-module/relayer/local/mocks"
)

func TestGetBlocksDataFromLocal(t *testing.T) {
	logger := log.NewNopLogger()

	mockLocalProvider := new(localmocks.CosmosProvider)

	relayer := &relayer.Relayer{
		LocalProvider: mockLocalProvider,
		Logger:        logger,
	}

	blockHeights := []int64{1, 2, 3}

	mockProtoBlock1 := &tmproto.Block{
		Header: tmproto.Header{
			Height: 1,
		},
	}
	mockProtoBlock2 := &tmproto.Block{
		Header: tmproto.Header{
			Height: 2,
		},
	}
	mockProtoBlock3 := &tmproto.Block{
		Header: tmproto.Header{
			Height: 3,
		},
	}

	expectedBlockData1, _ := mockProtoBlock1.Marshal()
	expectedBlockData2, _ := mockProtoBlock2.Marshal()
	expectedBlockData3, _ := mockProtoBlock3.Marshal()

	// Initialize expectedBlockData with the capacity needed for all three blocks
	expectedBlockData := make([]byte, 0, len(expectedBlockData1)+len(expectedBlockData2)+len(expectedBlockData3))

	// Append each block data, ensuring to reassign the result
	expectedBlockData = append(expectedBlockData, expectedBlockData1...)
	expectedBlockData = append(expectedBlockData, expectedBlockData2...)
	expectedBlockData = append(expectedBlockData, expectedBlockData3...)

	mockLocalProvider.On("GetBlockAtHeight", mock.Anything, int64(1)).Return(&coretypes.ResultBlock{
		Block: &tmtypes.Block{
			Header: tmtypes.Header{Height: 1},
		},
	}, nil)
	mockLocalProvider.On("GetBlockAtHeight", mock.Anything, int64(2)).Return(&coretypes.ResultBlock{
		Block: &tmtypes.Block{
			Header: tmtypes.Header{Height: 2},
		},
	}, nil)
	mockLocalProvider.On("GetBlockAtHeight", mock.Anything, int64(3)).Return(&coretypes.ResultBlock{
		Block: &tmtypes.Block{
			Header: tmtypes.Header{Height: 3},
		},
	}, nil)

	result := relayer.GetBlocksDataFromLocal(sdk.Context{}, blockHeights)

	assert.Equal(t, expectedBlockData, result)

	mockLocalProvider.AssertExpectations(t)
}
