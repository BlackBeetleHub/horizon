package benefits

import (
	"github.com/stellar/horizon/db2/core"
	"github.com/stellar/horizon/paths"
	"github.com/pkg/errors"
)

func GeneratePossibleExchanges(listBuying []core.Asset) ([]paths.CoreExchange, error) {
	var result []paths.CoreExchange
	for i := 0; i < len(listBuying); i++ {
		for t:= i + 1; t < len(listBuying); t++ {
			result = append(result, paths.CoreExchange{listBuying[i], listBuying[t]})
		}
	}
	if len(result) != 0 {
		return result, nil
	}
	return result, errors.New("Error, the count of possible exchanges is zero.")
}






