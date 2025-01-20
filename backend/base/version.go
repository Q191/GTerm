package base

import (
	"fmt"
	"runtime/debug"
	"strings"
)

const (
	commitPrefix  = "https://github.com/MisakaTAT/GTerm/commit/"
	releasePrefix = "https://github.com/MisakaTAT/Shiro/releases/tag/"
	versionPrefix = "v"
)

var (
	// Version 编译时使用 ldflags 覆盖版本信息
	Version = "Unknown"
)

// formatVersion 格式化版本号
func formatVersion(version string) string {
	if version == "Unknown" {
		return version
	}
	// 如果版本号不是以 v 开头，则添加
	if !strings.HasPrefix(version, versionPrefix) {
		return versionPrefix + version
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
			Version = setting.Value[:7] // 只使用前7位commit hash
		}
	}
	Version = formatVersion(Version)
}

// GetVersion 获取版本号
func GetVersion() string {
	return Version
}

// GetCommitURL 获取commit URL
func GetCommitURL() string {
	if Version == "Unknown" || strings.HasPrefix(Version, versionPrefix) {
		return ""
	}
	return fmt.Sprintf("%s%s", commitPrefix, Version)
}

// GetReleaseURL 获取release URL
func GetReleaseURL() string {
	if !strings.HasPrefix(Version, versionPrefix) {
		return ""
	}
	return fmt.Sprintf("%s%s", releasePrefix, Version)
} 