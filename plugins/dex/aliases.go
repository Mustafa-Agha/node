package dex

import (
	"github.com/Mustafa-Agha/node/plugins/dex/order"
	"github.com/Mustafa-Agha/node/plugins/dex/store"
	"github.com/Mustafa-Agha/node/plugins/dex/types"
)

// type MsgList = list.Msg
// type TradingPair = types.TradingPair

type TradingPairMapper = store.TradingPairMapper
type DexKeeper = order.DexKeeper

var NewTradingPairMapper = store.NewTradingPairMapper
var NewDexKeeper = order.NewDexKeeper

const DefaultCodespace = types.DefaultCodespace
