package relayer

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/spf13/cobra"
	"github.com/vitwit/avail-da-module/types"
)

func (r *Relayer) StartBlobLifeCycle(msg types.MsgSubmitBlobRequest, cmd *cobra.Command) { // TODO: deprecate

	fmt.Println("inside like cycle.................", msg, cmd == nil)
	clientCtx, err := client.GetClientTxContext(cmd)
	if err != nil {
		fmt.Println("error in start blob life cycle", err)
		return
	}

	msg.ValidatorAddress = clientCtx.GetFromAddress().String()
	err = tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), &msg)
	if err != nil {
		fmt.Println("error in start blob life cycle", err)
		return
	}
	fmt.Println("broadcast success..........")
}
