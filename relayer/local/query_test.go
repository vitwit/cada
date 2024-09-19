package local_test

import (
	"context"
	"fmt"

	"github.com/cometbft/cometbft/libs/bytes"
	coretypes "github.com/cometbft/cometbft/rpc/core/types"
	"github.com/stretchr/testify/mock"
)

type MockRPCClient struct {
	mock.Mock
}

func (m *MockRPCClient) Block(ctx context.Context, height *int64) (*coretypes.ResultBlock, error) {
	args := m.Called(ctx, height)

	if block, ok := args.Get(0).(*coretypes.ResultBlock); ok {
		return block, args.Error(1)
	}

	return nil, args.Error(1)
}

func (m *MockRPCClient) Status(ctx context.Context) (*coretypes.ResultStatus, error) {
	args := m.Called(ctx)

	if status, ok := args.Get(0).(*coretypes.ResultStatus); ok {
		return status, args.Error(1)
	}

	return nil, args.Error(1)
}

func (m *MockRPCClient) ABCIQuery(ctx context.Context, path string, data bytes.HexBytes) (*coretypes.ResultABCIQuery, error) {
	args := m.Called(ctx, path, data)
	return args.Get(0).(*coretypes.ResultABCIQuery), args.Error(1)
}

func (s *CosmosProviderTestSuite) TestQueryABCI_Success() {
	path := "/store"
	data := bytes.HexBytes("queryData")

	result := &coretypes.ResultABCIQuery{}

	s.mockRPCClient.On("ABCIQuery", mock.Anything, path, data).Return(result, nil)

	resp, err := s.cosmosProvider.QueryABCI(context.Background(), path, data)

	s.Require().NoError(err)
	s.Require().NotNil(resp)
}

func (s *CosmosProviderTestSuite) TestGetBlockAtHeight_Success() {
	height := int64(10)
	expectedBlock := &coretypes.ResultBlock{}
	s.mockRPCClient.On("Block", mock.Anything, &height).Return(expectedBlock, nil)

	block, err := s.cosmosProvider.GetBlockAtHeight(context.Background(), height)

	s.Require().NoError(err)
	s.Require().Equal(expectedBlock, block)
}

func (s *CosmosProviderTestSuite) TestGetBlockAtHeight_Error() {
	height := int64(10)
	expectedError := fmt.Errorf("error querying block")
	s.mockRPCClient.On("Block", mock.Anything, &height).Return(nil, expectedError)

	block, err := s.cosmosProvider.GetBlockAtHeight(context.Background(), height)

	s.Require().Error(err)
	s.Require().Nil(block)
}

func (s *CosmosProviderTestSuite) TestStatus_Success() {
	expectedStatus := &coretypes.ResultStatus{
		// Populate with expected status data
	}
	s.mockRPCClient.On("Status", mock.Anything).Return(expectedStatus, nil)

	status, err := s.cosmosProvider.Status(context.Background())

	s.Require().NoError(err)
	s.Require().Equal(expectedStatus, status)
}

func (s *CosmosProviderTestSuite) TestStatus_Error() {
	expectedError := fmt.Errorf("error querying status")
	s.mockRPCClient.On("Status", mock.Anything).Return(nil, expectedError)

	status, err := s.cosmosProvider.Status(context.Background())

	s.Require().Error(err)
	s.Require().Nil(status)
	s.Require().Equal("error querying status: error querying status", err.Error())
}
