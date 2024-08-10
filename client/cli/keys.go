package cli

import (
	"fmt"
	"path/filepath"

	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/server"
	"github.com/cosmos/cosmos-sdk/version"

	availblob "github.com/vitwit/avail-da-module"
	"github.com/vitwit/avail-da-module/relayer"
)

// NewKeysCmd returns a root CLI command handler for all x/tiablob keys commands.
func NewKeysCmd() *cobra.Command {
	keysCmd := &cobra.Command{
		Use:                        availblob.ModuleName,
		Short:                      availblob.ModuleName + " keys subcommands",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	keysCmd.AddCommand(
		NewKeysAddCmd(),
		NewKeysRestoreCmd(),
		// NewKeysShowCmd(),
		// NewKeysDeleteCmd(),
	)

	return keysCmd
}

// NewKeysAddCmd returns a CLI command handler for creating a key for posting blocks to Celestia.
func NewKeysAddCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "add",
		Short: "Generate a key for posting blocks or feegranting on avail",
		Args:  cobra.NoArgs,
		Example: fmt.Sprintf(` 
$ %s keys availblob add
`, version.AppName),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			serverCtx := server.GetServerContextFromCmd(cmd)
			cfg := relayer.AvailConfigFromAppOpts(serverCtx.Viper)
			_ = cfg

			keyDir := filepath.Join(clientCtx.HomeDir, "keys")

			_ = keyDir

			return nil
		},
	}
	return cmd
}

// NewKeysRestoreCmd returns a CLI command handler for restoring a key for posting blocks to Celestia.
func NewKeysRestoreCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "restore [mnemonic]",
		Short: "Generate a key for posting blocks to Celestia",
		Args:  cobra.ExactArgs(1),
		Example: fmt.Sprintf(` 
$ %s keys availblob restore "pattern match caution ..."
`, version.AppName),
		RunE: func(cmd *cobra.Command, args []string) error {

			return nil
		},
	}

	return cmd
}
