package horizon

import (
	"github.com/stellar/horizon/simplepath"
	"github.com/stellar/horizon/assetspath"
)

func initPathFinding(app *App) {
	app.paths = &simplepath.Finder{app.CoreQ()}
}

func initBenefitChecker(app *App){
	app.benefits = &assetspath.BenefitsChecker{app.CoreQ()}
}

func init() {
	appInit.Add("path-finder", initPathFinding, "app-context", "log", "core-db")
	appInit.Add("benefit-checker", initBenefitChecker, "app-context", "log", "core-db")
}