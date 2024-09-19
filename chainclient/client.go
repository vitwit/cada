package chainclient

const (
	KeyringBackendTest = "test"
)

// ChainClient represents a client used for interacting with a network.
// It contains configurations and credentials necessary for broadcasting transactions and interacting with the chain.
type ChainClient struct {
	Address            string `json:"address"`
	AddressPrefix      string `json:"account_address_prefix"`
	RPC                string `json:"rpc"`
	Key                string `json:"key"`
	Mnemonic           string `json:"mnemonic"`
	KeyringServiceName string `json:"keyring_service_name"`
	HDPath             string `json:"hd_path"`
	Enabled            bool   `json:"enabled"`
	ChainName          string `json:"chain_name"`
	Denom              string `json:"denom"`
}
