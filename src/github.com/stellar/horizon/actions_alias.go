package horizon

import (
	"github.com/stellar/horizon/db2/core"
	"github.com/stellar/horizon/render/hal"
	"github.com/stellar/horizon/render/sse"
	"github.com/stellar/horizon/resource"
)

// DataShowAction renders a account summary found by its address.
type AliasShowAction struct {
	Action
	Address 	 	string
	AliasAddress 	string
	Aliases 		[]core.Alias
	Resource		resource.Aliases
}

// JSON is a method for actions.JSON
func (action *AliasShowAction) JSON() {
	action.Do(
		action.loadParams,
		action.loadRecord,
		action.loadResource,
		func() {
			hal.Render(action.W, action.Resource)
		},
	)
}

	// SSE is a method for actions.SSE
func (action *AliasShowAction) SSE(stream sse.Stream) {
	action.Do(
		action.loadParams,
		action.loadRecord,
		func() {
			stream.Send(sse.Event{Data: action.Aliases})
		},
	)
}

func (action *AliasShowAction) loadParams() {
	action.Address = action.GetString("account_id")
}

func (action *AliasShowAction) loadRecord() {
	action.Err = action.CoreQ().
		AliasesByAddress(&action.Aliases, action.Address)
	if action.Err != nil {
		return
	}
}

func (action *AliasShowAction) loadResource() {
	action.Err = action.Resource.Populate(
		action.Ctx,
		action.Aliases,
	)
}
