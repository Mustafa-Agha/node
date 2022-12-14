package commands

import (
	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client"

	"github.com/Mustafa-Agha/node/wire"
)

const (
	flagSymbol = "symbol"
	flagAmount = "amount"
)

func AddCommands(cmd *cobra.Command, cdc *wire.Codec) {

	tokenCmd := &cobra.Command{
		Use:   "token",
		Short: "issue or view tokens",
		Long:  ``,
	}

	cmdr := Commander{Cdc: cdc}
	tokenCmd.AddCommand(
		client.PostCommands(
			issueTokenCmd(cmdr),
			mintTokenCmd(cmdr),
			burnTokenCmd(cmdr),
			freezeTokenCmd(cmdr),
			unfreezeTokenCmd(cmdr),
			timeLockCmd(cmdr),
			timeUnlockCmd(cmdr),
			timeRelockCmd(cmdr),
			initiateHTLTCmd(cmdr),
			depositHTLTCmd(cmdr),
			claimHTLTCmd(cmdr),
			refundHTLTCmd(cmdr),
			transferOwnershipCmd(cmdr),
		)...)

	tokenCmd.AddCommand(
		client.GetCommands(
			listTokensCmd,
			getTokenInfoCmd(cmdr),
			queryTimeLocksCmd(cmdr),
			queryTimeLockCmd(cmdr),
			querySwapCmd(cmdr),
			querySwapsByRecipientCmd(cmdr),
			querySwapsByCreatorCmd(cmdr))...)

	tokenCmd.AddCommand(
		client.PostCommands(MultiSendCmd(cdc))...,
	)

	tokenCmd.AddCommand(
		client.PostCommands(
			issueMiniTokenCmd(cmdr),
			issueTinyTokenCmd(cmdr),
			setTokenURICmd(cmdr))...,
	)

	tokenCmd.AddCommand(client.LineBreak)

	cmd.AddCommand(tokenCmd)

}
