package paths

import (
	"github.com/stellar/go/xdr"
	"github.com/stellar/horizon/db2/core"
)

// Query is a query for paths
type Query struct {
	DestinationAddress string
	DestinationAsset   xdr.Asset
	DestinationAmount  xdr.Int64
	SourceAssets       []xdr.Asset
}

type Exchange struct {
	DestinationAsset   xdr.Asset
	SourceAsset		   xdr.Asset
}

type CoreExchange struct {
	Dest    core.Asset
	Source  core.Asset
}

func (exc *CoreExchange) ToExchange() Exchange{
	var exchange Exchange
	exchange.SourceAsset = exc.Source.ToXdrAsset()
	exchange.DestinationAsset = exc.Dest.ToXdrAsset()
	return exchange
}

// Path is the interface that represents a single result returned
// by a path finder.
type Path interface {
	Path() []xdr.Asset
	Source() xdr.Asset
	Destination() xdr.Asset
	// Cost returns an amount (which may be estimated), delimited in the Source assets
	// that is suitable for use as the `sendMax` field for a `PathPaymentOp` struct.
	Cost(amount xdr.Int64) (xdr.Int64, error)
	MaxCost() (result xdr.Int64, err error)
}

type Paths struct {
	Paths []Path
}
// Finder finds paths.
type Finder interface {
	Find(Query) ([]Path, error)
	FindFromExchange(exchange Exchange) ([]Path, error)
}