package types

import (
	"github.com/MisakaTAT/GTerm/backend/enums"
)

type FileTransferItemInfo struct {
	Name        string `json:"name"`
	Size        int64  `json:"size"`
	IsDir       bool   `json:"isDir"`
	ModTime     string `json:"modTime"`
	Permissions string `json:"permissions"`
	Owner       string `json:"owner"`
	Group       string `json:"group"`
	Extended    string `json:"extended"`
}

type FileList struct {
	Files        []*FileTransferItemInfo `json:"files"`
	AbsolutePath string                  `json:"absolutePath"`
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
