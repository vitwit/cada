package cli

import (
	"log"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/spf13/cobra"
	types "github.com/vitwit/avail-da-module/types"
)

// GetQueryCmd returns the root query command for the avail-da module.
func GetQueryCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   types.ModuleName,
		Short: "Querying commands for the avail-da module",
		RunE:  client.ValidateCmd,
	}

	cmd.AddCommand(GetLatestBlobStatusInfo())

	return cmd
}

// GetLatestBlobStatusInfo returns a command to query the latest status of blob submissions.
func GetLatestBlobStatusInfo() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get-da-status",
		Short: "Shows what range of blocks are being submitted and their status",
		Long: `Shows what range of blocks are being submitted and their status,
		`,
		Example: "simd query cada get-da-status",
		Args:    cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, _ []string) error {
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
