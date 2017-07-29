package benefits

import (
	"github.com/stellar/horizon/paths"
	"github.com/stellar/horizon/db2/core"
	"github.com/stellar/horizon/assetspath"
	"strconv"
)

type Benefit struct {
	benefitChecker paths.BenefitsChecker
	PossibleExchanges []paths.CoreExchange
	BenefitExchanges []BenefitExchange
	q *core.Q
}

type BenefitExchange struct {
	To paths.Path
	Back paths.Path
	AmountBenefit int64
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
	benefit.Start()
	return nil
}

func (benefit *Benefit) Start () {
	benefit.BenefitExchanges = benefit.SearchBenefits()
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

func (benefit *Benefit) SearchBenefits () []BenefitExchange {
	var benefitExchanges []BenefitExchange

	pExcenges := &benefit.PossibleExchanges

	for i:=0; i <len(*pExcenges); i++ {
		res := benefit.SearchBenefitsInExchange((*pExcenges)[i].ToExchange())
		if len(res) != 0 {
			benefitExchanges = append(benefitExchanges, res...)
		}
		//benefitExchanges = append(benefitExchanges, res...)
	}
	return benefitExchanges
}

func (benefit *Benefit) SearchBenefitsInExchange(exchange paths.Exchange) []BenefitExchange {
	var result []BenefitExchange
	fronts,err := benefit.GetPathsFromExchange(exchange)
	if err!=nil {
		return result
	}
	backs, err := benefit.GetBackPathsFromExchange(exchange)
	if err != nil {
		return result
	}
	for i:=0; i< len(fronts); i++ {
		for t:=0; t< len(backs); t++ {
			isBenefit,err := benefit.isBenefitPaths(fronts[i],backs[t])
			if err != nil {
				// TODO: make validator
			}
			if isBenefit {
				result = append(result,
					BenefitExchange{ To:fronts[i], Back:backs[t], AmountBenefit:1 })
			}
		}
	}
	return result
}

func (benefit *Benefit) isBenefitPaths(front, back paths.Path) (bool, error) {

	maxDistFront, err := front.MaxCost()
	if (err !=nil || maxDistFront == 0){
		return false, err
	}

	maxSourceFront, err := front.Cost(maxDistFront)
	if (err != nil || maxSourceFront == 0){
		return false, err
	}

	maxDistBack, err := back.MaxCost()
	if (err != nil || maxDistBack == 0) {
		return false, err
	}

	maxSourceBack, err := back.Cost(maxDistBack)
	if ( err != nil || maxSourceBack == 0) {
		return false, err
	}

	if maxDistFront > maxSourceBack {
		maxDistFront = maxSourceBack
		maxSourceFront,err = front.Cost(maxDistFront)
		if (err != nil || maxSourceFront == 0) {
			return false, err
		}
		return (maxSourceFront < maxDistBack), err
	}

	if maxSourceBack > maxDistFront {
		maxDistBack = maxSourceFront
		maxSourceBack, err = back.Cost(maxDistBack)
		if (err !=nil || maxSourceFront == 0) {
			return false, err
		}
		return (maxSourceBack < maxDistFront), err
	}
	return false, err
}
