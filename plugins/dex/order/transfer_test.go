package order

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestTradeTransfers_Sort(t *testing.T) {
	e := TradeTransfers{
		{inAsset: "ABC", outAsset: "CE", Oid: "1"},
		{inAsset: "ABC", outAsset: "BTC", Oid: "2"},
		{inAsset: "XYZ", outAsset: "BTC", Oid: "3"},
		{inAsset: "XYZ", outAsset: "CE", Oid: "4"},
		{inAsset: "ABC", outAsset: "XYZ", Oid: "5"},
		{inAsset: "BTC", outAsset: "CE", Oid: "6"},
		{inAsset: "CE", outAsset: "BTC", Oid: "7"},
		{inAsset: "CE", outAsset: "ABC", Oid: "8"},
		{inAsset: "ABC", outAsset: "CE", Oid: "9"},
		{inAsset: "ABC", outAsset: "BTC", Oid: "10"},
	}
	e.Sort()
	require.Equal(t, TradeTransfers{
		{inAsset: "CE", outAsset: "ABC", Oid: "8"},
		{inAsset: "CE", outAsset: "BTC", Oid: "7"},
		{inAsset: "ABC", outAsset: "CE", Oid: "1"},
		{inAsset: "ABC", outAsset: "CE", Oid: "9"},
		{inAsset: "BTC", outAsset: "CE", Oid: "6"},
		{inAsset: "XYZ", outAsset: "CE", Oid: "4"},
		{inAsset: "ABC", outAsset: "BTC", Oid: "2"},
		{inAsset: "ABC", outAsset: "BTC", Oid: "10"},
		{inAsset: "ABC", outAsset: "XYZ", Oid: "5"},
		{inAsset: "XYZ", outAsset: "BTC", Oid: "3"},
	}, e)
}

func TestExpireTransfers_Sort(t *testing.T) {
	e := ExpireTransfers{
		{inAsset: "ABC", Symbol: "ABC_CE", Oid: "1"},
		{inAsset: "ABC", Symbol: "ABC_BTC", Oid: "2"},
		{inAsset: "XYZ", Symbol: "XYZ_BTC", Oid: "3"},
		{inAsset: "XYZ", Symbol: "XYZ_CE", Oid: "4"},
		{inAsset: "ABC", Symbol: "ABC_XYZ", Oid: "5"},
		{inAsset: "BTC", Symbol: "CE_BTC", Oid: "6"},
		{inAsset: "CE", Symbol: "CE_BTC", Oid: "7"},
		{inAsset: "CE", Symbol: "ABC_CE", Oid: "8"},
		{inAsset: "ABC", Symbol: "ABC_CE", Oid: "9"},
		{inAsset: "ABC", Symbol: "ABC_BTC", Oid: "10"},
	}
	e.Sort()
	require.Equal(t, ExpireTransfers{
		{inAsset: "CE", Symbol: "ABC_CE", Oid: "8"},
		{inAsset: "CE", Symbol: "CE_BTC", Oid: "7"},
		{inAsset: "ABC", Symbol: "ABC_CE", Oid: "1"},
		{inAsset: "ABC", Symbol: "ABC_CE", Oid: "9"},
		{inAsset: "ABC", Symbol: "ABC_BTC", Oid: "2"},
		{inAsset: "ABC", Symbol: "ABC_BTC", Oid: "10"},
		{inAsset: "ABC", Symbol: "ABC_XYZ", Oid: "5"},
		{inAsset: "BTC", Symbol: "CE_BTC", Oid: "6"},
		{inAsset: "XYZ", Symbol: "XYZ_CE", Oid: "4"},
		{inAsset: "XYZ", Symbol: "XYZ_BTC", Oid: "3"},
	}, e)
}
