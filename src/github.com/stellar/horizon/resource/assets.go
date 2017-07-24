package resource

import (

"github.com/stellar/horizon/db2/core"

"golang.org/x/net/context"
)

func (this *Assets) Populate( ctx context.Context, assetsSelling []core.Asset, assetsBuyings []core.Asset) (err error) {

	this.BuyingAssets = make([]Asset, len(assetsBuyings))
	this.SellingAssets = make([]Asset, len(assetsSelling))
	for i, asset := range assetsBuyings {
		s, _ := core.AssetFromDB(asset.AssetType, asset.AssetCode.String, asset.Issuer.String)
		this.BuyingAssets[i].Populate(ctx, s)
	}
	for i, asset := range assetsSelling {
		s, _ := core.AssetFromDB(asset.AssetType, asset.AssetCode.String, asset.Issuer.String)
		this.SellingAssets[i].Populate(ctx, s)
	}
	return
}