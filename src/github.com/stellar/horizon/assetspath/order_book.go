package assetspath


import (
	"errors"
	"math/big"

	sq "github.com/Masterminds/squirrel"
	"github.com/stellar/go/xdr"
	"github.com/stellar/horizon/assets"
	"github.com/stellar/horizon/db2/core"
)

// ErrNotEnough represents an error that occurs when pricing a trade on an
// orderbook.  This error occurs when the orderbook cannot fulfill the
// requested amount.
var ErrNotEnough = errors.New("not enough depth")

type orderBook struct {
	Selling xdr.Asset
	Buying  xdr.Asset
	Q       *core.Q
}


func (ob *orderBook) MaxReciveCount (source xdr.Asset, sourceAmount xdr.Int64) (result xdr.Int64, err error) {

	var tmpSourceAmount int64
	var canBuy int64
	canBuy = 0
	tmpSourceAmount = int64(sourceAmount)

	sql, sqlBuildError := ob.GetSelectBuilderForCost(source)
	inverted := assets.Equals(source, ob.Buying)
	if sqlBuildError != nil {
		return
	}

	rows, err := ob.Q.Query(sql)
	if err != nil {
		return
	}
	defer rows.Close()

	for rows.Next() {
		// load data from the row
		var available, pricen, priced, offerid int64
		if inverted {
			err = rows.Scan(&available, &priced, &pricen, &offerid)
			available = mul(available, pricen, priced)
		} else {
			err = rows.Scan(&available, &pricen, &priced, &offerid)
		}

		if err != nil {
			break
		}
		av := mul(available, priced, pricen)
		if av > tmpSourceAmount {
			canBuy+= mul(tmpSourceAmount, priced, pricen)
			return
		}
		canBuy += mul(available, priced, pricen)
		tmpSourceAmount -= av
	}

	return xdr.Int64(canBuy), nil
}

func (ob *orderBook) MaxAvailebleCost(source xdr.Asset) (result xdr.Int64, err error) {
	// load offers from the two assets

	sql, sqlBuildError := ob.GetSelectBuilderForCost(source)
	inverted := assets.Equals(source, ob.Buying)
	if sqlBuildError != nil {
		return
	}

	rows, err := ob.Q.Query(sql)
	if err != nil {
		return
	}
	defer rows.Close()

	var resMax int64

	for rows.Next() {
		// load data from the row
		var available, pricen, priced, offerid int64
		if inverted {
			err = rows.Scan(&available, &priced, &pricen, &offerid)
			available = mul(available, pricen, priced)
		} else {
			err = rows.Scan(&available, &pricen, &priced, &offerid)
		}
		if err != nil {
			break
		}
		resMax += available
	}

	return xdr.Int64(resMax), nil
}

func (ob *orderBook)GetSelectBuilderForCost(source xdr.Asset) (sq.SelectBuilder, error) {
	var (
		// selling/buying types
		st, bt xdr.AssetType
		// selling/buying codes
		sc, bc string
		// selling/buying issuers
		si, bi string
	)

	err := ob.Selling.Extract(&st, &sc, &si)

	err = ob.Buying.Extract(&bt, &bc, &bi)


	sql := sq.
	Select("amount", "pricen", "priced", "offerid").
		From("offers").
		Where(sq.Eq{
		"sellingassettype":               st,
		"COALESCE(sellingassetcode, '')": sc,
		"COALESCE(sellingissuer, '')":    si}).
		Where(sq.Eq{
		"buyingassettype":               bt,
		"COALESCE(buyingassetcode, '')": bc,
		"COALESCE(buyingissuer, '')":    bi})

	inverted := assets.Equals(source, ob.Buying)

	if !inverted {
		sql = sql.OrderBy("price ASC")
	} else {
		sql = sql.OrderBy("price DESC")
	}

	return sql, err
}

func (ob *orderBook) Cost(source xdr.Asset, sourceAmount xdr.Int64) (result xdr.Int64, err error) {
	// load offers from the two assets

	sql, sqlBuildError := ob.GetSelectBuilderForCost(source)
	inverted := assets.Equals(source, ob.Buying)
	if sqlBuildError != nil {
		return
	}

	rows, err := ob.Q.Query(sql)
	if err != nil {
		return
	}
	defer rows.Close()

	var (
		needed = int64(sourceAmount)
		cost   int64
	)

	check, _ := ob.MaxAvailebleCost(source)
	if check > sourceAmount {
		println("True check > sourceAmount")
	}
	for rows.Next() {
		// load data from the row
		var available, pricen, priced, offerid int64
		if inverted {
			err = rows.Scan(&available, &priced, &pricen, &offerid)
			available = mul(available, pricen, priced)
		} else {
			err = rows.Scan(&available, &pricen, &priced, &offerid)
		}
		if err != nil {
			return
		}

		if available >= needed {
			cost += mul(needed, pricen, priced)
			result = xdr.Int64(cost)
			return
		}

		cost += mul(available, pricen, priced)
		needed -= available
	}

	err = ErrNotEnough
	return
}

// mul multiplies the input amount by the input price
func mul(amount int64, pricen int64, priced int64) int64 {
	var r, n, d big.Int

	r.SetInt64(amount)
	n.SetInt64(pricen)
	d.SetInt64(priced)

	r.Mul(&r, &n)
	r.Quo(&r, &d)
	return r.Int64()
}
