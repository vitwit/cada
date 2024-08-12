package cli

import (
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/spf13/cobra"
	availblob "github.com/vitwit/avail-da-module"
)

func GetQueryCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   availblob.ModuleName,
		Short: "Querying commands for the avail-da module",
		RunE:  client.ValidateCmd,
	}

	return cmd
}
