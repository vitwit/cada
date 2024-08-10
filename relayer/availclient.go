package relayer

import (
	"fmt"

	gsrpc "github.com/centrifuge/go-substrate-rpc-client/v4"
)

// AvailClient is the client that handles data submission
type AvailClient struct {
	api    *gsrpc.SubstrateAPI
	config AvailConfig
}

// NewAvailClient initializes a new AvailClient
func NewAvailClient(config AvailConfig) (*AvailClient, error) {
	api, err := gsrpc.NewSubstrateAPI(config.AppRpcURL)
	if err != nil {
		return nil, fmt.Errorf("cannot create api:%w", err)
	}

	return &AvailClient{api: api, config: config}, nil
}
