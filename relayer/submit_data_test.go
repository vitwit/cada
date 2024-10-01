package relayer_test

import (
	"testing"

	"cosmossdk.io/log"
	"github.com/stretchr/testify/assert"
	relayer "github.com/vitwit/avail-da-module/relayer"
	"github.com/vitwit/avail-da-module/relayer/avail"
	mocks "github.com/vitwit/avail-da-module/relayer/avail/mocks"
	cadatypes "github.com/vitwit/avail-da-module/x/cada/types"
)

func TestSubmitDataToAvailClient(t *testing.T) {
	logger := log.NewNopLogger()

	mockDAClient := new(mocks.DA)

	relayer := &relayer.Relayer{
		AvailDAClient: mockDAClient,
		Logger:        logger,
		AvailConfig:   cadatypes.AvailConfiguration{AppID: 1},
	}

	data := []byte("test data")
	blocks := []int64{1, 2, 3}

	expectedBlockInfo := avail.BlockMetaData{
		BlockHash:   "hash123",
		BlockNumber: 42,
		Hash:        "hash456",
	}

	t.Run("success", func(t *testing.T) {
		mockDAClient.On("Submit", data).Return(expectedBlockInfo, nil)

		blockInfo, err := relayer.SubmitDataToAvailClient(data, blocks)

		assert.NoError(t, err)
		assert.Equal(t, expectedBlockInfo, blockInfo)

		mockDAClient.AssertExpectations(t)
	})
}
