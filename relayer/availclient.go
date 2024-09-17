package relayer

// AvailClient is the client that handles data submission
type AvailClient struct {
	config AvailConfig
}

// NewAvailClient initializes a new AvailClient
func NewAvailClient(config AvailConfig) (*AvailClient, error) {
	// api, err := gsrpc.NewSubstrateAPI(config.AppRpcURL)
	// if err != nil {
	// 	return nil, fmt.Errorf("cannot create api:%w", err)
	// }

	return &AvailClient{config: config}, nil
}
