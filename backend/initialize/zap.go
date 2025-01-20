package initialize

import (
	"fmt"
	"github.com/MisakaTAT/GTerm/backend/consts"
	rotateLogs "github.com/lestrrat-go/file-rotatelogs"
	"github.com/vrischmann/userdir"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"io"
	"os"
	"path"
	"time"
)

func InitZap() *zap.Logger {
	logPath := path.Join(userdir.GetConfigHome(), consts.ApplicationName, fmt.Sprintf("%s.log", consts.ApplicationName))
	config := zapcore.EncoderConfig{
		MessageKey:   "msg",
		LevelKey:     "level",
		TimeKey:      "ts",
		CallerKey:    "file",
		LineEnding:   zapcore.DefaultLineEnding,
		EncodeLevel:  zapcore.CapitalColorLevelEncoder,
		EncodeCaller: zapcore.ShortCallerEncoder,
		EncodeTime:   customEncodeTime,
		EncodeDuration: func(d time.Duration, enc zapcore.PrimitiveArrayEncoder) {
			enc.AppendInt64(int64(d) / 1000000)
		},
	}

	// TODO: A required privilege is not held by the client.
	// Windows 下创建日志文件会有权限问题，之后再看看怎么解决。
	logWriter := zapcore.AddSync(writer(logPath))
	core := zapcore.NewTee(
		zapcore.NewCore(zapcore.NewConsoleEncoder(config), logWriter, zap.InfoLevel),
		zapcore.NewCore(zapcore.NewConsoleEncoder(config), zapcore.NewMultiWriteSyncer(zapcore.AddSync(os.Stdout)), zap.InfoLevel),
	)
	return zap.New(core, zap.AddCaller(), zap.AddStacktrace(zap.WarnLevel))
}

func writer(filename string) io.Writer {
	hook, err := rotateLogs.New(
		filename+".%Y-%m-%d-%H",
		rotateLogs.WithLinkName(filename),
		rotateLogs.WithMaxAge(time.Hour*24*7),
		rotateLogs.WithRotationTime(time.Hour),
	)
	if err != nil {
		panic(fmt.Errorf("rotate logs init failed: %v", err))
	}
	return hook
}

func customEncodeTime(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(fmt.Sprintf("[%s] %s", consts.ApplicationName, t.Format("2006-01-02 15:04:05")))
}
