package cli

import (
	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client"

	availblob "github.com/vitwit/avail-da-module"
)

// NewTxCmd
func NewTxCmd() *cobra.Command {
	txCmd := &cobra.Command{
		Use:                        availblob.ModuleName,
		Short:                      availblob.ModuleName + " transaction subcommands",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	return txCmd
}
