package account

import (
	"github.com/cosmos/cosmos-sdk/x/auth"

	app "github.com/Mustafa-Agha/node/common/types"
	"github.com/Mustafa-Agha/node/plugins/account/scripts"
)

func InitPlugin(appp app.ChainApp, accountKeeper auth.AccountKeeper) {
	// add msg handlers
	for route, handler := range routes(accountKeeper) {
		appp.GetRouter().AddRoute(route, handler)
	}

	//register transfer memo checker
	scripts.RegisterTransferMemoCheckScript(accountKeeper)
}
