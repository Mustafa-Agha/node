package account

import (
	"github.com/Mustafa-Agha/node/wire"
)

// Register concrete types on wire codec
func RegisterWire(cdc *wire.Codec) {
	cdc.RegisterConcrete(SetAccountFlagsMsg{}, "scripts/SetAccountFlagsMsg", nil)
}
