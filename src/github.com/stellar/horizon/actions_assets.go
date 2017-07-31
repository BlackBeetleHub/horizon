package horizon

import (
	"github.com/stellar/horizon/resource"
	"github.com/stellar/horizon/render/hal"
	"github.com/stellar/horizon/db2/core"
)

// PathIndexAction provides path finding
type AssetsShowAction struct {
	Action
	Resource resource.Assets
	ResourceBuy []core.Asset
	ResourceSell []core.Asset
}

// JSON implements actions.JSON
func (action *AssetsShowAction) JSON() {
	action.Do(
		action.loadAssets,
		action.loadResource,
		func() {
			hal.Render(action.W, action.Resource)
		},
	)
}

func (action *AssetsShowAction) loadAssets() {
	action.Err = action.CoreQ().AssetsForSelling(
		&action.ResourceSell,
	)

	if action.Err != nil {
		println("Not load Assets for Selling")
		return
	}

	action.Err = action.CoreQ().AssetsForBuying(
		&action.ResourceBuy,
	)

	if action.Err != nil {
		println("Not load Assets for Buying")
		return
	}
}

func (action *AssetsShowAction) loadResource() {
	action.Err = action.Resource.Populate(action.Ctx , action.ResourceBuy, action.ResourceSell)
	if action.Err != nil {
		println("Populate assets failed")
		return;
	}
}

