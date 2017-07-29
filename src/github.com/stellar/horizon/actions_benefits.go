package horizon

import (
	"github.com/stellar/horizon/render/hal"
	"github.com/stellar/horizon/resource"
	"github.com/stellar/horizon/benefits"
)

type BenefitsShowAction struct {
	Action
	Page    hal.BasePage
	Records []benefits.BenefitExchange
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
}

func (action *BenefitsShowAction) loadAssets(){
	action.Err = action.Bfts.Init(action.CoreQ())
	if action.Err != nil {
		println(action.Err)
		return
	}
}

func (action *BenefitsShowAction) loadResource(){
	action.Records = action.Bfts.BenefitExchanges
}

func (action *BenefitsShowAction) loadPage() {
	action.Page.Init()
	for _, p := range action.Records {
		var res resource.BenefitExchange
		action.Err = res.Populate(action.Ctx, p)
		if action.Err != nil {
			return
		}
		action.Page.Add(res)
	}
}