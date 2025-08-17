//go:build wireinject

package cmd

import (
	"github.com/Q191/GTerm/backend/initialize"
	"github.com/Q191/GTerm/backend/services"
	"github.com/google/wire"
)

func NewApp() *App {
	wire.Build(
		AppSet,
		initialize.InitDefaultContext,
		initialize.InitDatabase,
		initialize.InitHTTPServer,
		initialize.InitLogger,
		services.ProviderSet,
	)
	return new(App)
}
