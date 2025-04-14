package types

import (
	"github.com/MisakaTAT/GTerm/backend/enums"
	"io"
)

type FileTransferItemInfo struct {
	Name        string `json:"name"`
	Size        int64  `json:"size"`
	IsDir       bool   `json:"isDir"`
	ModTime     string `json:"modTime"`
	Permissions string `json:"permissions"`
}

type FileList struct {
	Files        []FileTransferItemInfo `json:"files"`
	AbsolutePath string                 `json:"absolutePath"`
}

type FileTransferTask struct {
	ID          string                      `json:"id"`
	Source      string                      `json:"source"`
	Destination string                      `json:"destination"`
	Size        int64                       `json:"size"`
	Transferred int64                       `json:"transferred"`
	IsUpload    bool                        `json:"isUpload"`
	Status      enums.FileTransferTaskState `json:"status"`
	Error       *string                     `json:"error"`
}

type ProgressReader struct {
	Reader           io.Reader
	TotalSize        int64
	BytesRead        int64
	ProgressCallback func(int64, int64)
}

func (pr *ProgressReader) Read(p []byte) (int, error) {
	n, err := pr.Reader.Read(p)
	if n > 0 {
		pr.BytesRead += int64(n)
		if pr.ProgressCallback != nil {
			pr.ProgressCallback(pr.BytesRead, pr.TotalSize)
		}
	}
	return n, err
}
