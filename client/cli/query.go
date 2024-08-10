package cli

import (
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/spf13/cobra"
	availblob1 "github.com/vitwit/avail-da-module"
	"github.com/vitwit/avail-da-module/types"
)

func GetQueryCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   availblob1.ModuleName,
		Short: "Querying commands for the availblob module",
		RunE:  client.ValidateCmd,
	}

	cmd.AddCommand(
		GetValidatorsInfo(),
	)

	return cmd
}

func GetValidatorsInfo() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get-validators",
		Short: "Get registered Celestia validator addresses",
		Long: `Get registered Celestia validator addresses,
		`,
		Args: cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}
			queryClient := types.NewQueryClient(clientCtx)

			req := &types.QueryValidatorsRequest{}
			res, _ := queryClient.Validators(cmd.Context(), req)

			return clientCtx.PrintProto(res)
		},
	}
	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}
