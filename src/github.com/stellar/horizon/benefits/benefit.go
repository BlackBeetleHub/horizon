package benefits

import (
	"github.com/stellar/horizon/paths"
	"github.com/stellar/horizon/db2/core"
	"github.com/stellar/horizon/simplepath"
)

type Benefits struct {
	benefitChecker    paths.Finder
	PossibleExchanges []paths.CoreExchange
	BenefitExchanges  []BenefitExchange
	q                 *core.Q
}

func (benefit *Benefits) Init(q *core.Q) error {
	var coreAssests []core.Asset
	benefit.benefitChecker = &simplepath.Finder{q}
	err := q.AssetsForBuying(&coreAssests)
	if err != nil {
		return err
	}
	benefit.InitPossibleExchanges(coreAssests)
	benefit.CheckValidExchanges()
	return err
}

func (benefit *Benefits) Start() {
	benefit.BenefitExchanges = benefit.SearchBenefits()
}

// Create all possible exchanges without back.
// Exemple: if we have 6 assets - we'll give 15 Exchanges.
// For check result: C (n,k) where k = 2 , n = len([]core.Asset); C (6,2) = 15.
func (benefit *Benefits) InitPossibleExchanges(listBuying []core.Asset) {
	var result []paths.CoreExchange
	for i := 0; i < len(listBuying); i++ {
		for t := i + 1; t < len(listBuying); t++ {
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

// Check possible Exchanges. It deletes exchange if it doesn't find path forwar or back.
// Not all assets possible get back.
func (benefit *Benefits) CheckValidExchanges() {
	var validateResultWay []paths.CoreExchange
	for i := 0; i < len(benefit.PossibleExchanges); i++ {
		if benefit.ChackValidExchange(&benefit.PossibleExchanges[i]) {
			validateResultWay = append(validateResultWay, benefit.PossibleExchanges[i])
		}
	}
	benefit.PossibleExchanges = validateResultWay
}

func (benefit *Benefits) ChackValidExchange(ce *paths.CoreExchange) bool {
	exchange:= ce.ToExchange()
	resBuy, _ := benefit.benefitChecker.FindFromExchange(exchange)
	if len(resBuy) == 0 {
		return false
	}
	exchange = SwapExchangeAssets(exchange)
	resSell, _ := benefit.benefitChecker.FindFromExchange(exchange)
	if len(resSell) == 0 {
		return false
	}
	return true
}

func (benefit *Benefits) GetPathsFromExchange(exchange paths.Exchange) (result []paths.Path, err error) {
	return benefit.benefitChecker.FindFromExchange(exchange)
}

func (benefit *Benefits) GetBackPathsFromExchange(exchange paths.Exchange) (result []paths.Path, err error) {
	reverseExchange := SwapExchangeAssets(exchange)
	return benefit.benefitChecker.FindFromExchange(reverseExchange)
}

func (benefit *Benefits) SearchBenefits() []BenefitExchange {
	var benefitExchanges []BenefitExchange

	pExcenges := &benefit.PossibleExchanges

	for i := 0; i < len(*pExcenges); i++ {
		res := benefit.SearchBenefitsInExchange((*pExcenges)[i].ToExchange())
		if len(res) != 0 {
			benefitExchanges = append(benefitExchanges, res...)
		}
	}
	return benefitExchanges
}

func (benefit *Benefits) SearchBenefitsInExchange(exchange paths.Exchange) []BenefitExchange {
	var result []BenefitExchange
	fronts, err := benefit.GetPathsFromExchange(exchange)
	if err != nil {
		return result
	}
	backs, err := benefit.GetBackPathsFromExchange(exchange)
	if err != nil {
		return result
	}
	for i := 0; i < len(fronts); i++ {
		for t := 0; t < len(backs); t++ {
			isBenefit, profit, err := benefit.isBenefitPaths(fronts[i], backs[t])
			if err != nil {
				continue
			}
			if isBenefit {
				result = append(result,
					BenefitExchange{To: fronts[i], Back: backs[t], Profit: profit })
			}
		}
	}
	return result
}

func (benefit *Benefits) isBenefitPaths(front, back paths.Path) (bool, int64, error) {

	maxDistFront, err := front.MaxCost()
	if (err != nil || maxDistFront == 0) {
		return false, 0, err
	}

	maxSourceFront, err := front.Cost(maxDistFront)
	if (err != nil || maxSourceFront == 0) {
		return false, 0, err
	}

	maxDistBack, err := back.MaxCost()
	if (err != nil || maxDistBack == 0) {
		return false, 0, err
	}

	maxSourceBack, err := back.Cost(maxDistBack)
	if ( err != nil || maxSourceBack == 0) {
		return false, 0, err
	}

	if maxDistFront > maxSourceBack {
		maxDistFront = maxSourceBack
		maxSourceFront, err = front.Cost(maxDistFront)
		if (err != nil || maxSourceFront == 0) {
			return false, 0, err
		}
		if (maxSourceFront < maxDistBack) {
			return true, int64(maxDistBack - maxSourceFront), err
		}
	}

	if (maxSourceBack > maxDistFront && maxSourceFront > maxDistBack) {
		return false, 0, err
	}

	if (maxSourceBack > maxDistFront) {
		maxDistBack = maxSourceFront
		maxSourceBack, err = back.Cost(maxDistBack)
		if (err != nil || maxSourceFront == 0) {
			return false,0 , err
		}
		if (maxSourceBack < maxDistFront) {
			maxDistFront = maxSourceBack
			maxSourceFront, err = front.Cost(maxDistFront)
			if err!=nil {
				return false,0 , err
			}
			if (maxSourceFront < maxDistBack) {
				return true, int64(maxDistBack - maxSourceFront), err
			}
		}
	}
	return false,0, err
}
