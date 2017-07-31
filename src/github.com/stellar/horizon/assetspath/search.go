package assetspath

import (
	"github.com/stellar/go/xdr"
	"github.com/stellar/horizon/paths"
	smPath "github.com/stellar/horizon/simplepath"
)

type search struct {
	Exchange  		paths.Exchange
	BenefitsChecker *BenefitsChecker

	queue   []*smPath.PathNode
	targets map[string]bool
	visited map[string]bool

	Err     error
	Results []paths.Path
}

func (s *search) Init() {
	s.queue = []*smPath.PathNode{
		&smPath.PathNode{
			Asset: s.Exchange.DestinationAsset,
			Tail:  nil,
			Q:     s.BenefitsChecker.Q,
		},
	}

	s.targets = map[string]bool{}
	s.targets[s.Exchange.SourceAsset.String()] = true
	s.visited = map[string]bool{}
	s.Err = nil
	s.Results = nil
}

func (s *search) Run() {
	if s.Err != nil {
		return
	}

	for s.hasMore() {
		s.runOnce()
	}
}

func (s *search) pop() *smPath.PathNode {
	next := s.queue[0]
	s.queue = s.queue[1:]
	return next
}

func (s *search) hasMore() bool {
	if s.Err != nil {
		return false
	}

	if len(s.Results) > 40 {
		return false
	}

	return len(s.queue) > 0
}

func (s *search) isTarget(id string) bool {
	_, found := s.targets[id]
	return found
}

func (s *search) visit(id string) bool {
	if _, found := s.visited[id]; found {
		return false
	}

	s.visited[id] = true
	return true
}

func (s *search) runOnce() {
	cur := s.pop()
	id := cur.Asset.String()

	if s.isTarget(id) {
		s.Results = append(s.Results, cur)
	}

	if !s.visit(id) {
		return
	}

	if cur.Depth() > 7 {
		return
	}

	s.extendSearch(cur)
}

func (s *search) extendSearch(cur *smPath.PathNode) {
	var connected []xdr.Asset
	s.Err = s.BenefitsChecker.Q.ConnectedAssets(&connected, cur.Asset)
	if s.Err != nil {
		return
	}

	for _, a := range connected {
		newPath := &smPath.PathNode{
			Asset: a,
			Tail:  cur,
			Q:     s.BenefitsChecker.Q,
		}

		s.queue = append(s.queue, newPath)
	}
}
