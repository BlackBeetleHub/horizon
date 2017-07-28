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
	benefitChecker paths.BenefitsChecker
	PossibleExchanges []paths.CoreExchange
	PossiblePaths []Pathes
	ResultPaths []Pathes
	q *core.Q
}

func (benefit *Benefit) Init(q *core.Q) error {
	var coreAssests []core.Asset
	benefit.benefitChecker = &assetspath.BenefitsChecker {q}
	err := q.AssetsForBuying(&coreAssests)
	if err != nil {
		return err
	}
	benefit.InitPossibleExchanges(coreAssests)
	benefit.CheckValidExchanges()
	return nil
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

func (benefit *Benefit) CheckValidExchanges() {
	println("Before validate Exchanges" + strconv.Itoa(len(benefit.PossibleExchanges)))
	var validateResultWay []paths.CoreExchange
	for i:=0; i < len(benefit.PossibleExchanges); i++ {
		var exchange paths.Exchange
		short := benefit.PossibleExchanges[i]
		exchange.SourceAsset = short.Source.ToXdrAsset()
		exchange.DestinationAsset = short.Dest.ToXdrAsset()
		resBuy, _:= benefit.benefitChecker.Find(exchange)
		if len(resBuy) == 0 {
			continue;
		}
		exchange = SwapExchangeAssets(exchange)
		resSell, _:= benefit.benefitChecker.Find(exchange)
		if len(resSell) == 0 {
			continue;
		}
		validateResultWay = append(validateResultWay, short)
	}
	benefit.PossibleExchanges = validateResultWay
	println("After validate Exchanges" + strconv.Itoa(len(benefit.PossibleExchanges)))
}

func (benefit *Benefit) GetPathsFromExchange(exchange paths.Exchange) (result []paths.Path, err error){
	return benefit.benefitChecker.Find(exchange)
}

func (benefit *Benefit) GetBackPathsFromExchange(exchange paths.Exchange) (result []paths.Path, err error){
	reverseExchange := SwapExchangeAssets(exchange)
	return benefit.benefitChecker.Find(reverseExchange)
}

func (benefit *Benefit) SearchBenefitInPath() {

}