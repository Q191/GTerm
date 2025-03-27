package services

import "github.com/google/wire"

var SetProvider = wire.NewSet(
	TerminalSrvSet,
	PreferencesSrvSet,
	GroupSrvSet,
	ConnectionSrvSet,
	MetadataSrvSet,
	CredentialSrvSet,
	WebsocketSrvSet,
)
