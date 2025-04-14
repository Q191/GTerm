package enums

import "strings"

type FileTransferTaskState string

const (
	Pending   FileTransferTaskState = "Pending"
	Progress  FileTransferTaskState = "Progress"
	Completed FileTransferTaskState = "Completed"
	Error     FileTransferTaskState = "Error"
)

var FileTransferTaskStateEnums = []FileTransferTaskState{Pending, Progress, Completed, Error}

func (a FileTransferTaskState) TSName() string {
	return strings.ToUpper(string(a))
}
