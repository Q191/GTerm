package types

import (
	"fmt"

	"github.com/MisakaTAT/GTerm/backend/enums"
)

type Message struct {
	Type        enums.TerminalType `json:"type"`
	Content     any                `json:"content,omitempty"`
	Message     string             `json:"message,omitempty"`
	Details     string             `json:"details,omitempty"`
	Code        string             `json:"code,omitempty"`
	Host        string             `json:"host,omitempty"`
	Fingerprint string             `json:"fingerprint,omitempty"`
}

type Fingerprint struct {
	Type   enums.TerminalType `json:"type"`
	Accept bool               `json:"accept"`
}

type FingerprintError struct {
	Host        string `json:"host"`
	Fingerprint string `json:"fingerprint"`
}

func (e *FingerprintError) Error() string {
	return fmt.Sprintf("unknown host fingerprint: %s", e.Fingerprint)
}
