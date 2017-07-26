package assetspath

import (
	"github.com/stellar/horizon/db2/core"
	"github.com/stellar/horizon/log"
	"github.com/stellar/horizon/paths"
)

// Finder implements the paths.Finder interface and searchs for
// payment paths using a simple breadth first search of the offers table of a stellar-core.
//
// This implementation is not meant to be fast or to provide the lowest costs paths, but
// rather is meant to be a simple implementation that gives usable paths.
type BenefitsChecker struct {
	Q *core.Q
}

// ensure the struct is paths.Finder compliant
var _ paths.BenefitsChecker = &BenefitsChecker{}

// Find performs a path find with the provided query.
func (f *BenefitsChecker) Find(q paths.Exchange) (result []paths.Path, err error) {
	println("Yes, find benefitsChecker))")
	log.WithField("source_asset", q.SourceAsset).
		WithField("destination_asset", q.DestinationAsset).
		Info("Starting pathfind")

	s := &search{
		Exchange:  q,
		BenefitsChecker: f,
	}

	s.Init()
	s.Run()

	result, err = s.Results, s.Err

	log.WithField("found", len(s.Results)).
		WithField("err", s.Err).
		Info("Finished pathfind")
	return
}