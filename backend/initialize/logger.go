package initialize

import (
	"context"

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
	ctx context.Context
}

func ProvideLogger() Logger {
	return &LoggerWrapper{}
}

func (l *LoggerWrapper) SetContext(ctx context.Context) {
	l.ctx = ctx
}

func (l *LoggerWrapper) SetLogLevel(level logger.LogLevel) {
	runtime.LogSetLogLevel(l.ctx, level)
}

func (l *LoggerWrapper) Debug(msg string, args ...any) {
	if len(args) == 0 {
		runtime.LogDebug(l.ctx, msg)
		return
	}
	runtime.LogDebugf(l.ctx, msg, args...)
}

func (l *LoggerWrapper) Info(msg string, args ...any) {
	if len(args) == 0 {
		runtime.LogInfo(l.ctx, msg)
		return
	}
	runtime.LogInfof(l.ctx, msg, args...)
}

func (l *LoggerWrapper) Warn(msg string, args ...any) {
	if len(args) == 0 {
		runtime.LogWarning(l.ctx, msg)
		return
	}
	runtime.LogWarningf(l.ctx, msg, args...)
}

func (l *LoggerWrapper) Error(msg string, args ...any) {
	if len(args) == 0 {
		runtime.LogError(l.ctx, msg)
		return
	}
	runtime.LogErrorf(l.ctx, msg, args...)
}

func (l *LoggerWrapper) Fatal(msg string, args ...any) {
	if len(args) == 0 {
		runtime.LogFatal(l.ctx, msg)
		return
	}
	runtime.LogFatalf(l.ctx, msg, args...)
}
