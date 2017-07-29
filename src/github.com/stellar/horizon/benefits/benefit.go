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
	benefit.Start()
	return nil
}

func (benefit *Benefit) Start () {
	first := benefit.PossibleExchanges[12]
	exchange:= first.ToExchange()
	firstPaths, err := benefit.GetPathsFromExchange(exchange)
	if err != nil {
		println(err)
		println("GetPathsFromExchange error")
		return
	}
	backPaths, err := benefit.GetBackPathsFromExchange(exchange)

	if err != nil {
		println(err)
		println("GetBackPathsFromExchange error")
		return
	}

	from:= firstPaths[0]
	to:=backPaths[0]

	res , err :=benefit.isBenefit(from, to)
	if err != nil {
		println(err)
	}
	println(res)
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

func (benefit *Benefit) isBenefit(front, back paths.Path) (bool, error) {

	println("Front :")
	assets:=front.Path()
	println(front.Source().String())
	for i:=0;i < len(assets);i++ {
		println(assets[i].String())
	}
	println(front.Destination().String())

	println("Back :")
	println(back.Source().String())
	assets =back.Path()
	for i:=0;i < len(assets);i++ {
		println(assets[i].String())
	}
	println(back.Destination().String())
	maxDistFront, err := front.MaxCost()
	if (err !=nil || maxDistFront == 0){
		println("maxDistFront")
		return false, err
	}
	maxSourceFront, err := front.Cost(maxDistFront)
	if (err != nil || maxSourceFront == 0){
		println("maxSourceFront")
		return false, err
	}
	maxDistBack, err := back.MaxCost()
	if (err != nil || maxDistBack == 0) {
		println("maxDistBack")
		return false, err
	}
	maxSourceBack, err := back.Cost(maxDistBack)
	if ( err != nil || maxSourceBack == 0) {
		println("maxSourceBack")
		return false, err
	}
	print("MaxSourceFront: ")
	println(maxSourceFront)
	print("maxDistFront: ")
	println(maxDistFront)
	print("MaxSourceBack: ")
	println(maxSourceBack)
	print("MaxDistBack: ")
	println(maxDistBack)
	println("Return true always!!!!")

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
		print("MaxSourceFront: ")
		println(maxSourceFront)
		print("maxDistFront: ")
		println(maxDistFront)
		print("MaxSourceBack: ")
		println(maxSourceBack)
		print("MaxDistBack: ")
		println(maxDistBack)
		if (err !=nil || maxSourceFront == 0) {
			return false, err
		}
		return (maxSourceBack < maxDistFront), err
	}
	return false, err
}
