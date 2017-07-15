package resource

import (
	//"encoding/base64"
	//"fmt"

	"github.com/stellar/horizon/db2/core"
	//"github.com/stellar/horizon/db2/history"
	//"github.com/stellar/horizon/httpx"
	//"github.com/stellar/horizon/render/hal"
	"golang.org/x/net/context"
)

func (this *Alias) Populate(ctx context.Context, row core.Alias) {
	this.AliasID = row.AliasID
}

// Populate fills out the resource's fields
func (this *Aliases) Populate(
	ctx context.Context,
	cl []core.Alias,
) (err error) {
	// populate data
	this.Aliases = make([]Alias, len(cl))
	for i, s := range cl {
		this.Aliases[i].Populate(ctx, s)
	}

	return
}