package horizon

import (
	"github.com/stellar/horizon/render/hal"
	"github.com/stellar/horizon/db2/core"
	"github.com/stellar/horizon/paths"
	"github.com/stellar/horizon/resource"
	"github.com/stellar/horizon/benefits"
)

type BenefitsShowAction struct {
	Action
	Query   paths.Exchange
	PossibleExchanges []paths.CoreExchange
	Records []paths.Path
	ListAssets []core.Asset
	Page    hal.BasePage
	Bfts benefits.Benefit
}

func (action *BenefitsShowAction) JSON()  {
	action.Do(
		action.loadAssets,
		action.loadQuery,
		action.loadResource,
		action.loadPage,
		func(){
			hal.Render(action.W,action.Page)
		},
	)
}

func (action *BenefitsShowAction) loadQuery() {
	var asset, distAsset core.Asset
	for i:=0; i < len(action.ListAssets); i++ {
		if(action.ListAssets[i].AssetCode.String == "EURT"){
			asset = action.ListAssets[i]
		}
		if(action.ListAssets[i].AssetCode.String == "BTC"){
			distAsset = action.ListAssets[i]
		}
	}

	var tmp []paths.Path
	tmp, _ = action.App.benefits.Find(action.Query)
	if len(tmp) == 0 {
		println("Bad pathing")
	}
	println(asset.AssetCode.String)
	println(distAsset.AssetCode.String)
	//print(action.CoreQ().MaxExchangeCounter(asset,distAsset))
	action.Query.SourceAsset, action.Err = core.AssetFromDB(distAsset.AssetType, distAsset.AssetCode.String, distAsset.Issuer.String)
	if action.Err != nil {
		println(action.Err)
		return;
	}
	action.Query.DestinationAsset, action.Err = core.AssetFromDB(asset.AssetType, asset.AssetCode.String, asset.Issuer.String)
	if action.Err != nil {
		println(action.Err)
		return;
	}
}

func (action *BenefitsShowAction) loadAssets(){
	action.Err = action.CoreQ().AssetsForBuying(
		&action.ListAssets,
	)
	action.PossibleExchanges, action.Err = benefits.GeneratePossibleExchanges(action.ListAssets)
	action.Bfts.InitPossibleExchanges(action.ListAssets)
	action.Bfts.CheckValidExchanges(action.CoreQ())
	if action.Err != nil {
		println(action.Err)
		return
	}
}

func (action *BenefitsShowAction) loadResource(){
	action.Records, action.Err = action.App.benefits.Find(action.Query)
}

func (action *BenefitsShowAction) loadPage() {
	action.Page.Init()
	for _, p := range action.Records {
		var res resource.BasePath
		action.Err = res.PopulateBenefit(action.Ctx, action.Query, p)
		if action.Err != nil {
			return
		}
		action.Page.Add(res)
	}
}