package benefits

import "github.com/stellar/horizon/paths"

type BenefitExchange struct {
	To paths.Path
	Back paths.Path
	Profit int64
}

