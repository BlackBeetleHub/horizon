package benefits

import (
	"github.com/stellar/horizon/paths"
	"github.com/stellar/horizon/db2/core"
	"github.com/stellar/horizon/assetspath"
	"strconv"
)

type Pathes struct{
	Paths []paths.Path
}

type Benefit struct {

	PossibleExchanges []paths.CoreExchange
	PossiblePaths []Pathes
	ResultPaths []Pathes
}

func (benefit *Benefit) InitPossibleExchanges (listBuying []core.Asset) {
	var result []paths.CoreExchange
	for i := 0; i < len(listBuying); i++ {
		for t:= i + 1; t < len(listBuying); t++ {
			result = append(result, paths.CoreExchange{listBuying[i], listBuying[t]})
		}
	}
	if len(result) == 0 {
		panic("PossibleExchanges is zero!")
		return;
	}
	benefit.PossibleExchanges = result
}

func SwapExchangeAssets(exchange paths.Exchange) paths.Exchange {
	return paths.Exchange{exchange.SourceAsset, exchange.DestinationAsset}
}

func (benefit *Benefit) CheckValidExchanges(q *core.Q) {
	println("Before validate Exchanges" + strconv.Itoa(len(benefit.PossibleExchanges)))
	mBenefits := &assetspath.BenefitsChecker{q}
	var validateResultWay []paths.CoreExchange
	for i:=0; i < len(benefit.PossibleExchanges); i++ {
		var exchange paths.Exchange
		short := benefit.PossibleExchanges[i]
		exchange.SourceAsset = short.Source.ToXdrAsset()
		exchange.DestinationAsset = short.Dest.ToXdrAsset()
		resBuy, _:= mBenefits.Find(exchange)
		if len(resBuy) == 0 {
			continue;
		}
		exchange = SwapExchangeAssets(exchange)
		resSell, _:= mBenefits.Find(exchange)
		if len(resSell) == 0 {
			continue;
		}
		validateResultWay = append(validateResultWay, short)
	}
	benefit.PossibleExchanges = validateResultWay
	println("After validate Exchanges" + strconv.Itoa(len(benefit.PossibleExchanges)))
}

func GetAmountForOperation() int64 {
	return  1
}

func GenerateBackPaths() {

}

func GenerateForwardPaths() {

}

func MaxCostForwardPath() {

}

func MaxCostBackPath(){

}

func (benefit *Benefit) SearchBenefitInPath() {

}