package services

import (
	"github.com/OpenToolkitLab/GTerm/backend/pkg/base"
	"github.com/google/wire"
	"go.uber.org/zap"
	"runtime"
)

var PreferencesSrvSet = wire.NewSet(wire.Struct(new(PreferencesSrv), "*"))

type PreferencesSrv struct {
	Logger *zap.Logger
}

func (s *PreferencesSrv) GTermVer() string {
	return base.Version
}

func (s *PreferencesSrv) GOOS() string {
	return runtime.GOOS
}

func (s *PreferencesSrv) IsDarwin() bool {
	return s.GOOS() == "darwin"
}
