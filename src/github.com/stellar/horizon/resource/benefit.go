package resource

import (
	"github.com/stellar/horizon/paths"
	"context"
)

func (this *BasePath) PopulateBenefit(ctx context.Context, q paths.Exchange, p paths.Path) (err error) {

	err = p.Source().Extract(
		&this.SourceAssetType,
		&this.SourceAssetCode,
		&this.SourceAssetIssuer)

	if err != nil {
		return
	}

	err = p.Destination().Extract(
		&this.DestinationAssetType,
		&this.DestinationAssetCode,
		&this.DestinationAssetIssuer)

	if err != nil {
		return
	}

	path := p.Path()

	this.Path = make([]Asset, len(path))

	for i, a := range path {
		err = a.Extract(
			&this.Path[i].Type,
			&this.Path[i].Code,
			&this.Path[i].Issuer)
		if err != nil {
			return
		}
	}

	return
}

func (this BasePath) PagingToken() string {
	return ""
}

func (this *BenefitExchange)