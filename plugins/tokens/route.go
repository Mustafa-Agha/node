package tokens

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/bank"

	"github.com/Mustafa-Agha/node/plugins/tokens/burn"
	"github.com/Mustafa-Agha/node/plugins/tokens/freeze"
	"github.com/Mustafa-Agha/node/plugins/tokens/issue"
	"github.com/Mustafa-Agha/node/plugins/tokens/ownership"
	"github.com/Mustafa-Agha/node/plugins/tokens/seturi"
	"github.com/Mustafa-Agha/node/plugins/tokens/store"
	"github.com/Mustafa-Agha/node/plugins/tokens/swap"
	"github.com/Mustafa-Agha/node/plugins/tokens/timelock"
)

func Routes(tokenMapper store.Mapper, accKeeper auth.AccountKeeper, keeper bank.Keeper,
	timeLockKeeper timelock.Keeper, swapKeeper swap.Keeper) map[string]sdk.Handler {
	routes := make(map[string]sdk.Handler)
	routes[issue.Route] = issue.NewHandler(tokenMapper, keeper)
	routes[burn.BurnRoute] = burn.NewHandler(tokenMapper, keeper)
	routes[freeze.FreezeRoute] = freeze.NewHandler(tokenMapper, accKeeper, keeper)
	routes[timelock.MsgRoute] = timelock.NewHandler(timeLockKeeper)
	routes[swap.AtomicSwapRoute] = swap.NewHandler(swapKeeper)
	routes[seturi.SetURIRoute] = seturi.NewHandler(tokenMapper)
	routes[ownership.Route] = ownership.NewHandler(tokenMapper, keeper)
	return routes
}
