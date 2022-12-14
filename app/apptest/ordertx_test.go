package apptest

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	abci "github.com/tendermint/tendermint/abci/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/fees"
	"github.com/cosmos/cosmos-sdk/x/auth"

	"github.com/Mustafa-Agha/node/common/utils"
	o "github.com/Mustafa-Agha/node/plugins/dex/order"
	"github.com/Mustafa-Agha/node/plugins/dex/types"
)

type level struct {
	price utils.Fixed8
	qty   utils.Fixed8
}

func getOrderBook(pair string) ([]level, []level, bool) {
	buys := make([]level, 0)
	sells := make([]level, 0)
	orderbooks, pendingMatch := testApp.DexKeeper.GetOrderBookLevels(pair, 5)
	for _, l := range orderbooks {
		if l.BuyPrice != 0 {
			buys = append(buys, level{price: l.BuyPrice, qty: l.BuyQty})
		}
		if l.SellPrice != 0 {
			sells = append(sells, level{price: l.SellPrice, qty: l.SellQty})
		}
	}
	return buys, sells, pendingMatch
}

func genOrderID(add sdk.AccAddress, seq int64, ctx sdk.Context, am auth.AccountKeeper) string {
	acc := am.GetAccount(ctx, add)
	if acc.GetSequence() != seq {
		err := acc.SetSequence(seq)
		if err != nil {
			panic(err)
		}
		am.SetAccount(ctx, acc)
	}
	oid := fmt.Sprintf("%X-%d", add, seq)
	return oid
}

func newTestFeeConfig() o.FeeConfig {
	feeConfig := o.NewFeeConfig()
	feeConfig.FeeRateNative = 500
	feeConfig.FeeRate = 1000
	feeConfig.ExpireFeeNative = 2e4
	feeConfig.ExpireFee = 1e5
	feeConfig.IOCExpireFeeNative = 1e4
	feeConfig.IOCExpireFee = 5e4
	feeConfig.CancelFeeNative = 2e4
	feeConfig.CancelFee = 1e5
	return feeConfig
}

func Test_handleNewOrder_CheckTx(t *testing.T) {
	assert := assert.New(t)
	ctx := testApp.NewContext(sdk.RunTxModeCheck, abci.Header{})
	InitAccounts(ctx, testApp)
	testApp.DexKeeper.PairMapper.AddTradingPair(ctx, types.NewTradingPair("BTC-000", "CE", 1e8))

	am := testApp.AccountKeeper
	acc := Account(0)
	acc2 := Account(1)
	add := acc.GetAddress()
	add2 := acc2.GetAddress()
	msg := o.NewNewOrderMsg(add, genOrderID(add, 0, ctx, am), 1, "BTC-000_CE", 355e8, 100e8)
	res, e := testClient.CheckTxSync(msg, testApp.Codec)
	assert.NotEqual(uint32(0), res.Code)
	assert.Nil(e)
	assert.Regexp(".*do not have enough token to lock.*", res.GetLog())
	assert.Equal(int64(500e8), GetAvail(ctx, add, "CE"))
	assert.Equal(int64(0), GetLocked(ctx, add, "CE"))
	assert.Equal(int64(200e8), GetAvail(ctx, add, "BTC-000"))
	assert.Equal(int64(0), GetLocked(ctx, add, "BTC-000"))

	msg = o.NewNewOrderMsg(add, genOrderID(add, 0, ctx, am), 1, "BTC-000_CE", 355e8, 1e8)
	res, e = testClient.CheckTxSync(msg, testApp.Codec)
	assert.Equal(uint32(0), res.Code)
	assert.Nil(e)
	assert.Equal(int64(145e8), GetAvail(ctx, add, "CE"))
	assert.Equal(int64(355e8), GetLocked(ctx, add, "CE"))
	assert.Equal(int64(200e8), GetAvail(ctx, add, "BTC-000"))
	assert.Equal(int64(0), GetLocked(ctx, add, "BTC-000"))

	// using acc2

	msg = o.NewNewOrderMsg(add2, genOrderID(add2, 0, ctx, am), 2, "BTC-000_CE", 355e8, 250e8)
	res, e = testClient.CheckTxSync(msg, testApp.Codec)
	assert.NotEqual(uint32(0), res.Code)
	assert.Nil(e)
	assert.Regexp(".*do not have enough token to lock.*", res.GetLog())
	assert.Equal(int64(500e8), GetAvail(ctx, add2, "CE"))
	assert.Equal(int64(0), GetLocked(ctx, add2, "CE"))
	assert.Equal(int64(200e8), GetAvail(ctx, add2, "BTC-000"))
	assert.Equal(int64(0), GetLocked(ctx, add2, "BTC-000"))

	msg = o.NewNewOrderMsg(add2, genOrderID(add2, 0, ctx, am), 2, "BTC-000_CE", 355e8, 200e8)
	res, e = testClient.CheckTxSync(msg, testApp.Codec)
	assert.Equal(uint32(0), res.Code)
	assert.Nil(e)
	assert.Equal(int64(500e8), GetAvail(ctx, add2, "CE"))
	assert.Equal(int64(0), GetLocked(ctx, add2, "CE"))
	assert.Equal(int64(0), GetAvail(ctx, add2, "BTC-000"))
	assert.Equal(int64(200e8), GetLocked(ctx, add2, "BTC-000"))
}

func Test_handleNewOrder_DeliverTx(t *testing.T) {
	assert := assert.New(t)
	testClient.cl.BeginBlockSync(abci.RequestBeginBlock{})
	ctx := testApp.NewContext(sdk.RunTxModeDeliver, abci.Header{})
	InitAccounts(ctx, testApp)
	testApp.DexKeeper.ClearOrderBook("BTC-000_CE")
	tradingPair := types.NewTradingPair("BTC-000", "CE", 1e8)
	testApp.DexKeeper.PairMapper.AddTradingPair(ctx, tradingPair)
	testApp.DexKeeper.AddEngine(tradingPair)
	testApp.DexKeeper.GetEngines()["BTC-000_CE"].LastMatchHeight = -1

	tradingPair2 := types.NewTradingPair("ETH-001", "CE", 1e8)
	testApp.DexKeeper.PairMapper.AddTradingPair(ctx, tradingPair2)
	testApp.DexKeeper.AddEngine(tradingPair2)
	testApp.DexKeeper.GetEngines()["ETH-001_CE"].LastMatchHeight = -1

	add := Account(0).GetAddress()
	oid := fmt.Sprintf("%X-0", add)
	msg := o.NewNewOrderMsg(add, oid, 1, "BTC-000_CE", 355e8, 1e8)

	res, e := testClient.DeliverTxSync(msg, testApp.Codec)
	t.Logf("res is %v and error is %v", res, e)
	assert.Equal(uint32(0), res.Code)
	assert.Nil(e)
	buys, sells, pendingMatch := getOrderBook("BTC-000_CE")
	assert.Equal(1, len(buys))
	assert.Equal(0, len(sells))
	assert.Equal(true, pendingMatch)
	assert.Equal(utils.Fixed8(355e8), buys[0].price)
	assert.Equal(utils.Fixed8(1e8), buys[0].qty)

	buys, sells, pendingMatch = getOrderBook("ETH-001_CE")
	assert.Equal(0, len(buys))
	assert.Equal(0, len(sells))
	assert.Equal(false, pendingMatch)
}

func Test_Match(t *testing.T) {
	assert := assert.New(t)
	testClient.cl.BeginBlockSync(abci.RequestBeginBlock{})
	ctx := testApp.NewContext(sdk.RunTxModeDeliver, abci.Header{})
	InitAccounts(ctx, testApp)
	testApp.DexKeeper.ClearOrderBook("BTC-000_CE")
	ethPair := types.NewTradingPair("ETH-000", "CE", 97e8)
	testApp.DexKeeper.PairMapper.AddTradingPair(ctx, ethPair)
	testApp.DexKeeper.AddEngine(ethPair)
	testApp.DexKeeper.GetEngines()["ETH-000_CE"].LastMatchHeight = -1
	btcPair := types.NewTradingPair("BTC-000", "CE", 96e8)
	testApp.DexKeeper.PairMapper.AddTradingPair(ctx, btcPair)
	testApp.DexKeeper.AddEngine(btcPair)
	testApp.DexKeeper.GetEngines()["BTC-000_CE"].LastMatchHeight = -1
	testApp.DexKeeper.FeeManager.UpdateConfig(newTestFeeConfig())

	// setup accounts
	am := testApp.AccountKeeper
	acc := Account(0)
	acc2 := Account(1)
	acc3 := Account(2)
	add := acc.GetAddress()
	add2 := acc2.GetAddress()
	add3 := acc3.GetAddress()
	ResetAccounts(ctx, testApp, 100000e8, 100000e8, 100000e8)

	/*	--------------------------------------------------------------
		SUM    SELL    PRICE    BUY    SUM    EXECUTION    IMBALANCE
		1500           102      300    300    300          -1200
		1500           101             300    300          -1200
		1500           100      100    400    400          -1100
		1500           99       200    600    600          -900
		1500   250     98       300    900    900          -600
		1250   250     97              900    900          -350
		1000   1000    96              900    900          -100*
	*/
	t.Log(GetAvail(ctx, add, "BTC-000"))
	t.Log(GetAvail(ctx, add, "CE"))
	msg := o.NewNewOrderMsg(add, genOrderID(add, 0, ctx, am), 1, "BTC-000_CE", 102e8, 300e8)
	res, e := testClient.DeliverTxSync(msg, testApp.Codec)
	t.Log(GetAvail(ctx, add, "BTC-000"))
	t.Log(GetAvail(ctx, add, "CE"))
	msg = o.NewNewOrderMsg(add, genOrderID(add, 1, ctx, am), 1, "BTC-000_CE", 100e8, 100e8)
	res, e = testClient.DeliverTxSync(msg, testApp.Codec)
	t.Log(GetAvail(ctx, add, "BTC-000"))
	t.Log(GetAvail(ctx, add, "CE"))

	msg = o.NewNewOrderMsg(add2, genOrderID(add2, 0, ctx, am), 2, "BTC-000_CE", 96e8, 1000e8)
	res, e = testClient.DeliverTxSync(msg, testApp.Codec)
	msg = o.NewNewOrderMsg(add2, genOrderID(add2, 1, ctx, am), 2, "BTC-000_CE", 97e8, 250e8)
	res, e = testClient.DeliverTxSync(msg, testApp.Codec)
	msg = o.NewNewOrderMsg(add2, genOrderID(add2, 2, ctx, am), 2, "BTC-000_CE", 98e8, 250e8)
	res, e = testClient.DeliverTxSync(msg, testApp.Codec)
	msg = o.NewNewOrderMsg(add, genOrderID(add, 2, ctx, am), 1, "BTC-000_CE", 99e8, 200e8)
	res, e = testClient.DeliverTxSync(msg, testApp.Codec)
	t.Logf("res is %v and error is %v", res, e)
	msg = o.NewNewOrderMsg(add, genOrderID(add, 3, ctx, am), 1, "BTC-000_CE", 98e8, 300e8)
	res, e = testClient.DeliverTxSync(msg, testApp.Codec)
	buys, sells, pendingMatch := getOrderBook("BTC-000_CE")
	assert.Equal(4, len(buys))
	assert.Equal(3, len(sells))

	assert.Equal(true, pendingMatch)
	testApp.DexKeeper.MatchAndAllocateSymbols(ctx, nil, false)
	buys, sells, pendingMatch = getOrderBook("BTC-000_CE")

	assert.Equal(0, len(buys))
	assert.Equal(3, len(sells))
	assert.Equal(false, pendingMatch)

	trades, lastPx := testApp.DexKeeper.GetLastTradesForPair("BTC-000_CE")
	assert.Equal(int64(96e8), lastPx)
	assert.Equal(4, len(trades))
	// total execution is 900e8 BTC-000 @ price 96e8, notional is 86400e8, fee is 43.2e8 CE
	assert.Equal(sdk.Coins{sdk.NewCoin("CE", 86.4e8)}, fees.Pool.BlockFees().Tokens)
	assert.Equal(int64(100900e8), GetAvail(ctx, add, "BTC-000"))
	assert.Equal(int64(13556.8e8), GetAvail(ctx, add, "CE"))
	assert.Equal(int64(0), GetLocked(ctx, add, "BTC-000"))
	assert.Equal(int64(0), GetLocked(ctx, add, "CE"))
	assert.Equal(int64(98500e8), GetAvail(ctx, add2, "BTC-000"))
	assert.Equal(int64(186356.8e8), GetAvail(ctx, add2, "CE"))
	assert.Equal(int64(600e8), GetLocked(ctx, add2, "BTC-000"))
	assert.Equal(int64(0), GetLocked(ctx, add2, "CE"))

	// test ETH-000_CE pair
	/*	--------------------------------------------------------------
		SUM    SELL    PRICE    BUY    SUM    EXECUTION    IMBALANCE
		110            102      30     30     30           -80
		110            101      10     40     40           -70
		110            100             40     40           -70
		110            99       50     90     90           -20
		110    10      98              90     90           -20
		100    50      97              90     90           -10*
		50             96       15     105    50           55
		50     50      95              105    50           55
	*/

	msg = o.NewNewOrderMsg(add2, genOrderID(add2, 3, ctx, am), 1, "ETH-000_CE", 102e8, 30e8)
	res, e = testClient.DeliverTxSync(msg, testApp.Codec)
	msg = o.NewNewOrderMsg(add2, genOrderID(add2, 4, ctx, am), 1, "ETH-000_CE", 101e8, 10e8)
	res, e = testClient.DeliverTxSync(msg, testApp.Codec)
	msg = o.NewNewOrderMsg(add3, genOrderID(add3, 0, ctx, am), 2, "ETH-000_CE", 95e8, 50e8)
	res, e = testClient.DeliverTxSync(msg, testApp.Codec)
	msg = o.NewNewOrderMsg(add3, genOrderID(add3, 1, ctx, am), 2, "ETH-000_CE", 98e8, 10e8)
	res, e = testClient.DeliverTxSync(msg, testApp.Codec)
	msg = o.NewNewOrderMsg(add3, genOrderID(add3, 2, ctx, am), 2, "ETH-000_CE", 97e8, 50e8)
	res, e = testClient.DeliverTxSync(msg, testApp.Codec)
	msg = o.NewNewOrderMsg(add2, genOrderID(add2, 5, ctx, am), 1, "ETH-000_CE", 96e8, 15e8)
	res, e = testClient.DeliverTxSync(msg, testApp.Codec)
	msg = o.NewNewOrderMsg(add2, genOrderID(add2, 6, ctx, am), 1, "ETH-000_CE", 99e8, 50e8)
	res, e = testClient.DeliverTxSync(msg, testApp.Codec)
	t.Logf("res is %v and error is %v", res, e)

	buys, sells, _ = getOrderBook("BTC-000_CE")
	assert.Equal(0, len(buys))
	assert.Equal(3, len(sells))
	buys, sells, _ = getOrderBook("ETH-000_CE")
	assert.Equal(4, len(buys))
	assert.Equal(3, len(sells))

	testApp.DexKeeper.MatchAndAllocateSymbols(ctx, nil, false)
	buys, sells, _ = getOrderBook("ETH-000_CE")

	t.Logf("buys: %v", buys)
	t.Logf("sells: %v", sells)
	assert.Equal(1, len(buys))
	assert.Equal(2, len(sells))
	buys, sells, _ = getOrderBook("BTC-000_CE")
	assert.Equal(0, len(buys))
	assert.Equal(3, len(sells))
	trades, lastPx = testApp.DexKeeper.GetLastTradesForPair("ETH-000_CE")
	assert.Equal(int64(97e8), lastPx)
	assert.Equal(4, len(trades))
	// total execution is 90e8 ETH @ price 97e8, notional is 8730e8
	// fee for this round is 8.73e8 CE, totalFee is 95.13e8 CE
	assert.Equal(sdk.Coins{sdk.NewCoin("CE", 95.13e8)}, fees.Pool.BlockFees().Tokens)
	fees.Pool.Clear()
	assert.Equal(int64(100900e8), GetAvail(ctx, add, "BTC-000"))
	assert.Equal(int64(13556.8e8), GetAvail(ctx, add, "CE"))
	assert.Equal(int64(0), GetLocked(ctx, add, "BTC-000"))
	assert.Equal(int64(0), GetLocked(ctx, add, "CE"))
	assert.Equal(int64(98500e8), GetAvail(ctx, add2, "BTC-000"))
	assert.Equal(int64(600e8), GetLocked(ctx, add2, "BTC-000"))
	// for buy, still locked = 15*96=1440, spent 8730
	// so reserve 1440+8730 = 10170
	// fee is 4.365e8 CE
	assert.Equal(int64(176182.435e8), GetAvail(ctx, add2, "CE"))
	assert.Equal(int64(1440e8), GetLocked(ctx, add2, "CE"))
	assert.Equal(int64(100090e8), GetAvail(ctx, add2, "ETH-000"))
	assert.Equal(int64(0), GetLocked(ctx, add2, "ETH-000"))
	assert.Equal(int64(100000e8), GetAvail(ctx, add3, "BTC-000"))
	assert.Equal(int64(0), GetLocked(ctx, add3, "BTC-000"))
	assert.Equal(int64(108725.635e8), GetAvail(ctx, add3, "CE"))
	assert.Equal(int64(0), GetLocked(ctx, add3, "CE"))
	assert.Equal(int64(99890e8), GetAvail(ctx, add3, "ETH-000"))
	assert.Equal(int64(20e8), GetLocked(ctx, add3, "ETH-000"))
	fees.Pool.Clear()
}

func Test_handleCancelOrder_CheckTx(t *testing.T) {
	assert := assert.New(t)
	testClient.cl.BeginBlockSync(abci.RequestBeginBlock{})
	ctx := testApp.NewContext(sdk.RunTxModeDeliver, abci.Header{})
	InitAccounts(ctx, testApp)
	testApp.DexKeeper.ClearOrderBook("BTC-000_CE")
	tradingPair := types.NewTradingPair("BTC-000", "CE", 1e8)
	testApp.DexKeeper.PairMapper.AddTradingPair(ctx, tradingPair)
	testApp.DexKeeper.AddEngine(tradingPair)
	testApp.DexKeeper.GetEngines()["BTC-000_CE"].LastMatchHeight = -1
	testApp.DexKeeper.FeeManager.UpdateConfig(newTestFeeConfig())

	// setup accounts
	add := Account(0).GetAddress()
	oid := fmt.Sprintf("%X-0", add)
	add2 := Account(1).GetAddress()

	msg := o.NewCancelOrderMsg(add, "BTC-000_CE", "DOESNOTEXIST-0")
	res, e := testClient.DeliverTxSync(msg, testApp.Codec)
	assert.Regexp(".*Failed to find order \\[DOESNOTEXIST-0\\].*", res.GetLog())
	assert.Nil(e)
	newMsg := o.NewNewOrderMsg(add, oid, 1, "BTC-000_CE", 355e8, 1e8)
	res, e = testClient.DeliverTxSync(newMsg, testApp.Codec)
	assert.Equal(uint32(0), res.Code)
	assert.Nil(e)
	assert.Equal(int64(145e8), GetAvail(ctx, add, "CE"))
	assert.Equal(int64(355e8), GetLocked(ctx, add, "CE"))
	assert.Equal(int64(200e8), GetAvail(ctx, add, "BTC-000"))
	assert.Equal(int64(0), GetLocked(ctx, add, "BTC-000"))
	msg = o.NewCancelOrderMsg(add2, "BTC-000_CE", oid)
	res, e = testClient.DeliverTxSync(msg, testApp.Codec)
	assert.Regexp(".*does not belong to transaction sender.*", res.GetLog())
	msg = o.NewCancelOrderMsg(add, "BTC-000_CE", oid)
	res, e = testClient.DeliverTxSync(msg, testApp.Codec)
	assert.Equal(uint32(0), res.Code)
	assert.Nil(e)
	assert.Equal(int64(500e8-2e4), GetAvail(ctx, add, "CE"))
	assert.Equal(int64(0), GetLocked(ctx, add, "CE"))
	assert.Equal(int64(200e8), GetAvail(ctx, add, "BTC-000"))
	assert.Equal(int64(0), GetLocked(ctx, add, "BTC-000"))
}
