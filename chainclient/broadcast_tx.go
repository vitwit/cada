package chainclient

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	cometrpc "github.com/cometbft/cometbft/rpc/client/http"
	clitx "github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/crypto/keyring"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/vitwit/avail-da-module/types"
)

func GetBinPath() string {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.Fatal(err)
	}

	availdHomePath := filepath.Join(homeDir, ".availsdk")
	fmt.Println("availdHonmePath.......", availdHomePath)
	return availdHomePath
}

func ExecuteTX(ctx sdk.Context, msg types.MsgUpdateBlobStatusRequest, cdc codec.BinaryCodec, rpcAddr string) error {
	// Define keyring and RPC client configuration

	// homePath := "/home/vitwit/.availsdk"
	homePath := GetBinPath()
	key := os.Getenv("KEY")
	keyName := key

	if key == "" {
		keyName = "alice"
	}

	// Create a keyring
	kr, err := keyring.New(sdk.KeyringServiceName(), keyring.BackendTest, homePath, os.Stdin, cdc.(codec.Codec))
	if err != nil {
		return fmt.Errorf("error creating keyring: %w", err)
	}

	info, err := kr.Key(keyName)
	if err != nil {
		return err
	}
	valAddr, err := info.GetAddress()
	if err != nil {
		return fmt.Errorf("error while getting account address : %w", err)
	}

	// Create an RPC client
	rpcClient, err := cometrpc.NewWithTimeout(rpcAddr, "/websocket", 3)
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
