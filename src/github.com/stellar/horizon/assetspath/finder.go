package assetspath

import (
	"github.com/stellar/horizon/db2/core"
	"github.com/stellar/horizon/paths"
)

type BenefitsChecker struct {
	Q *core.Q
}

var _ paths.BenefitsChecker = &BenefitsChecker{}

func (f *BenefitsChecker) Find(q paths.Exchange) (result []paths.Path, err error) {
	s := &search{
		Exchange:  q,
		BenefitsChecker: f,
	}

	s.Init()
	s.Run()

	result, err = s.Results, s.Err

	return
}