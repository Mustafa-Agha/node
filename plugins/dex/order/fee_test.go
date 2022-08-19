package order

import (
	"testing"

	"github.com/stretchr/testify/require"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/Mustafa-Agha/node/common/testutils"
	"github.com/Mustafa-Agha/node/common/types"
	"github.com/Mustafa-Agha/node/plugins/dex/matcheng"
	dextype "github.com/Mustafa-Agha/node/plugins/dex/types"
)

func NewTestFeeConfig() FeeConfig {
	feeConfig := NewFeeConfig()
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

func feeManagerCalcTradeFeeForSingleTransfer(t *testing.T, symbol string) {
	ctx, am, keeper := setup()
	keeper.FeeManager.UpdateConfig(NewTestFeeConfig())
	keeper.AddEngine(dextype.NewTradingPair(symbol, "TNT", 1e7))
	keeper.AddEngine(dextype.NewTradingPair("TNT", "XYZ-111", 1e7))
	_, acc := testutils.NewAccount(ctx, am, 0)
	tran := Transfer{
		inAsset:  symbol,
		in:       1000,
		outAsset: "TNT",
		out:      100,
	}
	// no enough tnt or native fee rounding to 0
	fee := keeper.FeeManager.calcTradeFeeFromTransfer(acc.GetCoins(), &tran, keeper.engines)
	require.Equal(t, sdk.Coins{{symbol, 1}}, fee.Tokens)
	_, acc = testutils.NewAccount(ctx, am, 100)
	fee = keeper.FeeManager.calcTradeFeeFromTransfer(acc.GetCoins(), &tran, keeper.engines)
	require.Equal(t, sdk.Coins{{symbol, 1}}, fee.Tokens)

	tran = Transfer{
		inAsset:  symbol,
		in:       1000000,
		outAsset: "TNT",
		out:      10000,
	}
	_, acc = testutils.NewAccount(ctx, am, 1)
	fee = keeper.FeeManager.calcTradeFeeFromTransfer(acc.GetCoins(), &tran, keeper.engines)
	require.Equal(t, sdk.Coins{{symbol, 1000}}, fee.Tokens)
	_, acc = testutils.NewAccount(ctx, am, 100)
	fee = keeper.FeeManager.calcTradeFeeFromTransfer(acc.GetCoins(), &tran, keeper.engines)
	require.Equal(t, sdk.Coins{{"TNT", 5}}, fee.Tokens)

	tran = Transfer{
		inAsset:  "TNT",
		in:       100,
		outAsset: symbol,
		out:      1000,
	}
	_, acc = testutils.NewAccount(ctx, am, 100)
	fee = keeper.FeeManager.calcTradeFeeFromTransfer(acc.GetCoins(), &tran, keeper.engines)
	require.Equal(t, sdk.Coins{{"TNT", 0}}, fee.Tokens)

	tran = Transfer{
		inAsset:  "TNT",
		in:       10000,
		outAsset: symbol,
		out:      100000,
	}
	fee = keeper.FeeManager.calcTradeFeeFromTransfer(acc.GetCoins(), &tran, keeper.engines)
	require.Equal(t, sdk.Coins{{"TNT", 5}}, fee.Tokens)

	tran = Transfer{
		inAsset:  symbol,
		in:       100000,
		outAsset: "XYZ-111",
		out:      100000,
	}
	acc.SetCoins(sdk.Coins{{symbol, 1000000}, {"TNT", 100}})
	fee = keeper.FeeManager.calcTradeFeeFromTransfer(acc.GetCoins(), &tran, keeper.engines)
	require.Equal(t, sdk.Coins{{"TNT", 5}}, fee.Tokens)
	tran = Transfer{
		inAsset:  "XYZ-111",
		in:       100000,
		outAsset: symbol,
		out:      100000,
	}
	acc.SetCoins(sdk.Coins{{"XYZ-111", 1000000}, {"TNT", 1000}})
	fee = keeper.FeeManager.calcTradeFeeFromTransfer(acc.GetCoins(), &tran, keeper.engines)
	require.Equal(t, sdk.Coins{{"TNT", 500}}, fee.Tokens)
}

func TestFeeManager_calcTradeFeeForSingleTransfer(t *testing.T) {
	setChainVersion()
	defer resetChainVersion()
	symbol := "ABC-000"
	feeManagerCalcTradeFeeForSingleTransfer(t, symbol)
}

func TestFeeManager_calcTradeFeeForSingleTransferMini(t *testing.T) {
	setChainVersion()
	defer resetChainVersion()
	symbol := "ABC-000M"
	feeManagerCalcTradeFeeForSingleTransfer(t, symbol)
}

func TestFeeManager_CalcTradesFee(t *testing.T) {
	setChainVersion()
	defer resetChainVersion()
	ctx, am, keeper := setup()
	keeper.FeeManager.UpdateConfig(NewTestFeeConfig())
	keeper.AddEngine(dextype.NewTradingPair("ABC-000", "TNT", 1e7))
	keeper.AddEngine(dextype.NewTradingPair("XYZ-111", "TNT", 2e7))
	keeper.AddEngine(dextype.NewTradingPair("ABC-000", "BTC", 1e4))
	keeper.AddEngine(dextype.NewTradingPair("XYZ-111", "BTC", 2e4))
	keeper.AddEngine(dextype.NewTradingPair("TNT", "BTC", 5e5))
	keeper.AddEngine(dextype.NewTradingPair("ABC-000", "XYZ-111", 6e7))
	keeper.AddEngine(dextype.NewTradingPair("ZYX-000M", "TNT", 1e8))

	tradeTransfers := TradeTransfers{
		{inAsset: "ABC-000", outAsset: "TNT", Oid: "1", in: 1e5, out: 2e4, Trade: &matcheng.Trade{}},
		{inAsset: "ABC-000", outAsset: "BTC", Oid: "2", in: 3e5, out: 4e1, Trade: &matcheng.Trade{}},
		{inAsset: "XYZ-111", outAsset: "BTC", Oid: "3", in: 2e6, out: 4e2, Trade: &matcheng.Trade{}},
		{inAsset: "XYZ-111", outAsset: "TNT", Oid: "4", in: 1e7, out: 2e6, Trade: &matcheng.Trade{}},
		{inAsset: "ABC-000", outAsset: "XYZ", Oid: "5", in: 8e6, out: 5e6, Trade: &matcheng.Trade{}},
		{inAsset: "BTC", outAsset: "TNT", Oid: "6", in: 1e8, out: 500e8, Trade: &matcheng.Trade{}},
		{inAsset: "TNT", outAsset: "BTC", Oid: "7", in: 300e8, out: 7e7, Trade: &matcheng.Trade{}},
		{inAsset: "TNT", outAsset: "ABC-000", Oid: "8", in: 5e8, out: 60e8, Trade: &matcheng.Trade{}},
		{inAsset: "ABC-000", outAsset: "TNT", Oid: "9", in: 7e6, out: 5e5, Trade: &matcheng.Trade{}},
		{inAsset: "ABC-000", outAsset: "BTC", Oid: "10", in: 6e5, out: 8e1, Trade: &matcheng.Trade{}},
		{inAsset: "ZYX-000M", outAsset: "TNT", Oid: "11", in: 2e7, out: 2e6, Trade: &matcheng.Trade{}},
	}
	_, acc := testutils.NewAccount(ctx, am, 0)
	_ = acc.SetCoins(sdk.Coins{
		{"ABC-000", 100e8},
		{"TNT", 15251400},
		{"BTC", 10e8},
		{"XYZ-111", 100e8},
		{"ZYX-000M", 100e8},
	})
	fees := keeper.FeeManager.CalcTradesFee(acc.GetCoins(), tradeTransfers, keeper.engines)
	require.Equal(t, "ABC-000:8000;TNT:15251305;BTC:100000;XYZ-111:2000;ZYX-000M:20000", fees.String())
	require.Equal(t, "TNT:250000", tradeTransfers[0].Fee.String())
	require.Equal(t, "TNT:15000000", tradeTransfers[1].Fee.String())
	require.Equal(t, "TNT:10", tradeTransfers[2].Fee.String())
	require.Equal(t, "TNT:250", tradeTransfers[3].Fee.String())
	require.Equal(t, "BTC:100000", tradeTransfers[4].Fee.String())
	require.Equal(t, "TNT:1000", tradeTransfers[5].Fee.String())
	require.Equal(t, "ZYX-000M:20000", tradeTransfers[6].Fee.String())
	require.Equal(t, "TNT:15", tradeTransfers[7].Fee.String())
	require.Equal(t, "TNT:30", tradeTransfers[8].Fee.String())
	require.Equal(t, "ABC-000:8000", tradeTransfers[9].Fee.String())
	require.Equal(t, "XYZ-111:2000", tradeTransfers[10].Fee.String())

	require.Equal(t, sdk.Coins{
		{"ABC-000", 100e8},
		{"TNT", 15251400},
		{"BTC", 10e8},
		{"XYZ-111", 100e8},
		{"ZYX-000M", 100e8},
	}, acc.GetCoins())
}

func TestFeeManager_CalcExpiresFee(t *testing.T) {
	setChainVersion()
	defer resetChainVersion()
	ctx, am, keeper := setup()
	keeper.FeeManager.UpdateConfig(NewTestFeeConfig())
	keeper.AddEngine(dextype.NewTradingPair("ABC-000", "TNT", 1e7))
	keeper.AddEngine(dextype.NewTradingPair("XYZ-111", "TNT", 2e7))
	keeper.AddEngine(dextype.NewTradingPair("TNT", "BTC", 5e5))
	keeper.AddEngine(dextype.NewTradingPair("ZYX-000M", "TNT", 1e8))

	// in TNT
	expireTransfers := ExpireTransfers{
		{inAsset: "ABC-000", Symbol: "ABC-000_TNT", Oid: "1"},
		{inAsset: "ABC-000", Symbol: "ABC-000_BTC", Oid: "2"},
		{inAsset: "XYZ-111", Symbol: "XYZ-111_BTC", Oid: "3"},
		{inAsset: "XYZ-111", Symbol: "XYZ-111_TNT", Oid: "4"},
		{inAsset: "ABC-000", Symbol: "ABC-000_XYZ-111", Oid: "5"},
		{inAsset: "BTC", Symbol: "TNT_BTC", Oid: "6"},
		{inAsset: "TNT", Symbol: "TNT_BTC", Oid: "7"},
		{inAsset: "TNT", Symbol: "ABC-000_TNT", Oid: "8"},
		{inAsset: "ABC-000", Symbol: "ABC-000_TNT", Oid: "9"},
		{inAsset: "ABC-000", Symbol: "ABC-000_BTC", Oid: "10"},
		{inAsset: "ZYX-000M", Symbol: "ZYX-000M_BTC", Oid: "11"},
	}
	_, acc := testutils.NewAccount(ctx, am, 0)
	_ = acc.SetCoins(sdk.Coins{
		{"ABC-000", 100e8},
		{"TNT", 120000},
		{"BTC", 10e8},
		{"XYZ-111", 800000},
		{"ZYX-000M", 900000},
	})
	fees := keeper.FeeManager.CalcExpiresFee(acc.GetCoins(), eventFullyExpire, expireTransfers, keeper.engines, nil)
	require.Equal(t, "ABC-000:1000000;TNT:120000;BTC:500;XYZ-111:800000;ZYX-000M:100000", fees.String())
	require.Equal(t, "TNT:20000", expireTransfers[0].Fee.String())
	require.Equal(t, "TNT:20000", expireTransfers[1].Fee.String())
	require.Equal(t, "TNT:20000", expireTransfers[2].Fee.String())
	require.Equal(t, "TNT:20000", expireTransfers[3].Fee.String())
	require.Equal(t, "TNT:20000", expireTransfers[4].Fee.String())
	require.Equal(t, "TNT:20000", expireTransfers[5].Fee.String())
	require.Equal(t, "ABC-000:1000000", expireTransfers[6].Fee.String())
	require.Equal(t, "BTC:500", expireTransfers[7].Fee.String())
	require.Equal(t, "XYZ-111:500000", expireTransfers[8].Fee.String())
	require.Equal(t, "XYZ-111:300000", expireTransfers[9].Fee.String())
	require.Equal(t, "ZYX-000M:100000", expireTransfers[10].Fee.String())
	require.Equal(t, sdk.Coins{
		{"ABC-000", 100e8},
		{"TNT", 120000},
		{"BTC", 10e8},
		{"XYZ-111", 800000},
		{"ZYX-000M", 900000},
	}, acc.GetCoins())
}

func TestFeeManager_calcTradeFee(t *testing.T) {
	setChainVersion()
	defer resetChainVersion()
	symbol := "ABC-000"
	feeManagerCalcTradeFee(t, symbol)
}

func TestFeeManager_calcTradeFeeMini(t *testing.T) {
	setChainVersion()
	defer resetChainVersion()
	symbol := "ABC-000M"
	feeManagerCalcTradeFee(t, symbol)
}

func feeManagerCalcTradeFee(t *testing.T, symbol string) {
	ctx, am, keeper := setup()
	keeper.FeeManager.UpdateConfig(NewTestFeeConfig())
	keeper.AddEngine(dextype.NewTradingPair(symbol, "TNT", 1e7))
	// TNT
	_, acc := testutils.NewAccount(ctx, am, 0)
	// the tradeIn amount is large enough to make the fee > 0
	tradeIn := sdk.NewCoin(types.NativeTokenSymbol, 100e8)
	fee := keeper.FeeManager.CalcTradeFee(acc.GetCoins(), tradeIn, keeper.engines)
	require.Equal(t, sdk.Coins{sdk.NewCoin(types.NativeTokenSymbol, 5e6)}, fee.Tokens)
	// small tradeIn amount
	tradeIn = sdk.NewCoin(types.NativeTokenSymbol, 100)
	fee = keeper.FeeManager.CalcTradeFee(acc.GetCoins(), tradeIn, keeper.engines)
	require.Equal(t, sdk.Coins{sdk.NewCoin(types.NativeTokenSymbol, 0)}, fee.Tokens)

	// !TNT
	_, acc = testutils.NewAccount(ctx, am, 100)
	// has enough tnt
	tradeIn = sdk.NewCoin(symbol, 1000e8)
	acc.SetCoins(sdk.Coins{sdk.NewCoin(types.NativeTokenSymbol, 1e8)})
	fee = keeper.FeeManager.CalcTradeFee(acc.GetCoins(), tradeIn, keeper.engines)
	require.Equal(t, sdk.Coins{sdk.NewCoin(types.NativeTokenSymbol, 5e6)}, fee.Tokens)
	// no enough tnt
	acc.SetCoins(sdk.Coins{sdk.NewCoin(types.NativeTokenSymbol, 1e6)})
	fee = keeper.FeeManager.CalcTradeFee(acc.GetCoins(), tradeIn, keeper.engines)
	require.Equal(t, sdk.Coins{sdk.NewCoin(symbol, 1e8)}, fee.Tokens)

	// very high price to produce int64 overflow
	keeper.AddEngine(dextype.NewTradingPair(symbol, "TNT", 1e16))
	// has enough tnt
	tradeIn = sdk.NewCoin(symbol, 1000e8)
	acc.SetCoins(sdk.Coins{sdk.NewCoin(types.NativeTokenSymbol, 1e16)})
	fee = keeper.FeeManager.CalcTradeFee(acc.GetCoins(), tradeIn, keeper.engines)
	require.Equal(t, sdk.Coins{sdk.NewCoin(types.NativeTokenSymbol, 5e15)}, fee.Tokens)
	// no enough tnt, fee is within int64
	acc.SetCoins(sdk.Coins{sdk.NewCoin(types.NativeTokenSymbol, 1e15)})
	fee = keeper.FeeManager.CalcTradeFee(acc.GetCoins(), tradeIn, keeper.engines)
	require.Equal(t, sdk.Coins{sdk.NewCoin(symbol, 1e8)}, fee.Tokens)
	// no enough tnt, even the fee overflows
	tradeIn = sdk.NewCoin(symbol, 1e16)
	fee = keeper.FeeManager.CalcTradeFee(acc.GetCoins(), tradeIn, keeper.engines)
	require.Equal(t, sdk.Coins{sdk.NewCoin(symbol, 1e13)}, fee.Tokens)
}

func TestFeeManager_CalcFixedFee(t *testing.T) {
	setChainVersion()
	defer resetChainVersion()
	symbol1 := "ABC-000"
	symbol2 := "BTC-000"
	feeManagerCalcFixedFee(t, symbol1, symbol2)
}

func TestFeeManager_CalcFixedFeeMini(t *testing.T) {
	setChainVersion()
	defer resetChainVersion()
	symbol1 := "ABC-000M"
	symbol2 := "BTC-000M"
	feeManagerCalcFixedFee(t, symbol1, symbol2)
}

func feeManagerCalcFixedFee(t *testing.T, symbol1 string, symbol2 string) {
	ctx, am, keeper := setup()
	keeper.FeeManager.UpdateConfig(NewTestFeeConfig())
	_, acc := testutils.NewAccount(ctx, am, 1e4)
	keeper.AddEngine(dextype.NewTradingPair(symbol1, "TNT", 1e7))
	keeper.AddEngine(dextype.NewTradingPair("TNT", symbol2, 1e5))
	// in TNT
	// no enough TNT, but inAsset == TNT
	fee := keeper.FeeManager.CalcFixedFee(acc.GetCoins(), eventFullyExpire, types.NativeTokenSymbol, keeper.engines)
	require.Equal(t, sdk.Coins{sdk.NewCoin(types.NativeTokenSymbol, 1e4)}, fee.Tokens)
	// enough TNT
	acc.SetCoins(sdk.Coins{sdk.NewCoin(types.NativeTokenSymbol, 3e4)})
	fee = keeper.FeeManager.CalcFixedFee(acc.GetCoins(), eventFullyExpire, types.NativeTokenSymbol, keeper.engines)
	require.Equal(t, sdk.Coins{sdk.NewCoin(types.NativeTokenSymbol, 2e4)}, fee.Tokens)

	fee = keeper.FeeManager.CalcFixedFee(acc.GetCoins(), eventIOCFullyExpire, types.NativeTokenSymbol, keeper.engines)
	require.Equal(t, sdk.Coins{sdk.NewCoin(types.NativeTokenSymbol, 1e4)}, fee.Tokens)

	fee = keeper.FeeManager.CalcFixedFee(acc.GetCoins(), eventFullyCancel, types.NativeTokenSymbol, keeper.engines)
	require.Equal(t, sdk.Coins{sdk.NewCoin(types.NativeTokenSymbol, 2e4)}, fee.Tokens)

	// ABC-000_TNT, sell ABC-000
	fee = keeper.FeeManager.CalcFixedFee(acc.GetCoins(), eventFullyExpire, symbol1, keeper.engines)
	require.Equal(t, sdk.Coins{sdk.NewCoin(types.NativeTokenSymbol, 2e4)}, fee.Tokens)

	// No enough native token, but enough ABC-000
	acc.SetCoins(sdk.Coins{{Denom: types.NativeTokenSymbol, Amount: 1e4}, {Denom: symbol1, Amount: 1e8}})
	fee = keeper.FeeManager.CalcFixedFee(acc.GetCoins(), eventFullyExpire, symbol1, keeper.engines)
	require.Equal(t, sdk.Coins{sdk.NewCoin(symbol1, 1e6)}, fee.Tokens)

	// No enough native token and ABC-000
	acc.SetCoins(sdk.Coins{{Denom: types.NativeTokenSymbol, Amount: 1e4}, {Denom: symbol1, Amount: 1e5}})
	fee = keeper.FeeManager.CalcFixedFee(acc.GetCoins(), eventFullyExpire, symbol1, keeper.engines)
	require.Equal(t, sdk.Coins{sdk.NewCoin(symbol1, 1e5)}, fee.Tokens)

	// TNT_BTC-000, sell BTC-000
	acc.SetCoins(sdk.Coins{{Denom: symbol2, Amount: 1e4}})
	fee = keeper.FeeManager.CalcFixedFee(acc.GetCoins(), eventFullyExpire, symbol2, keeper.engines)
	require.Equal(t, sdk.Coins{sdk.NewCoin(symbol2, 1e2)}, fee.Tokens)

	// extreme prices
	keeper.AddEngine(dextype.NewTradingPair(symbol1, "TNT", 1))
	keeper.AddEngine(dextype.NewTradingPair("TNT", symbol2, 1e16))
	acc.SetCoins(sdk.Coins{{Denom: symbol1, Amount: 1e16}, {Denom: symbol2, Amount: 1e16}})
	fee = keeper.FeeManager.CalcFixedFee(acc.GetCoins(), eventFullyExpire, symbol1, keeper.engines)
	require.Equal(t, sdk.Coins{sdk.NewCoin(symbol1, 1e13)}, fee.Tokens)
	fee = keeper.FeeManager.CalcFixedFee(acc.GetCoins(), eventFullyExpire, symbol2, keeper.engines)
	require.Equal(t, sdk.Coins{sdk.NewCoin(symbol2, 1e13)}, fee.Tokens)
}

func TestFeeManager_calcTradeFeeForSingleTransfer_SupportBUSD(t *testing.T) {
	setChainVersion()
	defer resetChainVersion()
	ctx, am, keeper := setup()
	keeper.FeeManager.UpdateConfig(NewTestFeeConfig())
	keeper.SetBUSDSymbol("BUSD-BD1")

	// existing TNT -> BUSD trading pair
	keeper.AddEngine(dextype.NewTradingPair("TNT", "BUSD-BD1", 1e5))
	keeper.AddEngine(dextype.NewTradingPair("ABC-000", "BUSD-BD1", 1e7))
	keeper.AddEngine(dextype.NewTradingPair("BUSD-BD1", "XYZ-999", 1e6))

	// enough TNT, TNT will be collected
	_, acc := testutils.NewAccount(ctx, am, 1e5)

	// transferred in TNT
	tran := Transfer{
		inAsset: "TNT",
		in:      2e3,
	}
	fee := keeper.FeeManager.calcTradeFeeFromTransfer(acc.GetCoins(), &tran, keeper.engines)
	require.Equal(t, sdk.Coins{{"TNT", 1}}, fee.Tokens)

	// transferred in BUSD-BD1
	tran = Transfer{
		inAsset:  "BUSD-BD1",
		in:       1e3,
		outAsset: "ABC-000",
		out:      1e4,
	}
	fee = keeper.FeeManager.calcTradeFeeFromTransfer(acc.GetCoins(), &tran, keeper.engines)
	require.Equal(t, sdk.Coins{{"TNT", 5e2}}, fee.Tokens)

	// transferred in ABC-000
	tran = Transfer{
		inAsset:  "ABC-000",
		in:       1e3,
		outAsset: "BUSD-BD1",
		out:      100,
	}
	fee = keeper.FeeManager.calcTradeFeeFromTransfer(acc.GetCoins(), &tran, keeper.engines)
	require.Equal(t, sdk.Coins{{"TNT", 50}}, fee.Tokens)

	// transferred in XYZ-999
	tran = Transfer{
		inAsset:  "XYZ-999",
		in:       1e3,
		outAsset: "BUSD-BD1",
		out:      1e5,
	}
	fee = keeper.FeeManager.calcTradeFeeFromTransfer(acc.GetCoins(), &tran, keeper.engines)
	require.Equal(t, sdk.Coins{{"TNT", 5e4}}, fee.Tokens)

	// existing BUSD -> TNT trading pair
	ctx, am, keeper = setup()
	keeper.FeeManager.UpdateConfig(NewTestFeeConfig())
	keeper.AddEngine(dextype.NewTradingPair("BUSD-BD1", "TNT", 1e8))
	keeper.AddEngine(dextype.NewTradingPair("ABC-000", "BUSD-BD1", 1e7))
	keeper.AddEngine(dextype.NewTradingPair("BUSD-BD1", "XYZ-999", 1e6))

	// enough TNT, TNT will be collected
	_, acc = testutils.NewAccount(ctx, am, 1e10)

	// transferred in BUSD-BD1
	tran = Transfer{
		inAsset:  "BUSD-BD1",
		in:       1e4,
		outAsset: "ABC-000",
		out:      1e5,
	}
	fee = keeper.FeeManager.calcTradeFeeFromTransfer(acc.GetCoins(), &tran, keeper.engines)
	require.Equal(t, sdk.Coins{{"TNT", 5}}, fee.Tokens)

	// transferred in ABC-000
	tran = Transfer{
		inAsset:  "ABC-000",
		in:       1e6,
		outAsset: "BUSD-BD1",
		out:      1e5,
	}
	fee = keeper.FeeManager.calcTradeFeeFromTransfer(acc.GetCoins(), &tran, keeper.engines)
	require.Equal(t, sdk.Coins{{"TNT", 50}}, fee.Tokens)

	// transferred in XYZ-999
	tran = Transfer{
		inAsset:  "XYZ-999",
		in:       1e3,
		outAsset: "BUSD-BD1",
		out:      1e5,
	}
	fee = keeper.FeeManager.calcTradeFeeFromTransfer(acc.GetCoins(), &tran, keeper.engines)
	require.Equal(t, sdk.Coins{{"TNT", 50}}, fee.Tokens)
}

func TestFeeManager_CalcFixedFee_SupportBUSD(t *testing.T) {
	setChainVersion()
	defer resetChainVersion()
	ctx, am, keeper := setup()
	keeper.FeeManager.UpdateConfig(NewTestFeeConfig())
	keeper.SetBUSDSymbol("BUSD-BD1")

	// existing TNT -> BUSD trading pair
	_, acc := testutils.NewAccount(ctx, am, 0)
	keeper.AddEngine(dextype.NewTradingPair("TNT", "BUSD-BD1", 1e5))
	keeper.AddEngine(dextype.NewTradingPair("ABC-000", "BUSD-BD1", 1e7))
	keeper.AddEngine(dextype.NewTradingPair("BUSD-BD1", "XYZ-999", 1e6))

	// no enough TNT, the transferred-in asset will be collected
	// buy BUSD-BD1
	acc.SetCoins(sdk.Coins{{Denom: "BUSD-BD1", Amount: 1e4}})
	fee := keeper.FeeManager.CalcFixedFee(acc.GetCoins(), eventFullyExpire, "BUSD-BD1", keeper.engines)
	require.Equal(t, sdk.Coins{sdk.NewCoin("BUSD-BD1", 1e2)}, fee.Tokens)

	// buy ABC-000
	acc.SetCoins(sdk.Coins{{Denom: "ABC-000", Amount: 1e4}})
	fee = keeper.FeeManager.CalcFixedFee(acc.GetCoins(), eventFullyExpire, "ABC-000", keeper.engines)
	require.Equal(t, sdk.Coins{sdk.NewCoin("ABC-000", 1e3)}, fee.Tokens)

	// buy XYZ-999
	acc.SetCoins(sdk.Coins{{Denom: "XYZ-999", Amount: 1e4}})
	fee = keeper.FeeManager.CalcFixedFee(acc.GetCoins(), eventFullyExpire, "XYZ-999", keeper.engines)
	require.Equal(t, sdk.Coins{sdk.NewCoin("XYZ-999", 1)}, fee.Tokens)

	// existing BUSD -> TNT trading pair
	ctx, am, keeper = setup()
	keeper.FeeManager.UpdateConfig(NewTestFeeConfig())
	_, acc = testutils.NewAccount(ctx, am, 0)
	keeper.AddEngine(dextype.NewTradingPair("BUSD-BD1", "TNT", 1e9))
	keeper.AddEngine(dextype.NewTradingPair("ABC-000", "BUSD-BD1", 1e7))
	keeper.AddEngine(dextype.NewTradingPair("BUSD-BD1", "XYZ-999", 1e6))

	// no enough TNT, the transferred-in asset will be collected
	// buy BUSD-BD1
	acc.SetCoins(sdk.Coins{{Denom: "BUSD-BD1", Amount: 1e11}})
	fee = keeper.FeeManager.CalcFixedFee(acc.GetCoins(), eventFullyExpire, "BUSD-BD1", keeper.engines)
	require.Equal(t, sdk.Coins{sdk.NewCoin("BUSD-BD1", 1e4)}, fee.Tokens)

	// buy ABC-000
	acc.SetCoins(sdk.Coins{{Denom: "ABC-000", Amount: 1e10}})
	fee = keeper.FeeManager.CalcFixedFee(acc.GetCoins(), eventFullyExpire, "ABC-000", keeper.engines)
	require.Equal(t, sdk.Coins{sdk.NewCoin("ABC-000", 1e5)}, fee.Tokens)

	// buy XYZ-999
	acc.SetCoins(sdk.Coins{{Denom: "XYZ-999", Amount: 1e10}})
	fee = keeper.FeeManager.CalcFixedFee(acc.GetCoins(), eventFullyExpire, "XYZ-999", keeper.engines)
	require.Equal(t, sdk.Coins{sdk.NewCoin("XYZ-999", 1e2)}, fee.Tokens)
}
