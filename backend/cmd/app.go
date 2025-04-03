package cmd

import (
	"context"
	"github.com/MisakaTAT/GTerm/backend/enums"
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
	ConnectionSrv    *services.ConnectionSrv
	MetadataSrv      *services.MetadataSrv
	CredentialSrv    *services.CredentialSrv
	WebsocketSrv     *services.WebsocketSrv
}

func (a *App) Startup(ctx context.Context) {
	a.context = ctx
	http.Handle("/ws/terminal", http.HandlerFunc(a.WebsocketSrv.TerminalHandle))
	// runtime.LogSetLogLevel(ctx, logger.INFO)
}

func (a *App) Bind() (services []any) {
	services = append(services, a.TerminalSrv)
	services = append(services, a.PreferencesSrv)
	services = append(services, a.GroupSrv)
	services = append(services, a.ConnectionSrv)
	services = append(services, a.MetadataSrv)
	services = append(services, a.CredentialSrv)
	return
}

func (a *App) Enums() (es []any) {
	es = append(es, enums.AuthMethodEnums)
	es = append(es, enums.ConnProtocolEnums)
	es = append(es, enums.TerminalTypeEnums)
	return
}
