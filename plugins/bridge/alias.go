package bridge

import (
	"github.com/Mustafa-Agha/node/plugins/bridge/keeper"
	"github.com/Mustafa-Agha/node/plugins/bridge/types"
)

var (
	NewKeeper = keeper.NewKeeper
)

type (
	Keeper = keeper.Keeper

	TransferOutMsg = types.TransferOutMsg
	BindMsg        = types.BindMsg
	UnbindMsg      = types.UnbindMsg
)
