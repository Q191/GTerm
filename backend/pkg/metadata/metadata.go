package metadata

import (
	"strings"

	"github.com/Q191/GTerm/backend/pkg/exec"
	"golang.org/x/crypto/ssh"
)

type Metadata struct {
	exec *exec.Adapter
}

type SessionMetadata struct {
	Type      string            // 设备类型
	Vendor    string            // 厂商
	RawOutput map[string]string // 原始命令输出
}

var universalCommands = []string{
	"show version",                   // 网络设备通用
	"display version",                // 华为等
	"get system info",                // H3C等
	"uname -a",                       // Unix/Linux
	"ver",                            // Windows
	"cat /etc/os-release",            // Linux发行版
	"service-control --status --all", // vCenter
}

var networkVendorKeywords = map[string][]string{
	"cisco": {
		"cisco ios", "cisco ios-xe", "cisco ios-xr", "cisco nx-os",
		"cisco asa", "cisco", "cisco catalyst", "cisco nexus", "cisco meraki",
	},
	"huawei": {
		"huawei vrp", "huawei", "huawei ce", "huawei ne",
		"huawei s", "huawei ar", "huawei usg", "huawei ensp",
	},
	"fortinet": {
		"fortinet", "fortigate", "fortios", "fortimanager",
		"fortianalyzer", "forticlient", "fortiweb",
	},
	"mikrotik": {
		"mikrotik", "routeros", "routerboard", "mikrotik winbox",
	},
	"pfsense": {
		"pfsense", "opnsense", "m0n0wall", "netgate",
	},
	"juniper": {
		"juniper", "juniper junos", "juniper ex", "juniper mx", "juniper srx",
		"juniper qfx", "juniper t", "juniper ptx",
	},
	"hp": {
		"hp", "hpe", "hp procurve", "hp comware", "hpe aruba",
	},
	"dell": {
		"dell", "dell emc", "dell networking", "dell poweredge",
		"dell force10", "dell powerconnect",
	},
	"arista": {
		"arista", "arista eos", "arista cloudvision",
	},
	"h3c": {
		"h3c", "h3c os", "h3c switch", "h3c router", "h3c msr",
	},
	"checkpoint": {
		"checkpoint", "check point", "gaia", "checkpoint r80",
	},
	"paloalto": {
		"palo alto", "paloalto", "palo alto pan-os", "palo alto panorama",
	},
	"f5": {
		"f5", "f5 big-ip", "f5 ltm", "f5 gtm",
	},
	"citrix": {
		"citrix adc", "citrix netscaler", "citrix sdwan",
	},
	"extreme": {
		"extreme networks", "extreme", "exos",
	},
	"vmware": {
		"vmware", "esxi", "vcenter", "vsphere", "vmware nsx",
	},
	"sonicwall": {
		"sonicwall", "sonic wall", "sonicwall tz",
	},
	"ubiquiti": {
		"ubiquiti", "ubiquiti unifi", "ubiquiti edgerouter", "ubiquiti edgeswitch",
	},
	"zyxel": {
		"zyxel", "zy-xel", "zyxel zywall", "zyxel zynos",
	},
	"netgear": {
		"netgear", "netgear prosafe", "netgear gs",
	},
	"barracuda": {
		"barracuda", "barracuda networks", "barracuda firewall",
	},
	"tplink": {
		"tp-link", "tplink", "tp link", "tp-link jetstream",
	},
	"dlink": {
		"d-link", "dlink", "d link", "d-link nuclias",
	},
	"versa": {
		"versa networks", "versa", "versa director",
	},
	"sophos": {
		"sophos firewall", "sophos xg", "sophos sg",
	},
	"watchguard": {
		"watchguard", "watch guard", "watchguard firebox",
	},
}

var serverVendorKeywords = map[string][]string{
	// Linux发行版
	"redhat": {
		"red hat", "red hat rhel", "redhat",
	},
	"ubuntu": {
		"ubuntu", "ubuntu server", "ubuntu lts",
	},
	"centos": {
		"centos", "centos stream",
	},
	"debian": {
		"debian",
	},
	"opensuse": {
		"opensuse", "opensuse suse", "opensuse sles",
	},
	"fedora": {
		"fedora", "fedora server", "fedora workstation",
	},
	"almalinux": {
		"almalinux", "alma",
	},
	"kalilinux": {
		"kali", "kalilinux",
	},
	"archlinux": {
		"arch", "archlinux", "manjaro",
	},
	"rockylinux": {
		"rocky", "rockylinux",
	},
	"alpinelinux": {
		"alpine", "alpinelinux",
	},
	"gentoo": {
		"gentoo", "gentoo linux",
	},
	"raspberrypi": {
		"raspberry pi", "raspberry pi raspbian",
	},
	"linuxmint": {
		"mint", "linux mint",
	},
	"elementary": {
		"elementary", "elementary os",
	},
	"zorinos": {
		"zorin", "zorin os",
	},
	"popos": {
		"pop!_os", "pop os",
	},
	"oracle": {
		"oracle linux",
	},

	// BSD系统
	"freebsd": {
		"freebsd", "freebsd bsd", "freebsd 14",
		"freebsd 13", "freebsd 12",
	},
	"openbsd": {
		"openbsd", "openbsd open bsd", "openbsd 7",
		"openbsd current", "openbsd stable",
	},

	// 其他操作系统
	"windows": {
		"windows", "windows server", "windows desktop",
		"windows 10", "windows 11", "windows 2016",
		"windows 2019", "windows 2022",
	},
	"macos": {
		"macos", "mac os", "macos darwin", "apple",
		"macos ventura", "macos sonoma", "macos monterey",
		"macos big sur", "macos catalina",
	},
	"proxmox": {
		"proxmox", "proxmox ve", "proxmox virtual environment",
		"pve", "proxmox backup server", "proxmox 8",
	},
	"chromeos": {
		"chrome os", "chromeos", "chromium os",
		"google chromeos", "chrome os flex",
	},
}

func NewMetadata(client *ssh.Client) *Metadata {
	return &Metadata{exec: exec.New(client)}
}

func (m *Metadata) Parser() *SessionMetadata {
	info := &SessionMetadata{
		RawOutput: make(map[string]string),
	}

	for _, cmd := range universalCommands {
		output := m.exec.Run(cmd).Unwrap()
		if output == "" {
			continue
		}
		info.RawOutput[cmd] = output
		vendor, deviceType := m.identifyVendorAndType(output)
		info.Vendor = vendor
		info.Type = deviceType
		if info.Vendor != "" {
			return info
		}
	}

	return info
}

func (m *Metadata) identifyVendorAndType(output string) (string, string) {
	if output == "" {
		return "", "unknown"
	}

	lowOutput := strings.ToLower(output)

	var bestMatch struct {
		vendor     string
		deviceType string
		length     int
	}
	for i := 0; i < len(lowOutput); i++ {
		current := keywordTrie.root
		matchLength := 0
		for j := i; j < len(lowOutput); j++ {
			char := rune(lowOutput[j])
			if next, exists := current.children[char]; exists {
				current = next
				matchLength++
				if current.isEndOfWord && matchLength > bestMatch.length {
					bestMatch.vendor = current.vendor
					bestMatch.deviceType = current.deviceType
					bestMatch.length = matchLength
				}
			} else {
				break
			}
		}
	}

	if bestMatch.length > 0 {
		return bestMatch.vendor, bestMatch.deviceType
	}

	return "", "unknown"
}
