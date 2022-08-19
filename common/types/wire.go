package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/Mustafa-Agha/node/wire"
)

func RegisterWire(cdc *wire.Codec) {
	// Register AppAccount
	cdc.RegisterInterface((*sdk.Account)(nil), nil)
	cdc.RegisterInterface((*NamedAccount)(nil), nil)
	cdc.RegisterInterface((*IToken)(nil), nil)

	cdc.RegisterConcrete(&AppAccount{}, "tntchain/Account", nil)

	cdc.RegisterConcrete(&Token{}, "tntchain/Token", nil)
	cdc.RegisterConcrete(&MiniToken{}, "tntchain/MiniToken", nil)
}
