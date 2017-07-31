package resource

import (
	"github.com/stellar/horizon/paths"
	"context"
	"github.com/stellar/horizon/benefits"
	"strconv"
)

func (this *BasePath) PopulateBasePath(ctx context.Context, p paths.Path) (err error) {
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

func (this *BasePath) PopulateBenefit(ctx context.Context, q paths.Exchange, p paths.Path) (err error) {
	this.PopulateBasePath(ctx, p)
	return
}

func (this BasePath) PagingToken() string {
	return ""
}

func (this *BenefitExchange) Populate(ctx context.Context, bp benefits.BenefitExchange) (err error) {
	err = this.FromTo.PopulateBasePath(ctx, bp.To)
	if err != nil {
		return
	}
	err = this.ToFrom.PopulateBasePath(ctx, bp.Back)
	this.Profit = strconv.FormatInt(bp.Profit,10)
	return
}

func (this BenefitExchange) PagingToken() string {
	return ""
}