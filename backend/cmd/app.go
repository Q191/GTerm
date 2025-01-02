package cmd

import (
	"context"
	"github.com/MisakaTAT/GTerm/backend/initialize"
	"github.com/MisakaTAT/GTerm/backend/services"
	"github.com/google/wire"
	"go.uber.org/zap"
	"net/http"
)

var AppSet = wire.NewSet(wire.Struct(new(App), "*"))

type App struct {
	context          context.Context `wire:"-"`
	HTTPListenerPort *initialize.HTTPListenerPort
	Logger           *zap.Logger
	TerminalSrv      *services.TerminalSrv
	PreferencesSrv   *services.PreferencesSrv
	GroupSrv         *services.GroupSrv
	HostSrv          *services.HostSrv
}

func (a *App) Startup(ctx context.Context) {
	a.context = ctx
	http.Handle("/ws/terminal", http.HandlerFunc(a.TerminalSrv.Startup))
	// runtime.LogSetLogLevel(ctx, logger.INFO)
}

func (a *App) Bind() (services []any) {
	services = append(services, a.TerminalSrv)
	services = append(services, a.PreferencesSrv)
	services = append(services, a.GroupSrv)
	services = append(services, a.HostSrv)
	return
}
