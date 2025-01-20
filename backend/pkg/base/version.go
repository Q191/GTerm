package base

import (
	"runtime/debug"
	"strings"
)

const (
	commitPrefix  = "https://github.com/MisakaTAT/GTerm/commit/"
	releasePrefix = "https://github.com/MisakaTAT/Shiro/releases/tag/"
)

var (
	// Version 编译时使用 ldflags 覆盖版本信息
	Version        = "Unknown"
	VersionURL     = ""
	fullCommitHash = ""
)

func formatVersion(version string) string {
	if version == "Unknown" {
		return version
	}
	if !strings.HasPrefix(version, "v") {
		fullCommitHash = version
		return version[:7]
	}
	return version
}

func init() {
	if Version != "Unknown" {
		Version = formatVersion(Version)
		return
	}

	info, ok := debug.ReadBuildInfo()
	if !ok {
		return
	}

	Version = info.Main.Version
	for _, setting := range info.Settings {
		if setting.Key == "vcs.revision" {
			Version = setting.Value
		}
	}
	Version = formatVersion(Version)

	if !strings.HasPrefix(Version, "v") {
		VersionURL = commitPrefix + fullCommitHash
		return
	}
	VersionURL = releasePrefix + Version
}
