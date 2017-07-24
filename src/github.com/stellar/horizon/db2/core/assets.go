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
