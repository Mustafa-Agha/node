package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/tendermint/tendermint/crypto"
)

var (
	// tnt prefix address:  tnt1v8vkkymvhe2sf7gd2092ujc6hweta38xadu2pj
	// ttnt prefix address: ttnt1v8vkkymvhe2sf7gd2092ujc6hweta38xnc4wpr
	PegAccount = sdk.AccAddress(crypto.AddressHash([]byte("BinanceChainPegAccount")))
)
