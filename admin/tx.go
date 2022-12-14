package admin

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/bank"

	"github.com/Mustafa-Agha/node/common/runtime"
	"github.com/Mustafa-Agha/node/plugins/account"
	"github.com/Mustafa-Agha/node/plugins/bridge"
	"github.com/Mustafa-Agha/node/plugins/dex/order"
	list "github.com/Mustafa-Agha/node/plugins/dex/types"
	"github.com/Mustafa-Agha/node/plugins/tokens/burn"
	"github.com/Mustafa-Agha/node/plugins/tokens/freeze"
	"github.com/Mustafa-Agha/node/plugins/tokens/issue"
	"github.com/Mustafa-Agha/node/plugins/tokens/ownership"
	"github.com/Mustafa-Agha/node/plugins/tokens/seturi"
	"github.com/Mustafa-Agha/node/plugins/tokens/swap"
	"github.com/Mustafa-Agha/node/plugins/tokens/timelock"
)

var transferOnlyModeBlackList = []string{
	burn.BurnMsg{}.Type(),
	freeze.FreezeMsg{}.Type(),
	freeze.UnfreezeMsg{}.Type(),
	issue.IssueMsg{}.Type(),
	issue.MintMsg{}.Type(),
	order.NewOrderMsg{}.Type(),
	order.CancelOrderMsg{}.Type(),
	timelock.TimeLockMsg{}.Type(),
	timelock.TimeUnlockMsg{}.Type(),
	timelock.TimeRelockMsg{}.Type(),
	issue.IssueMiniMsg{}.Type(),
	issue.IssueTinyMsg{}.Type(),
	seturi.SetURIMsg{}.Type(),
	list.ListMsg{}.Type(),
	list.ListMiniMsg{}.Type(),
	ownership.TransferOwnershipMsg{}.Type(),
	swap.HTLTMsg{}.Type(),
	swap.DepositHTLTMsg{}.Type(),
	swap.ClaimHTLTMsg{}.Type(),
	swap.RefundHTLTMsg{}.Type(),
	account.SetAccountFlagsMsg{}.Type(),
	bridge.BindMsg{}.Type(),
	bridge.UnbindMsg{}.Type(),
	bridge.TransferOutMsg{}.Type(),
}

var TxBlackList = map[runtime.Mode][]string{
	runtime.TransferOnlyMode: transferOnlyModeBlackList,
	runtime.RecoverOnlyMode:  append(transferOnlyModeBlackList, bank.MsgSend{}.Type()),
}

func TxNotAllowedError() sdk.Error {
	return sdk.ErrInternal(fmt.Sprintf("The tx is not allowed, RunningMode: %v", runtime.GetRunningMode()))
}

func IsTxAllowed(tx sdk.Tx) bool {
	mode := runtime.GetRunningMode()
	if mode == runtime.NormalMode {
		return true
	}

	for _, msg := range tx.GetMsgs() {
		if !isMsgAllowed(msg, mode) {
			return false
		}
	}
	return true
}

func isMsgAllowed(msg sdk.Msg, mode runtime.Mode) bool {
	for _, msgType := range TxBlackList[mode] {
		if msgType == msg.Type() {
			return false
		}
	}

	return true
}
