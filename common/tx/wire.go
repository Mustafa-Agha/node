package tx

import (
	"github.com/cosmos/cosmos-sdk/x/auth"

	"github.com/Mustafa-Agha/node/wire"
)

func RegisterWire(cdc *wire.Codec) {
	cdc.RegisterConcrete(&auth.StdTx{}, "auth/StdTx", nil)
}
