package initialize

import (
	"github.com/wailsapp/wails/v2/pkg/logger"

	"github.com/google/wire"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

var InitLogger = wire.NewSet(ProvideLogger)

type Logger interface {
	Debug(msg string, args ...any)
	Info(msg string, args ...any)
	Warn(msg string, args ...any)
	Error(msg string, args ...any)
	Fatal(msg string, args ...any)
}

type LoggerWrapper struct {
	appContext *AppContext
}

func ProvideLogger(appContext *AppContext) Logger {
	return &LoggerWrapper{
		appContext: appContext,
	}
}

func (l *LoggerWrapper) SetLogLevel(level logger.LogLevel) {
	runtime.LogSetLogLevel(l.appContext.Context(), level)
}

func (l *LoggerWrapper) Debug(msg string, args ...any) {
	if len(args) == 0 {
		runtime.LogDebug(l.appContext.Context(), msg)
		return
	}
	runtime.LogDebugf(l.appContext.Context(), msg, args...)
}

func (l *LoggerWrapper) Info(msg string, args ...any) {
	if len(args) == 0 {
		runtime.LogInfo(l.appContext.Context(), msg)
		return
	}
	runtime.LogInfof(l.appContext.Context(), msg, args...)
}

func (l *LoggerWrapper) Warn(msg string, args ...any) {
	if len(args) == 0 {
		runtime.LogWarning(l.appContext.Context(), msg)
		return
	}
	runtime.LogWarningf(l.appContext.Context(), msg, args...)
}

func (l *LoggerWrapper) Error(msg string, args ...any) {
	if len(args) == 0 {
		runtime.LogError(l.appContext.Context(), msg)
		return
	}
	runtime.LogErrorf(l.appContext.Context(), msg, args...)
}

func (l *LoggerWrapper) Fatal(msg string, args ...any) {
	if len(args) == 0 {
		runtime.LogFatal(l.appContext.Context(), msg)
		return
	}
	runtime.LogFatalf(l.appContext.Context(), msg, args...)
}
