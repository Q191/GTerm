package metadata

import (
	"strings"

	"github.com/MisakaTAT/GTerm/backend/pkg/exec"
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
	"show version",        // 网络设备通用
	"display version",     // 华为等
	"get system info",     // H3C等
	"uname -a",            // Unix/Linux
	"ver",                 // Windows
	"cat /etc/os-release", // Linux发行版
}

var networkVendorKeywords = map[string][]string{
	"cisco": {
		"cisco ios", "cisco ios-xe", "cisco ios-xr", "cisco nx-os",
		"cisco asa", "cisco", "catalyst", "nexus", "aironet",
		"cisco firepower", "cisco meraki", "cisco packet tracer",
	},
	"huawei": {
		"huawei vrp", "huawei", "huawei ce", "huawei ne",
		"huawei s", "huawei ar", "huawei usg", "huawei ensp",
		"huawei cloudengine", "huawei secospace", "huawei oceanstor",
		"huawei fusionsphere", "huawei fusionaccess", "huawei fusioncloud",
	},
	"fortinet": {
		"fortinet", "fortigate", "fortios", "fortimanager",
		"fortianalyzer", "forticlient", "fortiweb", "fortimail",
		"fortisandbox", "fortiadc", "fortiddos", "fortiap",
	},
	"mikrotik": {
		"mikrotik", "routeros", "routerboard", "winbox",
		"mikrotik chr", "mikrotik switch", "mikrotik cloud",
	},
	"pfsense": {
		"pfsense", "opnsense", "m0n0wall", "pfctl",
		"pfsense ce", "pfsense plus", "netgate",
	},
	"juniper": {
		"juniper", "junos", "juniper ex", "juniper mx", "juniper srx",
		"juniper qfx", "juniper t", "juniper ptx", "juniper j",
		"juniper space", "juniper contrail", "junos pulse",
	},
	"hp": {
		"hp", "hpe", "hp procurve", "hp comware", "hpe aruba",
		"hp virtual connect", "hp proliant", "hpe simplivity",
		"hp bladesystem", "hp c-class", "hpe synergy",
	},
	"dell": {
		"dell", "dell emc", "dell networking", "dell poweredge",
		"dell force10", "dell powerconnect", "dell powervault",
		"dell powerswitch", "dell os10", "dell os9", "dell os6",
	},
	"arista": {
		"arista", "arista eos", "arista cloudvision",
		"arista data-center", "arista 7000", "arista 7500",
		"arista 7300", "arista 7280", "arista 7200", "arista 7150",
		"arista 7130", "arista 7050", "arista 7020", "arista 7010",
	},
	"h3c": {
		"h3c", "h3c os", "h3c switch", "h3c router", "h3c msr",
		"h3c wlan", "h3c secpath", "h3c s", "h3c e", "h3c r",
		"h3c comware", "h3c magic",
	},
	"checkpoint": {
		"checkpoint", "check point", "gaia", "checkpoint r80",
		"checkpoint r81", "checkpoint ngtp", "checkpoint ngfw",
		"checkpoint sandblast", "checkpoint quantum", "checkpoint harmony",
	},
	"paloalto": {
		"palo alto", "paloalto", "pan-os", "panos", "prisma",
		"panorama", "palo alto networks", "wildfire", "globalprotect",
		"palo alto vm-series", "palo alto pa-series",
	},
	"f5": {
		"f5", "f5 big-ip", "f5 ltm", "f5 gtm", "f5 asm",
		"f5 apm", "f5 afm", "f5 access", "f5 networks",
		"f5 silverline", "f5 viprion", "f5 i-series",
	},
	"citrix": {
		"citrix adc", "citrix netscaler", "citrix sdwan",
		"citrix gateway", "citrix application delivery",
	},
	"extreme": {
		"extreme networks", "extreme", "exos", "xos",
		"extreme summit", "extreme vsp", "extreme fabric engine",
		"extreme switching", "extreme wireless",
	},
	"vmware": {
		"vmware nsx", "vmware esxi networking", "vmware vcenter networking",
		"vmware vds", "vmware vsphere networking", "vmware vsan networking",
	},
	"sonicwall": {
		"sonicwall", "sonic wall", "sonicwall tz", "sonicwall nsa",
		"sonicwall supermassive", "sonicwall sonicos", "sonicwall capture",
	},
	"ubiquiti": {
		"ubiquiti", "unifi", "edgerouter", "edgeswitch", "edgemax",
		"airmax", "ubiquiti networks", "ubnt", "ubiquiti amplifi",
	},
	"zyxel": {
		"zyxel", "zy-xel", "zywall", "zynos", "zyxel nebula",
		"zyxel zywall", "zyxel gs", "zyxel xgs", "zyxel usg",
	},
	"netgear": {
		"netgear", "netgear prosafe", "netgear gs", "netgear xs",
		"netgear m", "netgear readynas", "netgear insight",
	},
	"barracuda": {
		"barracuda", "barracuda networks", "barracuda firewall",
		"barracuda waf", "barracuda cloudgen", "barracuda cgf",
	},
	"tplink": {
		"tp-link", "tplink", "tp link", "tp-link jetstream", "tp-link omada",
		"tp-link smart switch", "tp-link business", "tp-link safestream",
	},
	"dlink": {
		"d-link", "dlink", "d link", "d-link nuclias", "d-link des",
		"d-link dgs", "d-link dxs", "d-link dws",
	},
	"versa": {
		"versa networks", "versa", "versa director", "versa analytics",
		"versa secure sd-wan", "versa titan", "versa flexvnf",
	},
	"sophos": {
		"sophos firewall", "sophos xg", "sophos sg", "sophos utm",
		"sophos cyberoam", "sophos central", "sophos red",
	},
	"watchguard": {
		"watchguard", "watch guard", "watchguard firebox", "watchguard xtm",
		"watchguard dimension", "watchguard system manager",
	},
}

var serverVendorKeywords = map[string][]string{
	// Linux发行版
	"redhat": {
		"red hat enterprise", "red hat", "rhel", "redhat",
		"red hat enterprise linux", "rhel 9", "rhel 8", "rhel 7",
		"red hat linux", "redhat release", "redhat version",
	},
	"ubuntu": {
		"ubuntu", "ubuntu server", "xubuntu", "kubuntu",
		"ubuntu lts", "ubuntu core", "ubuntu mate", "ubuntu desktop",
		"ubuntu focal", "ubuntu jammy", "ubuntu bionic", "ubuntu noble",
		"ubuntu 24.04", "ubuntu 22.04", "ubuntu 20.04", "ubuntu 18.04",
	},
	"centos": {
		"centos stream", "centos", "centos linux",
		"centos 7", "centos 8", "centos 9", "centos release",
	},
	"debian": {
		"debian gnu", "debian", "debian linux", "debian testing",
		"debian stable", "debian sid", "debian bookworm", "debian bullseye",
		"debian buster", "debian stretch", "debian jessie",
		"debian 12", "debian 11", "debian 10", "debian 9",
	},
	"opensuse": {
		"opensuse leap", "opensuse tumbleweed", "opensuse", "suse",
		"suse linux", "suse enterprise", "sles", "opensuse microos",
		"suse manager", "sle", "sles 15", "sles 12", "leap 15",
	},
	"fedora": {
		"fedora server", "fedora workstation", "fedora", "fedora linux",
		"fedora core", "fedora rawhide", "fedora silverblue", "fedora iot",
		"fedora 39", "fedora 40", "fedora 38", "fedora 37",
	},
	"almalinux": {
		"alma linux", "almalinux", "alma", "almalinux 8", "almalinux 9",
		"almalinux foundation", "almalinux os",
	},
	"kalilinux": {
		"kali linux", "kalilinux", "kali", "kali rolling", "kali purple",
		"kali security", "kali penetration", "kali testing",
	},
	"archlinux": {
		"arch linux", "archlinux", "arch", "manjaro", "endeavouros",
		"garuda linux", "arco linux", "blackarch", "artix linux",
		"archcraft", "arch based", "rolling release", "pacman",
	},
	"rockylinux": {
		"rocky linux", "rockylinux", "rocky", "rocky 8", "rocky 9",
		"rocky enterprise", "rocky foundation",
	},
	"alpinelinux": {
		"alpine linux", "alpinelinux", "alpine", "alpine musl",
		"alpine edge", "alpine standard", "alpine minimal",
		"alpine 3.19", "alpine 3.18", "alpine 3.17",
	},
	"gentoo": {
		"gentoo linux", "gentoo", "gentoo portage", "gentoo prefix",
		"gentoo hardened", "gentoo stage3", "gentoo emerge",
	},
	"raspberrypi": {
		"raspberry pi os", "raspbian", "raspberry pi", "rpi os",
		"raspberry pi legacy", "raspberry pi desktop", "raspberry pi lite",
	},
	"linuxmint": {
		"mint", "linux mint", "mint cinnamon", "mint mate", "mint xfce",
		"mint 21", "mint 20", "mint vera", "mint victoria", "mint vanessa",
	},
	"elementary": {
		"elementary os", "elementary", "elementary horus", "elementary juno",
		"elementary loki", "elementary odin", "elementary 7", "elementary 6",
	},
	"zorinos": {
		"zorin os", "zorin", "zorin core", "zorin lite", "zorin pro",
		"zorin education", "zorin 17", "zorin 16", "zorin 15",
	},
	"popos": {
		"pop!_os", "pop os", "popos", "pop_os", "pop!_os 22.04", "pop!_os 20.04",
		"pop!_os system76", "pop!_os cosmic",
	},
	"oracle": {
		"oracle linux", "ol", "oracle enterprise linux", "oracle cloud",
		"oracle 7", "oracle 8", "oracle 9", "oracle autonomous",
	},
	"void": {
		"void linux", "void", "void musl", "void glibc",
		"void lxqt", "void xfce", "void runit",
	},
	"nixos": {
		"nixos", "nix", "nixpkgs", "nix-shell", "nix-env",
		"nix flakes", "nixos configuration", "nixos generation",
	},
	"clearlinux": {
		"clear linux", "clearlinux", "clear linux os", "clear linux project",
		"intel clear linux", "clear linux distro",
	},
	"slackware": {
		"slackware", "slackware linux", "slackware-current",
		"slackware 15", "slackware 14", "slackware lts",
	},
	"mageia": {
		"mageia", "mageia linux", "mageia 8", "mageia 9",
		"mageia cauldron", "mageia classic",
	},
	"mx": {
		"mx linux", "mxlinux", "mx", "mx fluxbox", "mx kde",
		"mx xfce", "antiX mx", "mx 23", "mx 21",
	},

	// BSD系统
	"freebsd": {
		"freebsd", "bsd", "freebsd 14", "freebsd 13", "freebsd 12",
		"freebsd release", "freebsd ports", "freebsd jail", "freebsd pkg",
	},
	"openbsd": {
		"openbsd", "open bsd", "openbsd 7", "openbsd current",
		"openbsd stable", "openbsd release", "openbsd ports",
	},
	"netbsd": {
		"netbsd", "net bsd", "netbsd 9", "netbsd 10",
		"netbsd pkgsrc", "netbsd release", "netbsd current",
	},
	"trueos": {
		"trueos", "true os", "truenas core", "truenas scale",
		"truenas enterprise", "ix systems", "freenas",
	},
	"dragonflybsd": {
		"dragonfly bsd", "dragonflybsd", "dragonfly", "dfly",
		"dragonfly 6", "dragonfly hammer", "dragonfly hammer2",
	},

	// 其他操作系统
	"windows": {
		"windows", "windows server", "windows desktop", "windows 10", "windows 11",
		"windows 2016", "windows 2019", "windows 2022", "windows azure",
		"microsoft windows", "win32", "powershell", "cmd.exe", "winnt",
		"windows nt", "windows build", "windows version",
	},
	"macos": {
		"macos", "mac os", "darwin", "apple", "macos ventura", "macos sonoma",
		"macos monterey", "macos big sur", "macos catalina", "macos mojave",
		"mac os x", "osx", "macbook", "imac", "mac mini", "mac studio",
		"apple silicon", "intel mac", "hackintosh",
	},
	"solaris": {
		"solaris", "oracle solaris", "sun solaris", "illumos", "openindiana",
		"omnios", "smartos", "solaris 11", "solaris 10", "solaris zones",
	},
	"aix": {
		"aix", "ibm aix", "aix unix", "aix power", "aix 7", "aix 6",
		"ibm power systems", "powervm", "aix lpar", "aix wpar",
	},
	"hpux": {
		"hp-ux", "hpux", "hp unix", "hp-ux itanium", "hp-ux 11i",
		"hp-ux 11.31", "hp-ux 11.23", "hp integrity",
	},
	"esxi": {
		"esxi", "vmware esxi", "vsphere", "esxi 8", "esxi 7", "esxi 6.7",
		"esxi hypervisor", "vmware vsan", "vmware vcenter", "esxcli",
	},
	"proxmox": {
		"proxmox", "proxmox ve", "proxmox virtual environment", "pve",
		"proxmox backup server", "proxmox mail gateway", "proxmox 8", "proxmox 7",
	},
	"android": {
		"android", "android os", "aosp", "android open source project",
		"android 14", "android 13", "android 12", "android 11",
	},
	"ios": {
		"ios", "apple ios", "ipados", "ios 17", "ios 16", "ios 15",
		"iphone os", "ipados 17", "ipados 16", "ipados 15",
	},
	"haiku": {
		"haiku", "haiku os", "beos", "haiku release", "haiku beta",
	},
	"reactos": {
		"reactos", "react os", "open source windows", "reactos 0.4",
	},
	"chromeos": {
		"chrome os", "chromeos", "chromium os", "google chromeos",
		"chrome os flex", "cloudready", "chrome os 115", "chrome os 114",
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
