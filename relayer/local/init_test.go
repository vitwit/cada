package local_test

import (
	"testing"

	"github.com/stretchr/testify/suite"
	local "github.com/vitwit/avail-da-module/relayer/local"
)

type CosmosProviderTestSuite struct {
	suite.Suite
	mockRPCClient  *MockRPCClient
	cosmosProvider *local.CosmosProvider
}

func (s *CosmosProviderTestSuite) SetupTest() {
	s.mockRPCClient = new(MockRPCClient)
	s.cosmosProvider = &local.CosmosProvider{
		Cdc:       nil,
		RPCClient: s.mockRPCClient,
	}
}

func TestCosmosProviderTestSuite(t *testing.T) {
	suite.Run(t, new(CosmosProviderTestSuite))
}
