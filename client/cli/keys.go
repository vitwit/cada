package cli

import (
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/spf13/cobra"
	types "github.com/vitwit/avail-da-module/types"
)

const (
	FlagPubKey   = "pubkey"
	FlagAddress  = "address"
	FlagLedger   = "ledger"
	FlagCoinType = "coin-type"
)

// NewKeysCmd returns a root CLI command handler for all cada keys commands.
func NewKeysCmd() *cobra.Command {
	keysCmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      types.ModuleName + " keys subcommands",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	keysCmd.AddCommand()

	return keysCmd
}
