package cli

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/cosmos/cosmos-sdk/version"

	availblob1 "github.com/vitwit/avail-da-module"
	"github.com/vitwit/avail-da-module/types"
)

// NewTxCmd returns a root CLI command handler for all x/availblob transaction commands.
func NewTxCmd() *cobra.Command {
	txCmd := &cobra.Command{
		Use:                        availblob1.ModuleName,
		Short:                      availblob1.ModuleName + " transaction subcommands",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	txCmd.AddCommand(
		NewSetAvailAddrCmd(),
	)

	return txCmd
}

// NewSetCelestiaAddrCmd returns a CLI command handler for creating a MsgSetCelestiaAddress transaction.
func NewSetAvailAddrCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "set-avail-address [celestia-addr] [celesita-val-address]",
		Short: "Set the Avail address & validator addresses for the validator",
		Args:  cobra.ExactArgs(2),
		Long:  `Set the Avail address and validator address of the validator for feegranting on Avail.`,
		Example: fmt.Sprintf(`
$ %s tx availblob set-avai-address availjzv52ewect8ntvwjs2za087yzl6y3smf5etf3n 5ECe3ANZA9HaxYexsV8yRGZXhzrTs68ScoYcYHwfLTQmzyki  --from keyname
`, version.AppName),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			txf, err := tx.NewFactoryCLI(clientCtx, cmd.Flags())
			if err != nil {
				return err
			}

			msg := &types.MsgSetAvailAddress{
				ValidatorAddress: args[1],
				AvailAddress:     args[0],
			}

			return tx.GenerateOrBroadcastTxWithFactory(clientCtx, txf, msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	_ = cmd.MarkFlagRequired(flags.FlagFrom)

	return cmd
}
