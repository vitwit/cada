package local_test

import (
	"context"
	"fmt"
	"testing"

	abci "github.com/cometbft/cometbft/abci/types"
	"github.com/cometbft/cometbft/libs/bytes"
	coretypes "github.com/cometbft/cometbft/rpc/core/types"
	"github.com/stretchr/testify/assert"
	provider "github.com/vitwit/avail-da-module/relayer/local"
	mocks "github.com/vitwit/avail-da-module/relayer/local/mocks"
)

func TestGetBlockAtHeight_Success(t *testing.T) {
	mockRPCClient := mocks.NewRPCClient(t)
	cosmosProvider := &provider.DefaultCosmosProvider{RPCClient: mockRPCClient}
	ctx := context.Background()
	height := int64(100)

	expectedBlock := &coretypes.ResultBlock{}
	mockRPCClient.On("Block", ctx, &height).Return(expectedBlock, nil)

	block, err := cosmosProvider.GetBlockAtHeight(ctx, height)

	assert.NoError(t, err)
	assert.Equal(t, expectedBlock, block)
	mockRPCClient.AssertCalled(t, "Block", ctx, &height)
}

func TestGetBlockAtHeight_Error(t *testing.T) {
	mockRPCClient := mocks.NewRPCClient(t)
	cosmosProvider := &provider.DefaultCosmosProvider{RPCClient: mockRPCClient}
	ctx := context.Background()
	height := int64(100)

	mockRPCClient.On("Block", ctx, &height).Return(nil, fmt.Errorf("block query error"))

	block, err := cosmosProvider.GetBlockAtHeight(ctx, height)

	assert.Nil(t, block)
	assert.Error(t, err)
	mockRPCClient.AssertCalled(t, "Block", ctx, &height)
}

func TestStatus_Success(t *testing.T) {
	mockRPCClient := mocks.NewRPCClient(t)
	cosmosProvider := &provider.DefaultCosmosProvider{RPCClient: mockRPCClient}
	ctx := context.Background()

	expectedStatus := &coretypes.ResultStatus{}
	mockRPCClient.On("Status", ctx).Return(expectedStatus, nil)

	status, err := cosmosProvider.Status(ctx)

	assert.NoError(t, err)
	assert.Equal(t, expectedStatus, status)
	mockRPCClient.AssertCalled(t, "Status", ctx)
}

func TestStatus_Error(t *testing.T) {
	mockRPCClient := mocks.NewRPCClient(t)
	cosmosProvider := &provider.DefaultCosmosProvider{RPCClient: mockRPCClient}
	ctx := context.Background()

	mockRPCClient.On("Status", ctx).Return(nil, fmt.Errorf("status query error"))

	status, err := cosmosProvider.Status(ctx)

	assert.Nil(t, status)
	assert.Error(t, err)
	mockRPCClient.AssertCalled(t, "Status", ctx)
}

func TestQueryABCI_Success(t *testing.T) {
	mockRPCClient := mocks.NewRPCClient(t)
	cosmosProvider := &provider.DefaultCosmosProvider{RPCClient: mockRPCClient}
	ctx := context.Background()
	path := "/custom/path"
	data := []byte("query-data")
	hexData := bytes.HexBytes(data)

	expectedResult := &coretypes.ResultABCIQuery{
		Response: abci.ResponseQuery{
			Code:  0,
			Value: data,
		},
	}
	mockRPCClient.On("ABCIQuery", ctx, path, hexData).Return(expectedResult, nil)

	response, err := cosmosProvider.QueryABCI(ctx, path, data)

	assert.NoError(t, err)
	assert.Equal(t, expectedResult.Response.Value, response)
	mockRPCClient.AssertCalled(t, "ABCIQuery", ctx, path, hexData)
}
