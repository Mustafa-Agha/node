package dex

import (
	"github.com/Mustafa-Agha/node/plugins/dex/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/gov"

	"github.com/Mustafa-Agha/node/plugins/dex/list"
	"github.com/Mustafa-Agha/node/plugins/dex/order"
	"github.com/Mustafa-Agha/node/plugins/tokens"
)

// Routes exports dex message routes
func Routes(dexKeeper *DexKeeper, tokenMapper tokens.Mapper, govKeeper gov.Keeper) map[string]sdk.Handler {
	routes := make(map[string]sdk.Handler)
	orderHandler := order.NewHandler(dexKeeper)
	routes[order.RouteNewOrder] = orderHandler
	routes[order.RouteCancelOrder] = orderHandler
	routes[types.ListRoute] = list.NewHandler(dexKeeper, tokenMapper, govKeeper)
	return routes
}
