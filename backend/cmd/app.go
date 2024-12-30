package cmd

import (
	"context"
	"github.com/OpenToolkitLab/GTerm/backend/initialize"
	"github.com/OpenToolkitLab/GTerm/backend/services"
	"github.com/google/wire"
	"go.uber.org/zap"
	"net/http"
)

var AppSet = wire.NewSet(wire.Struct(new(App), "*"))

type App struct {
	context          context.Context `wire:"-"`
	Database         *initialize.Database
	HTTPListenerPort *initialize.HTTPListenerPort
	Logger           *zap.Logger
	TerminalSrv      *services.TerminalSrv
	PreferencesSrv   *services.PreferencesSrv
}

func (a *App) Startup(ctx context.Context) {
	a.context = ctx
	http.Handle("/ws/terminal", http.HandlerFunc(a.TerminalSrv.Startup))
	// runtime.LogSetLogLevel(ctx, logger.INFO)
}

func (a *App) Bind() (services []any) {
	services = append(services, a.TerminalSrv)
	services = append(services, a.PreferencesSrv)
	return
}
