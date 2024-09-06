package cli

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"

	availblob "github.com/vitwit/avail-da-module"
	"github.com/vitwit/avail-da-module/keeper"
	"github.com/vitwit/avail-da-module/types"
)

// NewTxCmd creates and returns a Cobra command for transaction subcommands related to the availblob module.
func NewTxCmd(keeper *keeper.Keeper) *cobra.Command {
	txCmd := &cobra.Command{
		Use:                        availblob.ModuleName,
		Short:                      availblob.ModuleName + " transaction subcommands",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	// keeper.ClientCmd = txCmd

	txCmd.AddCommand(NewSubmitBlobCmd())
	txCmd.AddCommand(NewUpdateBlobStatusCmd(), InitKepperClientCmd(keeper))

	return txCmd
}

// InitKepperClientCmd creates and returns a Cobra command to initialize the keeper client.
func InitKepperClientCmd(keeper *keeper.Keeper) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "init-keeper-client", // TODO: remove this subcommand
		Short:   "initlialize a client to use in keeper",
		Example: "init-keeper-client",
		Args:    cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Println("setting keeper client.....", keeper.ClientCmd == nil)
			keeper.ClientCmd = cmd
			fmt.Println("setting keeper client.....", keeper.ClientCmd == nil)
			return nil
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}

// NewSubmitBlobCmd creates and returns a Cobra command to submit a blob with a specified range of blocks.
func NewSubmitBlobCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "submit-blob [from_block] [to_block]",
		Short:   "request to submit blob with blocks from [from] to [to]",
		Example: "submit-blob 11 15",
		Args:    cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			fromBlock, err := strconv.Atoi(args[0])
			if err != nil {
				return err
			}
			toBlock, err := strconv.Atoi(args[1])
			if err != nil {
				return err
			}

			msg := types.MsgSubmitBlobRequest{
				BlocksRange: &types.Range{
					From: uint64(fromBlock),
					To:   uint64(toBlock),
				},
				ValidatorAddress: clientCtx.GetFromAddress().String(),
			}

			// cmd.Marsha
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), &msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}

// NewUpdateBlobStatusCmd creates and returns a Cobra command to update the status of a blob within a specified range of blocks.
func NewUpdateBlobStatusCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use: "update-blob [from_block] [to_block] [status] [avail_height]",
		Short: `update blob status by giving blocks range and status(success|failure)
		 and the avail height at which the blob is stored`,
		Example: "update-blob 11 15 success 120",
		Args:    cobra.ExactArgs(4),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			fromBlock, err := strconv.Atoi(args[0])
			if err != nil {
				return err
			}
			toBlock, err := strconv.Atoi(args[1])
			if err != nil {
				return err
			}

			isSuccess, err := ParseStatus(args[2])
			if err != nil {
				return err
			}

			availHeight, err := strconv.Atoi(args[3])
			if err != nil {
				return err
			}

			msg := types.MsgUpdateBlobStatusRequest{
				BlocksRange: &types.Range{
					From: uint64(fromBlock),
					To:   uint64(toBlock),
				},
				ValidatorAddress: clientCtx.GetFromAddress().String(),
				AvailHeight:      uint64(availHeight),
				IsSuccess:        isSuccess,
			}
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), &msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}

// ParseStatus converts a string status to a boolean value indicating success or failure.
func ParseStatus(status string) (bool, error) {
	status = strings.ToUpper(status)
	if status == "SUCCESS" {
		return true, nil
	}

	if status == "FAILURE" {
		return false, nil
	}

	return false, errors.New("invalid status")
}
