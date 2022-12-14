package testutils

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/bank"

	"github.com/Mustafa-Agha/node/common/types"
	"github.com/Mustafa-Agha/node/plugins/dex/order"
	"github.com/Mustafa-Agha/node/plugins/tokens"
	"github.com/Mustafa-Agha/node/wire"
)

func MakeCodec() *wire.Codec {
	var cdc = wire.NewCodec()

	wire.RegisterCrypto(cdc) // Register crypto.
	bank.RegisterCodec(cdc)
	sdk.RegisterCodec(cdc) // Register Msgs
	tokens.RegisterWire(cdc)
	types.RegisterWire(cdc)
	cdc.RegisterConcrete(order.NewOrderMsg{}, "dex/NewOrder", nil)
	cdc.RegisterConcrete(order.CancelOrderMsg{}, "dex/CancelOrder", nil)

	cdc.RegisterConcrete(order.OrderBookSnapshot{}, "dex/OrderBookSnapshot", nil)
	cdc.RegisterConcrete(order.ActiveOrders{}, "dex/ActiveOrders", nil)

	return cdc
}
