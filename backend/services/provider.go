package services

import "github.com/google/wire"

var SetProvider = wire.NewSet(
	TerminalSrvSet,
	PreferencesSrvSet,
)
