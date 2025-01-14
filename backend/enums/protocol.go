package enums

type ConnProtocol = int

const (
	SSH ConnProtocol = iota
	Telnet
	RDP
	VNC
	Serial
)
