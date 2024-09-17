package chainclient

import (
	"bytes"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/tx"
)

const (
	KeyringBackendTest = "test"
)

// ChainClient is client to interact with SPN.
type ChainClient struct {
	factory            tx.Factory
	clientCtx          client.Context
	out                *bytes.Buffer
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