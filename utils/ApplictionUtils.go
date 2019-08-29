package utils

import (
	"github.com/goplume/core/configuration"
	"github.com/goplume/core/utils/stereotypes"
)

func Launch(app stereotypes.Application) {
	configuration.BootstrapConfiguration()
	app.InitApplication()
	app.RunApplication()

}
