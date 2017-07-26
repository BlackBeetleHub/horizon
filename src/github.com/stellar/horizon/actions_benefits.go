package horizon

import (
	"github.com/stellar/horizon/render/hal"
	"github.com/stellar/horizon/db2/core"
	"github.com/stellar/horizon/paths"
	"github.com/stellar/horizon/resource"
)

type BenefitsShowAction struct {
	Action
	Query   paths.Exchange
	Records []paths.Path
	ListAssets []core.Asset
	Page    hal.BasePage
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

func (action *BenefitsShowAction) loadQuery() { //TODO: to make load better amount.
	action.Query.DestinationAmount = 100000000
	asset := action.ListAssets[3]
	distAsset :=action.ListAssets[4]
	println(asset.AssetCode.String)
	println(distAsset.AssetCode.String)
	action.Query.SourceAsset, action.Err = core.AssetFromDB(distAsset.AssetType, distAsset.AssetCode.String, distAsset.Issuer.String)
	if action.Err != nil {
		println("Error loadQuery")
		return;
	}
	action.Query.DestinationAsset, action.Err = core.AssetFromDB(asset.AssetType, asset.AssetCode.String, asset.Issuer.String)
	if action.Err != nil {
		println("Error loadQuery")
		return;
	}
}

func (action *BenefitsShowAction) loadAssets(){
	action.Err = action.CoreQ().AssetsForBuying(
		&action.ListAssets,
	)
}

func (action *BenefitsShowAction) loadResource(){
	action.Records, action.Err = action.App.benefits.Find(action.Query)
}

func (action *BenefitsShowAction) loadPage() {
	action.Page.Init()
	for _, p := range action.Records {
		var res resource.Path
		action.Err = res.PopulateBenefit(action.Ctx, action.Query, p)
		if action.Err != nil {
			return
		}
		action.Page.Add(res)
	}
}