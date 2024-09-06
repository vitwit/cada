package chainclient

import (
	"fmt"
	"log"
	"os"

	cometrpc "github.com/cometbft/cometbft/rpc/client/http"
	clitx "github.com/cosmos/cosmos-sdk/client/tx"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/vitwit/avail-da-module/types"

	"path/filepath"

	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/crypto/keyring"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
)

// GetBinPath returns the path to the Avail SDK home directory within the user's home directory.
func GetBinPath() string {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.Fatal(err)
	}

	availdHomePath := filepath.Join(homeDir, ".availsdk")
	fmt.Println("availdHonmePath.......", availdHomePath)
	return availdHomePath
}

// ExecuteTX handles the creation and submission of a transaction to update blob status on the chain.
// It uses keyring and RPC client configurations to interact with the network.
func ExecuteTX(ctx sdk.Context, msg types.MsgUpdateBlobStatusRequest, cdc codec.BinaryCodec) error {
	// Define keyring and RPC client configuration

	// homePath := "/home/vitwit/.availsdk"
	homePath := GetBinPath()
	key := os.Getenv("KEY")
	fmt.Println("get key namee.........", key)
	if key == "" { //TODO : remove this later
		key = "alice"
	}
	keyName := key
	rpcAddress := "http://localhost:26657"

	// Create a keyring
	kr, err := keyring.New(sdk.KeyringServiceName(), keyring.BackendTest, homePath, os.Stdin, cdc.(codec.Codec))
	if err != nil {
		return fmt.Errorf("error creating keyring: %w", err)
	}

	info, err := kr.Key(keyName)
	valAddr, err := info.GetAddress()
	fmt.Println("after address................", valAddr)

	// Create an RPC client
	rpcClient, err := cometrpc.NewWithTimeout(rpcAddress, "/websocket", 3)
	if err != nil {
		return fmt.Errorf("error creating RPC client: %w", err)
	}

	// Create a new client context
	clientCtx := NewClientCtx(kr, rpcClient, ctx.ChainID(), cdc, homePath, valAddr)

	// Set the client context's from fields
	clientCtx.FromName = keyName
	clientCtx.FromAddress = valAddr

	// Fetch account number and sequence from the blockchain
	accountRetriever := authtypes.AccountRetriever{}
	account, err := accountRetriever.GetAccount(clientCtx, valAddr)
	if err != nil {
		return fmt.Errorf("error retrieving account: %w", err)
	}

	fmt.Println("account details......", account.GetAccountNumber(), account.GetSequence())

	// Set the correct account number and sequence
	factory := NewFactory(clientCtx).
		WithAccountNumber(account.GetAccountNumber()).
		WithSequence(account.GetSequence())

	// Create a transaction factory and set the validator address in the message
	// factory := NewFactory(clientCtx)
	msg.ValidatorAddress = valAddr.String()

	// Generate and broadcast the transaction
	if err := clitx.GenerateOrBroadcastTxWithFactory(clientCtx, factory, &msg); err != nil {
		return fmt.Errorf("error broadcasting transaction: %w", err)
	}

	return nil
}
