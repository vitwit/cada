package cli

import (
	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client"

	availblob "github.com/vitwit/avail-da-module"
)

const (
	FlagPubKey   = "pubkey"
	FlagAddress  = "address"
	FlagLedger   = "ledger"
	FlagCoinType = "coin-type"
)

// NewKeysCmd returns a root CLI command handler for all x/availblob keys commands.
func NewKeysCmd() *cobra.Command {
	keysCmd := &cobra.Command{
		Use:                        availblob.ModuleName,
		Short:                      availblob.ModuleName + " keys subcommands",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	keysCmd.AddCommand(
	// NewKeysAddCmd(),
	// NewKeysRestoreCmd(),
	// NewKeysShowCmd(),
	// NewKeysDeleteCmd(),
	)

	return keysCmd
}
