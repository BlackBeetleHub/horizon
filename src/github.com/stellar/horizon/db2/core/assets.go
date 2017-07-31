package core

import (
	sq "github.com/Masterminds/squirrel"
)

func (q *Q) AssetsForBuying(dest interface{}) error {
	return q.Select(dest, selectorAssetsForBuying)
}

func (q *Q) AssetsForSelling(dest interface{}) error{
	return q.Select(dest, selectorAssetsForSelling);
}

var selectorAssetsForSelling = sq.Select(
	"sellingassettype AS assettype",
	"sellingassetcode AS assetcode",
	"sellingissuer AS issuer",
).From("offers").GroupBy(
	"sellingassetcode",
	"sellingissuer",
	"sellingassettype",
)

var selectorAssetsForBuying = sq.Select(
	"buyingassettype AS assettype",
	"buyingassetcode AS assetcode",
	"buyingissuer AS issuer",
).From("offers").GroupBy(
	"buyingassetcode",
	"buyingissuer",
	"buyingassettype",
)


func (q *Q) MaxExchangeCounter(source, dist Asset) int64 {
	var maxForvard int64
	sql := sq.Select("SUM(amount) AS amounts").From("offers").
		Where("buyingassetcode = ?", source.AssetCode.String).
	Where("buyingassettype = ?", source.AssetType).
	Where("buyingissuer = ?", source.Issuer.String).
	Where("sellingassetcode = ?", dist.AssetCode.String).
	Where("sellingassettype = ?", dist.AssetType).
	Where("sellingissuer = ?", dist.Issuer.String)
	rows ,err := q.Query(sql)
	if err !=nil {
		println("error query")
		return 0
	}
	rows.Next()
	rows.Scan(&maxForvard)
	rows.Close()
	return  maxForvard
}

func (q *Q) MaxCountCanExchange(source, dist Asset) int64 {
	var maxSell,maxBuy int64
	maxSell = q.MaxExchangeCounter(source, dist)
	maxBuy = q.MaxExchangeCounter(dist, source)
	if maxBuy > maxSell {
		return  maxSell
	}
	return maxBuy
}