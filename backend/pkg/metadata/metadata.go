package metadata

import (
	"github.com/MisakaTAT/GTerm/backend/pkg/exec"
	"golang.org/x/crypto/ssh"
	"strconv"
	"strings"
)

type Metadata struct {
	KernelVersion string
	Arch          string
	CPU           *CPUInfo
	Memory        *MemoryInfo
	OSRelease     *OSRelease
	exec          *exec.Adapter
}

type OSRelease struct {
	ID         string
	Name       string
	Version    string
	VersionID  string
	PrettyName string
}

type CPUInfo struct {
	Cores uint   `json:"cores"` // CPU核心数
	Model string `json:"model"` // CPU型号
}

type MemoryInfo struct {
	Total uint `json:"total"` // 总内存(GB)
}

func NewMetadata(client *ssh.Client) *Metadata {
	return &Metadata{exec: exec.New(client)}
}

func (m *Metadata) Fetch() *Metadata {
	m.CPU = m.GetCPUInfo()
	m.Arch = m.GetArch()
	m.OSRelease = m.GetOSRelease()
	m.KernelVersion = m.GetKernelVersion()
	m.Memory = m.GetMemoryInfo()
	return m
}

func (m *Metadata) GetMemoryInfo() *MemoryInfo {
	var info MemoryInfo

	memInfo := m.exec.Run("cat /proc/meminfo").Unwrap()
	lines := strings.Split(memInfo, "\n")
	var memValues = make(map[string]uint64)

	for _, line := range lines {
		fields := strings.Fields(line)
		if len(fields) < 2 {
			continue
		}
		key := strings.TrimSuffix(fields[0], ":")
		if val, err := strconv.ParseUint(fields[1], 10, 64); err == nil {
			memValues[key] = val
		}
	}

	toGB := func(kb uint64) uint {
		return uint(float64(kb)/(1000.0*1000.0) + 0.5)
	}

	info.Total = toGB(memValues["MemTotal"])

	return &info
}

func (m *Metadata) GetArch() string {
	return m.exec.Run("arch").Unwrap()
}

func (m *Metadata) GetKernelVersion() string {
	return m.exec.Run("uname -r").Unwrap()
}

func (m *Metadata) GetHostname() string {
	return m.exec.Run("hostname").Unwrap()
}

func (m *Metadata) GetCPUInfo() *CPUInfo {
	var info CPUInfo

	model := m.exec.Run("grep 'model name' /proc/cpuinfo | head -n1 | cut -d ':' -f2").Unwrap()
	info.Model = strings.TrimSpace(model)

	cores := m.exec.Run("grep -c '^processor' /proc/cpuinfo").Unwrap()
	if coresInt, err := strconv.ParseUint(strings.TrimSpace(cores), 10, 64); err == nil {
		info.Cores = uint(coresInt)
	}
	return &info
}

func (m *Metadata) GetOSRelease() *OSRelease {
	var r OSRelease
	output := m.exec.Run("cat /etc/os-release").Unwrap()
	if output == "" {
		output = m.exec.Run("cat /etc/issue").Unwrap()
		if output != "" {
			r.Name = strings.TrimSpace(output)
			r.PrettyName = r.Name
		}
		return &r
	}
	lines := strings.Split(output, "\n")
	for _, line := range lines {
		if line = strings.TrimSpace(line); line == "" {
			continue
		}
		kv := strings.SplitN(line, "=", 2)
		if len(kv) != 2 {
			continue
		}
		field := strings.Trim(strings.TrimSpace(kv[0]), "\"")
		value := strings.Trim(strings.TrimSpace(kv[1]), "\"")
		switch field {
		case "NAME":
			r.Name = value
		case "VERSION":
			r.Version = value
		case "ID":
			r.ID = strings.ToLower(value)
		case "VERSION_ID":
			r.VersionID = value
		case "PRETTY_NAME":
			r.PrettyName = value
		}
	}
	return &r
}
