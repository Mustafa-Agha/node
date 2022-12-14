package bridge

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/Mustafa-Agha/node/plugins/bridge/types"
)

func Routes(keeper Keeper) map[string]sdk.Handler {
	routes := make(map[string]sdk.Handler)
	routes[types.RouteBridge] = NewHandler(keeper)
	return routes
}
