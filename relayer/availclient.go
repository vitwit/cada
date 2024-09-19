package relayer

import "github.com/vitwit/avail-da-module/types"

// AvailClient is the client that handles data submission
type AvailClient struct {
	config types.AvailConfiguration
}

// NewAvailClient initializes a new AvailClient
func NewAvailClient(config types.AvailConfiguration) (*AvailClient, error) {
	// api, err := gsrpc.NewSubstrateAPI(config.AppRpcURL)
	// if err != nil {
	// 	return nil, fmt.Errorf("cannot create api:%w", err)
	// }

	return &AvailClient{config: config}, nil
}
