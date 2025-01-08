package metadata

import (
	"github.com/MisakaTAT/GTerm/backend/pkg/exec"
	"golang.org/x/crypto/ssh"
	"strings"
)

type Metadata struct {
	exec *exec.Adapter
}

type OSRelease struct {
	ID         string
	Name       string
	Version    string
	VersionID  string
	PrettyName string
}

func NewMetadata(client *ssh.Client) *Metadata {
	return &Metadata{exec: exec.New(client)}
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
