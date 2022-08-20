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

	cdc.RegisterConcrete(&AppAccount{}, "cechain/Account", nil)

	cdc.RegisterConcrete(&Token{}, "cechain/Token", nil)
	cdc.RegisterConcrete(&MiniToken{}, "cechain/MiniToken", nil)
}
