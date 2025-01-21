package enums

import "strings"

type ConnProtocol string

const (
	SSH    ConnProtocol = "SSH"
	Telnet ConnProtocol = "Telnet"
	RDP    ConnProtocol = "RDP"
	VNC    ConnProtocol = "VNC"
	Serial ConnProtocol = "Serial"
)

var ConnProtocolEnums = []ConnProtocol{SSH, Telnet, RDP, VNC, Serial}

func (c ConnProtocol) TSName() string {
	return strings.ToUpper(string(c))
}
