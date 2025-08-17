package cmd

import (
	"context"
	"net/http"

	"github.com/wailsapp/wails/v2/pkg/logger"

	"github.com/Q191/GTerm/backend/enums"
	"github.com/Q191/GTerm/backend/initialize"
	"github.com/Q191/GTerm/backend/services"
	"github.com/google/wire"
)

var AppSet = wire.NewSet(wire.Struct(new(App), "*"))

type App struct {
	AppContext       *initialize.AppContext
	HTTPListenerPort *initialize.HTTPListenerPort
	Logger           initialize.Logger
	TerminalSrv      *services.TerminalSrv
	PreferencesSrv   *services.PreferencesSrv
	GroupSrv         *services.GroupSrv
	ConnectionSrv    *services.ConnectionSrv
	MetadataSrv      *services.MetadataSrv
	CredentialSrv    *services.CredentialSrv
	WebsocketSrv     *services.WebsocketSrv
	FileTransferSrv  *services.FileTransferSrv
}

func (a *App) Startup(ctx context.Context) {
	a.AppContext.SetContext(ctx)

	if log, ok := a.Logger.(*initialize.LoggerWrapper); ok {
		log.SetLogLevel(logger.DEBUG)
	}

	http.Handle("/ws/terminal", http.HandlerFunc(a.WebsocketSrv.TerminalHandle))
}

func (a *App) Bind() (bd []any) {
	bd = append(bd, a.TerminalSrv)
	bd = append(bd, a.PreferencesSrv)
	bd = append(bd, a.GroupSrv)
	bd = append(bd, a.ConnectionSrv)
	bd = append(bd, a.MetadataSrv)
	bd = append(bd, a.CredentialSrv)
	bd = append(bd, a.FileTransferSrv)
	return
}

func (a *App) Enums() (es []any) {
	es = append(es, enums.AuthMethodEnums)
	es = append(es, enums.ConnProtocolEnums)
	es = append(es, enums.TerminalTypeEnums)
	es = append(es, enums.FileTransferTaskStateEnums)
	return
}
