package simplepath

import (
	"github.com/stellar/go/xdr"
	"github.com/stellar/horizon/paths"
)

// Search represents a single query against the simple finder.  It provides
// a place to store the results of the query, mostly for the purposes of code
// clarity.
//
// The Search struct is used as follows:
//
// 1.  Create an instance, ensuring the Query and Finder fields are set
// 2.  Call Init() to populate dependent fields in the struct with their initial values
// 3.  Call Run() to perform the Search.
//
type Search struct {
	Query  paths.Query
	Finder *Finder

	// Fields below are initialized by a call to Init() after
	// setting the fields above
	queue   []*PathNode
	targets map[string]bool
	visited map[string]bool

	//This fields below are initialized after the Search is run
	Err     error
	Results []paths.Path
}

// Init initialized the Search, setting fields on the struct used to
// hold state needed during the actual Search.
func (s *Search) Init() {
	s.queue = []*PathNode{
		&PathNode{
			Asset: s.Query.DestinationAsset,
			Tail:  nil,
			Q:     s.Finder.Q,
		},
	}
	println("check simplepath work")
	// build a map of asset's string representation to check if a given node
	// is one of the targets for our Search.  Unfortunately, xdr.Asset is not suitable
	// for use as a map key, and so we use its string representation.
	s.targets = map[string]bool{}
	for _, a := range s.Query.SourceAssets {
		s.targets[a.String()] = true
	}

	s.visited = map[string]bool{}
	s.Err = nil
	s.Results = nil
}

// Run triggers the Search, which will populate the Results and Err
// field for the Search after completion.
func (s *Search) Run() {
	if s.Err != nil {
		return
	}

	for s.hasMore() {
		s.runOnce()
	}
}

// pop removes the head from the Search queue, returning it to the caller
func (s *Search) pop() *PathNode {
	next := s.queue[0]
	s.queue = s.queue[1:]
	return next
}

// returns false if the Search should stop.
func (s *Search) hasMore() bool {
	if s.Err != nil {
		return false
	}

	if len(s.Results) > 40 {
		return false
	}

	return len(s.queue) > 0
}

// isTarget returns true if the asset id provided is one of the targets
// for this Search (i.e. one of the requesting account's trusted assets)
func (s *Search) isTarget(id string) bool {
	_, found := s.targets[id]
	return found
}

// visit returns true if the asset id provided has not been
// visited on this Search, after marking the id as visited
func (s *Search) visit(id string) bool {
	if _, found := s.visited[id]; found {
		return false
	}

	s.visited[id] = true
	return true
}

// runOnce processes the head of the Search queue, findings results
// and extending the Search as necessary.
func (s *Search) runOnce() {
	cur := s.pop()
	id := cur.Asset.String()

	if s.isTarget(id) {
		s.Results = append(s.Results, cur)
	}

	if !s.visit(id) {
		return
	}

	// A PathPaymentOp's path cannot be over 5 elements in length, and so
	// we abort our Search if the current linked list is over 7 (since the list
	// includes both source and destination in addition to the path)
	if cur.Depth() > 7 {
		return
	}

	s.extendSearch(cur)

}

func (s *Search) extendSearch(cur *PathNode) {
	// find connected assets
	var connected []xdr.Asset
	s.Err = s.Finder.Q.ConnectedAssets(&connected, cur.Asset)
	if s.Err != nil {
		return
	}

	for _, a := range connected {
		newPath := &PathNode{
			Asset: a,
			Tail:  cur,
			Q:     s.Finder.Q,
		}

		var hasEnough bool
		hasEnough, s.Err = s.hasEnoughDepth(newPath)
		if s.Err != nil {
			return
		}

		if !hasEnough {
			continue
		}

		s.queue = append(s.queue, newPath)
	}
}

func (s *Search) hasEnoughDepth(path *PathNode) (bool, error) {
	_, err := path.Cost(s.Query.DestinationAmount)
	if err == ErrNotEnough {
		return false, nil
	}
	return true, err
}
