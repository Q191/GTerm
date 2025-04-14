package sftp

import (
	"os"
	"path/filepath"
	"strings"
)

func GetFileName(path string) string {
	return filepath.Base(path)
}

func JoinPath(base, file string) string {
	if strings.HasPrefix(base, "/") {
		base = strings.TrimSuffix(base, "/")
		return base + "/" + file
	}
	return filepath.Join(base, file)
}

func GetFileSize(path string) (int64, error) {
	fileInfo, err := os.Stat(path)
	if err != nil {
		return 0, err
	}
	return fileInfo.Size(), nil
}
