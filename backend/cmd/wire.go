//go:build wireinject

package cmd

import (
	"github.com/OpenToolkitLab/GTerm/backend/initialize"
	"github.com/OpenToolkitLab/GTerm/backend/services"
	"github.com/google/wire"
)

func NewApp() *App {
	wire.Build(
		AppSet,
		initialize.InitDatabase,
		initialize.InitHTTPServer,
		initialize.InitZap,
		services.SetProvider,
	)
	return new(App)
}
