package services

import (
	"fmt"
	"runtime"
	"time"

	"github.com/Q191/GTerm/backend/initialize"
	"github.com/Q191/GTerm/backend/pkg/base"
	"github.com/google/wire"
)

var PreferencesSrvSet = wire.NewSet(wire.Struct(new(PreferencesSrv), "*"))

type PreferencesSrv struct {
	Logger initialize.Logger
}

func (s *PreferencesSrv) Version() string {
	return base.Version
}

func (s *PreferencesSrv) VersionURL() string {
	return base.VersionURL
}

func (s *PreferencesSrv) Copyright() string {
	return fmt.Sprintf("Copyright © 2025-%d Q191.\nAll rights reserved.", time.Now().Year())
}

func (s *PreferencesSrv) GOOS() string {
	return runtime.GOOS
}

func (s *PreferencesSrv) IsDarwin() bool {
	return s.GOOS() == "darwin"
}
