package apptest

import (
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	abci "github.com/tendermint/tendermint/abci/types"

	"github.com/Mustafa-Agha/node/common/utils"
	"github.com/Mustafa-Agha/node/plugins/dex/matcheng"
	"github.com/Mustafa-Agha/node/plugins/dex/order"
)

/*
test #1a: multiple buy orders (same price level) overflow int64 max
*/
func Test_Overflow_1a_new(t *testing.T) {
	assert := assert.New(t)

	_, ctx, accs := SetupTest_new(1)
	addr0 := accs[0].GetAddress()

	// 10 * 1e18 > int64 max
	for i := 0; i < 10; i++ {
		oid := GetOrderId(addr0, int64(i), ctx)
		msg := order.NewNewOrderMsg(addr0, oid, 1, "BTC-000_CE", 1, 1e18)
		res, err := testClient.DeliverTxSync(msg, testApp.Codec)
		assert.NoError(err)
		if i < 9 {
			assert.Equal(uint32(0), res.Code)
		} else {
			assert.True(strings.Contains(res.Log, "order quantity is too large to be placed on this price level"))
		}
	}

	buys, _ := GetOrderBook("BTC-000_CE")
	assert.Equal(utils.Fixed8(9e18), buys[0].qty)
}

/*
test #1b: multiple buy orders (diff price levels) overflow int64 max, init price 1
*/
func Test_Overflow_1b_new(t *testing.T) {
	assert := assert.New(t)

	addr, ctx, accs := SetupTest_new(1)
	addr0 := accs[0].GetAddress()
	addr1 := accs[1].GetAddress()

	ctx = UpdateContextC(addr, ctx, 1)

	/* sum of buy side overflowed as [10e18] > int64 max
	sum    sell    price    buy    sum      exec    imbal
	1e13   	       10*      1e18   1e18     1e13    the smallest abs
	1e13   	       9        1e18   2e18     1e13    -
	1e13           8        1e18   3e18     1e13    -
	1e13           7        1e18   4e18     1e13    -
	1e13           6        1e18   5e18     1e13    -
	1e13           5        1e18   6e18     1e13    -
	1e13           4        1e18   7e18     1e13    -
	1e13           3        1e18   8e18     1e13    -
	1e13           2        1e18   9e18     1e13    -
	1e13   1e13    1        1e18   [10e18]  1e13    the largest abs
	*/

	// although sum of buy side overflowed, in this case, match and allocation of orders can still be completed

	for i := 0; i < 10; i++ {
		oid := GetOrderId(addr0, int64(i), ctx)
		msg := order.NewNewOrderMsg(addr0, oid, 1, "BTC-000_CE", int64(i+1), 1e18)
		res, err := testClient.DeliverTxSync(msg, testApp.Codec)
		assert.NoError(err)
		assert.Equal(uint32(0), res.Code)
	}

	oidS := GetOrderId(addr1, 0, ctx)
	msgS1 := order.NewNewOrderMsg(addr1, oidS, 2, "BTC-000_CE", int64(1), 1e13)
	res, err := testClient.DeliverTxSync(msgS1, testApp.Codec)
	assert.NoError(err)
	assert.Equal(uint32(0), res.Code)

	buys, sells := GetOrderBook("BTC-000_CE")
	assert.Equal(10, len(buys))
	assert.Equal(1, len(sells))

	assert.Equal(int64(100000e8), GetAvail(ctx, addr0, "BTC-000"))
	assert.Equal(int64(94500e8), GetAvail(ctx, addr0, "CE"))
	assert.Equal(int64(5500e8), GetLocked(ctx, addr0, "CE"))
	assert.Equal(int64(0), GetAvail(ctx, addr1, "BTC-000"))
	assert.Equal(int64(100000e8), GetAvail(ctx, addr1, "CE"))
	assert.Equal(int64(100000e8), GetLocked(ctx, addr1, "BTC-000"))

	testClient.cl.EndBlockSync(abci.RequestEndBlock{})

	trades, lastPx := testApp.DexKeeper.GetLastTradesForPair("BTC-000_CE")
	assert.Equal(int64(10), lastPx)
	assert.Equal(1, len(trades))

	buys, sells = GetOrderBook("BTC-000_CE")
	assert.Equal(10, len(buys))
	assert.Equal(0, len(sells))

	assert.Equal(int64(200000e8), GetAvail(ctx, addr0, "BTC-000"))
	assert.Equal(int64(94499.99999500e8), GetAvail(ctx, addr0, "CE"))
	assert.Equal(int64(5499.9900e8), GetLocked(ctx, addr0, "CE"))
	assert.Equal(int64(0), GetAvail(ctx, addr1, "BTC-000"))
	assert.Equal(int64(100000.00999500e8), GetAvail(ctx, addr1, "CE"))
	assert.Equal(int64(0), GetLocked(ctx, addr1, "BTC-000"))
}

/*
test #1c: multiple buy orders (diff price levels) overflow int64 max, init price 1e4
*/
func Test_Overflow_1c_new(t *testing.T) {
	assert := assert.New(t)

	addr, ctx, accs := SetupTest_new(1e4)
	addr0 := accs[0].GetAddress()
	addr1 := accs[1].GetAddress()

	ctx = UpdateContextC(addr, ctx, 1)

	/* sum of buy side overflowed as [10e18] > int64 max
	sum    sell    price    buy    sum      exec    imbal
	1e9   	       10*      1e18   1e18     1e9    the smallest abs
	1e9   	       9        1e18   2e18     1e9    -
	1e9            8        1e18   3e18     1e9    -
	1e9            7        1e18   4e18     1e9    -
	1e9            6        1e18   5e18     1e9    -
	1e9            5        1e18   6e18     1e9    -
	1e9            4        1e18   7e18     1e9    -
	1e9            3        1e18   8e18     1e9    -
	1e9            2        1e18   9e18     1e9    -
	1e9    1e9     1        1e18   [10e18]  1e9    the largest abs
	*/

	// although sum of buy side overflowed, in this case, match and allocation of orders can still be completed

	for i := 0; i < 10; i++ {
		oid := GetOrderId(addr0, int64(i), ctx)
		msg := order.NewNewOrderMsg(addr0, oid, 1, "BTC-000_CE", int64(i+1), 1e18)
		res, err := testClient.DeliverTxSync(msg, testApp.Codec)
		assert.NoError(err)
		assert.Equal(uint32(0), res.Code)
	}

	oidS1 := GetOrderId(addr1, 0, ctx)
	msgS1 := order.NewNewOrderMsg(addr1, oidS1, 2, "BTC-000_CE", int64(1), 1e9)
	res, err := testClient.DeliverTxSync(msgS1, testApp.Codec)
	assert.NoError(err)
	assert.Equal(uint32(0), res.Code)

	buys, sells := GetOrderBook("BTC-000_CE")
	assert.Equal(10, len(buys))
	assert.Equal(1, len(sells))

	assert.Equal(int64(100000e8), GetAvail(ctx, addr0, "BTC-000"))
	assert.Equal(int64(94500e8), GetAvail(ctx, addr0, "CE"))
	assert.Equal(int64(5500e8), GetLocked(ctx, addr0, "CE"))
	assert.Equal(int64(99990e8), GetAvail(ctx, addr1, "BTC-000"))
	assert.Equal(int64(100000e8), GetAvail(ctx, addr1, "CE"))
	assert.Equal(int64(10e8), GetLocked(ctx, addr1, "BTC-000"))

	testClient.cl.EndBlockSync(abci.RequestEndBlock{})

	trades, lastPx := testApp.DexKeeper.GetLastTradesForPair("BTC-000_CE")
	assert.Equal(int64(10), lastPx)
	assert.Equal(1, len(trades))

	buys, sells = GetOrderBook("BTC-000_CE")
	assert.Equal(10, len(buys))
	assert.Equal(0, len(sells))

	// fee charged from receiving token btc-000, as fee in ce is < 1
	assert.Equal(int64(100009.9900e8), GetAvail(ctx, addr0, "BTC-000"))
	assert.Equal(int64(94500e8), GetAvail(ctx, addr0, "CE"))
	assert.Equal(int64(5499.99999900e8), GetLocked(ctx, addr0, "CE"))
	// in this case, it is expected that no fee charged for sell side
	assert.Equal(int64(99990e8), GetAvail(ctx, addr1, "BTC-000"))
	assert.Equal(int64(100000.00000100e8), GetAvail(ctx, addr1, "CE"))
	assert.Equal(int64(0), GetLocked(ctx, addr1, "BTC-000"))
}

/*
test #1d: additional test case using very cheap ce, not really overflow related
*/
func Test_Overflow_1d_new(t *testing.T) {
	assert := assert.New(t)

	addr, ctx, accs := SetupTest_new(1e18)
	addr0 := accs[0].GetAddress()
	addr1 := accs[1].GetAddress()
	ResetAccount(ctx, addr0, 2e18, 100000e8, 100000e8)
	ResetAccount(ctx, addr1, 2e18, 100000e8, 100000e8)

	ctx = UpdateContextC(addr, ctx, 1)

	oidB1 := GetOrderId(addr0, 0, ctx)
	msgB1 := order.NewNewOrderMsg(addr0, oidB1, 1, "BTC-000_CE", 1e18, 1)
	_, err := testClient.DeliverTxSync(msgB1, testApp.Codec)
	assert.NoError(err)

	oidS1 := GetOrderId(addr1, 0, ctx)
	msgS1 := order.NewNewOrderMsg(addr1, oidS1, 2, "BTC-000_CE", 1e18, 1)
	_, err = testClient.DeliverTxSync(msgS1, testApp.Codec)
	assert.NoError(err)

	buys, sells := GetOrderBook("BTC-000_CE")
	assert.Equal(1, len(buys))
	assert.Equal(1, len(sells))

	testClient.cl.EndBlockSync(abci.RequestEndBlock{})

	trades, lastPx := testApp.DexKeeper.GetLastTradesForPair("BTC-000_CE")
	assert.Equal(int64(1e18), lastPx)
	assert.Equal(1, len(trades))

	buys, sells = GetOrderBook("BTC-000_CE")
	assert.Equal(0, len(buys))
	assert.Equal(0, len(sells))

	assert.Equal(int64(10000000000001), GetAvail(ctx, addr0, "BTC-000"))
	assert.Equal(int64(1999999989995000000), GetAvail(ctx, addr0, "CE"))
	assert.Equal(int64(0), GetLocked(ctx, addr0, "CE"))
	assert.Equal(int64(9999999999999), GetAvail(ctx, addr1, "BTC-000"))
	assert.Equal(int64(2000000009995000000), GetAvail(ctx, addr1, "CE"))
	assert.Equal(int64(0), GetLocked(ctx, addr1, "BTC-000"))
}

/*
test #2a: multiple sell orders (same price level) overflow int64 max
*/
func Test_Overflow_2a_new(t *testing.T) {
	assert := assert.New(t)

	_, ctx, accs := SetupTest_new(1e18)
	addr0 := accs[0].GetAddress()

	for i := 0; i < 10; i++ {
		oid := GetOrderId(addr0, int64(i), ctx)
		msg := order.NewNewOrderMsg(addr0, oid, 2, "BTC-000_CE", 1e18, 1e8)
		res, err := testClient.DeliverTxSync(msg, testApp.Codec)
		assert.NoError(err)
		assert.Equal(uint32(0), res.Code)
	}

	_, sells := GetOrderBook("BTC-000_CE")
	assert.Equal(utils.Fixed8(10e8), sells[0].qty)
	// grand these orders, when the total amount (q*p + q*p + ... ) of a pair from one address is greater than int64 max
}

/*
test #2b: multiple sell orders (diff price levels) overflow int64 max
*/
func Test_Overflow_2b_new(t *testing.T) {
	assert := assert.New(t)

	_, ctx, accs := SetupTest_new(1e18)
	addr0 := accs[0].GetAddress()

	for i := 0; i < 5; i++ {
		oid := GetOrderId(addr0, int64(i), ctx)
		msg := order.NewNewOrderMsg(addr0, oid, 2, "BTC-000_CE", 1e18*int64(i+1), 1e8)
		res, err := testClient.DeliverTxSync(msg, testApp.Codec)
		assert.NoError(err)
		assert.Equal(uint32(0), res.Code)
	}

	_, sells := GetOrderBook("BTC-000_CE")
	assert.Equal(5, len(sells))
	// grand these orders, when the total amount (q*p + q*p + ... ) of a pair from one address is greater than int64 max
}

/*
test #3: non ce pair (with cheap ce)
*/
func Test_Overflow_3_new(t *testing.T) {
	assert := assert.New(t)

	addr, ctx, accs := SetupTest_new(1e18, 1e18, 10e8)
	addr0 := accs[0].GetAddress()
	addr1 := accs[1].GetAddress()

	ctx = UpdateContextC(addr, ctx, 1)

	oidB1 := GetOrderId(addr0, 0, ctx)
	msgB1 := order.NewNewOrderMsg(addr0, oidB1, 1, "BTC-000_ETH-000", 10e8, 1e8)
	_, err := testClient.DeliverTxSync(msgB1, testApp.Codec)
	assert.NoError(err)

	oidS1 := GetOrderId(addr1, 0, ctx)
	msgS1 := order.NewNewOrderMsg(addr1, oidS1, 2, "BTC-000_ETH-000", 10e8, 1e8)
	_, err = testClient.DeliverTxSync(msgS1, testApp.Codec)
	assert.NoError(err)

	buys, sells := GetOrderBook("BTC-000_ETH-000")
	assert.Equal(1, len(buys))
	assert.Equal(1, len(sells))

	testClient.cl.EndBlockSync(abci.RequestEndBlock{})

	trades, lastPx := testApp.DexKeeper.GetLastTradesForPair("BTC-000_ETH-000")
	assert.Equal(int64(10e8), lastPx)
	assert.Equal(1, len(trades))
	for i, trade := range trades {
		fmt.Printf("#%d: p: %d; q: %d; s: %d\n",
			i, trade.LastPx, trade.LastQty, trade.TickType)
	}
	assert.Equal(int64(1e8), trades[0].LastQty)
	assert.Equal(int8(matcheng.Neutral), trades[0].TickType)
	assert.Equal(int64(0.0010e8), trades[0].BuyerFee.Tokens[0].Amount)
	assert.Equal("BTC-000", trades[0].BuyerFee.Tokens[0].Denom)
	assert.Equal(int64(0.0100e8), trades[0].SellerFee.Tokens[0].Amount)
	assert.Equal("ETH-000", trades[0].SellerFee.Tokens[0].Denom)

	buys, sells = GetOrderBook("BTC-000_ETH-000")
	assert.Equal(0, len(buys))
	assert.Equal(0, len(sells))

	assert.Equal(int64(99990e8), GetAvail(ctx, addr0, "ETH-000"))
	assert.Equal(int64(100000.9990e8), GetAvail(ctx, addr0, "BTC-000"))
	// for buy side: insufficent ce (1x1e18 > 100000e8), so fee is deducted from btc-000 => 100001 - 1 * 0.001
	assert.Equal(int64(100000e8), GetAvail(ctx, addr0, "CE"))
	assert.Equal(int64(100009.9900e8), GetAvail(ctx, addr1, "ETH-000"))
	assert.Equal(int64(99999e8), GetAvail(ctx, addr1, "BTC-000"))
	// for sell side: it is overflowed (10x1e18 > int64 max), so fee is deducted from eth-000 => 10e8 * 0.001 = 0.01e8
	assert.Equal(int64(100000e8), GetAvail(ctx, addr1, "CE"))
}

/*
test #4: non ce pair (with expansive ce)
*/
func Test_Overflow_4_new(t *testing.T) {
	assert := assert.New(t)

	addr, ctx, accs := SetupTest_new(1, 1, 10e8)
	addr0 := accs[0].GetAddress()
	addr1 := accs[1].GetAddress()

	ctx = UpdateContextC(addr, ctx, 1)

	oidB1 := GetOrderId(addr0, 0, ctx)
	msgB1 := order.NewNewOrderMsg(addr0, oidB1, 1, "BTC-000_ETH-000", 10e8, 1e8)
	_, err := testClient.DeliverTxSync(msgB1, testApp.Codec)
	assert.NoError(err)

	oidS1 := GetOrderId(addr1, 0, ctx)
	msgS1 := order.NewNewOrderMsg(addr1, oidS1, 2, "BTC-000_ETH-000", 10e8, 1e8)
	_, err = testClient.DeliverTxSync(msgS1, testApp.Codec)
	assert.NoError(err)

	buys, sells := GetOrderBook("BTC-000_ETH-000")
	assert.Equal(1, len(buys))
	assert.Equal(1, len(sells))

	testClient.cl.EndBlockSync(abci.RequestEndBlock{})

	trades, lastPx := testApp.DexKeeper.GetLastTradesForPair("BTC-000_ETH-000")
	assert.Equal(int64(10e8), lastPx)
	assert.Equal(1, len(trades))

	buys, sells = GetOrderBook("BTC-000_ETH-000")
	assert.Equal(0, len(buys))
	assert.Equal(0, len(sells))

	assert.Equal(int64(100000.9990e8), GetAvail(ctx, addr0, "BTC-000"))
	assert.Equal(int64(99990e8), GetAvail(ctx, addr0, "ETH-000"))
	assert.Equal(int64(100000e8), GetAvail(ctx, addr0, "CE"))
	assert.Equal(int64(100009.9900e8), GetAvail(ctx, addr1, "ETH-000"))
	assert.Equal(int64(99999e8), GetAvail(ctx, addr1, "BTC-000"))
	assert.Equal(int64(100000e8), GetAvail(ctx, addr1, "CE"))
}

/*
test #5: sum of orders overflowed leads to unexpected trade failure
*/
func Test_Overflow_5_new(t *testing.T) {
	assert := assert.New(t)

	addr, ctx, accs := SetupTest_new(1)
	addr0 := accs[0].GetAddress()
	addr1 := accs[1].GetAddress()
	ResetAccount(ctx, addr0, 100000e8, 0, 0)
	ResetAccount(ctx, addr1, 100000e8, 9e18, 0)

	ctx = UpdateContextC(addr, ctx, 1)

	/* sum of buy side overflowed as 10e18 > int64 max
	sum    sell    price    buy    sum      exec    imbal
	9e18   	       5        5e18   5e18     5e18    -4e18
	9e18   	       4        5e18   [10e18]  9e18    unknown
	9e18   9e18	   3               [10e18]  9e18    unknown
	*/

	oidB1 := GetOrderId(addr0, 0, ctx)
	msgB1 := order.NewNewOrderMsg(addr0, oidB1, 1, "BTC-000_CE", 5, 5e18)
	_, err := testClient.DeliverTxSync(msgB1, testApp.Codec)
	assert.NoError(err)

	oidB2 := GetOrderId(addr0, 1, ctx)
	msgB2 := order.NewNewOrderMsg(addr0, oidB2, 1, "BTC-000_CE", 4, 5e18)
	_, err = testClient.DeliverTxSync(msgB2, testApp.Codec)
	assert.NoError(err)

	oidS1 := GetOrderId(addr1, 0, ctx)
	msgS1 := order.NewNewOrderMsg(addr1, oidS1, 2, "BTC-000_CE", 3, 9e18)
	_, err = testClient.DeliverTxSync(msgS1, testApp.Codec)
	assert.NoError(err)

	buys, sells := GetOrderBook("BTC-000_CE")
	assert.Equal(2, len(buys))
	assert.Equal(1, len(sells))

	testClient.cl.EndBlockSync(abci.RequestEndBlock{})

	trades, lastPx := testApp.DexKeeper.GetLastTradesForPair("BTC-000_CE")
	assert.Equal(int64(3), lastPx)
	assert.Equal(2, len(trades))

	buys, sells = GetOrderBook("BTC-000_CE")
	assert.Equal(1, len(buys))
	assert.Equal(0, len(sells))

	assert.Equal(int64(9e18), GetAvail(ctx, addr0, "BTC-000"))
	assert.Equal(int64(96898.6500e8), GetAvail(ctx, addr0, "CE"))
	assert.Equal(int64(400e8), GetLocked(ctx, addr0, "CE"))
	assert.Equal(int64(0), GetAvail(ctx, addr1, "BTC-000"))
	assert.Equal(int64(102698.6500e8), GetAvail(ctx, addr1, "CE"))
	assert.Equal(int64(0), GetLocked(ctx, addr1, "BTC-000"))
}
