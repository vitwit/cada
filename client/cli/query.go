package cli

import (
	"log"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/spf13/cobra"
	availblob "github.com/vitwit/avail-da-module"
	"github.com/vitwit/avail-da-module/types"
)

func GetQueryCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   availblob.ModuleName,
		Short: "Querying commands for the avail-da module",
		RunE:  client.ValidateCmd,
	}

	cmd.AddCommand(GetLatestBlobStatusInfo())

	return cmd
}

func GetLatestBlobStatusInfo() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get-da-status",
		Short: "Show what range of blocks are being submitted and thier status",
		Long: `Show what range of blocks are being submitted and thier status,
		`,
		Example: "simd query cada get-da-status",
		Args:    cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}
			queryClient := types.NewQueryClient(clientCtx)

			req := &types.QuerySubmittedBlobStatusRequest{}
			res, err := queryClient.SubmittedBlobStatus(cmd.Context(), req)
			if err != nil {
				log.Fatal(err)
			}

			return clientCtx.PrintProto(res)
		},
	}
	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}
