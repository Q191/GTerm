package services

import (
	"sync"

	"github.com/MisakaTAT/GTerm/backend/consts/messages"
	"github.com/MisakaTAT/GTerm/backend/initialize"
	"github.com/MisakaTAT/GTerm/backend/pkg/sftp"
	commonssh "github.com/MisakaTAT/GTerm/backend/pkg/ssh"
	"github.com/MisakaTAT/GTerm/backend/types"
	"github.com/MisakaTAT/GTerm/backend/utils/resp"
	"github.com/google/wire"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

var FileTransferSrvSet = wire.NewSet(wire.Struct(new(FileTransferSrv), "*"))

type FileTransferSrv struct {
	Logger           initialize.Logger
	ConnectionSrv    *ConnectionSrv
	AppContext       *initialize.AppContext
	SFTPHandler      *sftp.Handler `wire:"-"`
	SFTPHandlerMutex sync.Mutex    `wire:"-"`
}

func (s *FileTransferSrv) ConnectSFTP(connID uint) *resp.Resp {
	s.Logger.Info("Connecting to SFTP server, connection ID: %d", connID)

	s.SFTPHandlerMutex.Lock()
	defer s.SFTPHandlerMutex.Unlock()

	if s.SFTPHandler != nil && s.SFTPHandler.IsConnected {
		s.Logger.Info("Closing existing SFTP connection")
		if err := s.SFTPHandler.Close(); err != nil {
			s.Logger.Error("Error closing existing SFTP connection: %v", err)
		}
	}

	s.SFTPHandler = sftp.NewSFTPHandler(s.Logger)

	conn, err := s.ConnectionSrv.FindByID(connID)
	if err != nil {
		s.Logger.Error("Connection not found: %v, connID: %d", err, connID)
		return resp.FailWithMsg(err.Error())
	}

	conf := &commonssh.Config{
		Host:             conn.Host,
		Port:             conn.Port,
		User:             conn.Credential.Username,
		AuthMethod:       conn.Credential.AuthMethod,
		Password:         conn.Credential.Password,
		PrivateKey:       conn.Credential.PrivateKey,
		Passphrase:       conn.Credential.Passphrase,
		TrustUnknownHost: true,
	}

	if err = s.SFTPHandler.Connect(conf); err != nil {
		s.Logger.Error("Failed to connect to SFTP server: %v", err)
		return resp.FailWithMsg(err.Error())
	}

	s.Logger.Info("SFTP connection successful")
	return resp.OkWithCode(messages.Connected)
}

func (s *FileTransferSrv) DisconnectSFTP() *resp.Resp {
	s.Logger.Info("Disconnecting from SFTP server")

	s.SFTPHandlerMutex.Lock()
	defer s.SFTPHandlerMutex.Unlock()

	if s.SFTPHandler == nil || !s.SFTPHandler.IsConnected {
		s.Logger.Info("No active SFTP connection")
		return resp.Ok()
	}

	if err := s.SFTPHandler.Close(); err != nil {
		s.Logger.Error("Error disconnecting from SFTP server: %v", err)
		return resp.FailWithMsg(err.Error())
	}

	s.Logger.Info("SFTP connection closed")
	return resp.OkWithCode(messages.Disconnected)
}

func (s *FileTransferSrv) ListRemoteFiles(path string) *resp.Resp {
	s.Logger.Info("Listing remote files, path: %s", path)

	s.SFTPHandlerMutex.Lock()
	defer s.SFTPHandlerMutex.Unlock()

	if s.SFTPHandler == nil || !s.SFTPHandler.IsConnected {
		s.Logger.Error("Not connected to SFTP server")
		return resp.FailWithMsg("Not connected to SFTP server")
	}

	processedPath, err := s.SFTPHandler.ProcessPath(path)
	if err != nil {
		s.Logger.Error("Failed to process path %s: %v", path, err)
		return resp.FailWithMsg(err.Error())
	}
	s.Logger.Info("Using processed path for listing: %s", processedPath)

	files, err := s.SFTPHandler.ListRemoteFiles(processedPath)
	if err != nil {
		s.Logger.Error("Failed to list remote files: %v", err)
		return resp.FailWithMsg(err.Error())
	}

	response := &types.FileList{
		Files:        files,
		AbsolutePath: processedPath,
	}

	s.Logger.Info("Successfully listed remote files at path %s, found %d files/directories", processedPath, len(files))
	return resp.OkWithData(response)
}

func (s *FileTransferSrv) UploadFiles(localPaths []string, remotePath string) *resp.Resp {
	s.Logger.Info("Uploading files, remote path: %s", remotePath)

	s.SFTPHandlerMutex.Lock()
	defer s.SFTPHandlerMutex.Unlock()

	if s.SFTPHandler == nil || !s.SFTPHandler.IsConnected {
		s.Logger.Error("Not connected to SFTP server")
		return resp.FailWithMsg("Not connected to SFTP server")
	}

	processedPath, err := s.SFTPHandler.ProcessPath(remotePath)
	if err != nil {
		s.Logger.Error("Failed to process path %s: %v", remotePath, err)
		return resp.FailWithMsg(err.Error())
	}
	remotePath = processedPath
	s.Logger.Info("Using processed remote path: %s", remotePath)

	totalFiles := len(localPaths)
	completedFiles := 0

	for _, localPath := range localPaths {
		s.Logger.Info("Uploading file: %s -> %s", localPath, remotePath)

		fileName := sftp.GetFileName(localPath)
		remoteFilePath, err := s.SFTPHandler.JoinRemotePaths(remotePath, fileName)
		if err != nil {
			s.Logger.Error("Failed to join remote paths: %v", err)
			return resp.FailWithMsg(err.Error())
		}

		fileSize, err := sftp.GetFileSize(localPath)
		if err != nil {
			s.Logger.Error("Failed to get file size: %v", err)
		}

		err = s.SFTPHandler.UploadFile(localPath, remoteFilePath, func(transferred, total int64) {
			progress := float64(transferred) / float64(total) * 100
			s.Logger.Debug("Upload progress: %.2f%% (%d/%d)", progress, transferred, total)

			totalProgress := (float64(completedFiles*100) + progress) / float64(totalFiles)

			runtime.EventsEmit(s.AppContext.Context(), "transfer:progress", map[string]interface{}{
				"type":          "upload",
				"fileName":      fileName,
				"fileSize":      fileSize,
				"progress":      progress,
				"totalFiles":    totalFiles,
				"completed":     completedFiles,
				"total":         total,
				"transferred":   transferred,
				"totalProgress": totalProgress,
			})
		})

		if err != nil {
			s.Logger.Error("Failed to upload file: %v", err)
			return resp.FailWithMsg(err.Error())
		}

		completedFiles++
	}

	s.Logger.Info("Successfully uploaded %d files", len(localPaths))
	return resp.OkWithCode(messages.UploadSuccess)
}

func (s *FileTransferSrv) DownloadFiles(remotePaths []string, localPath string) *resp.Resp {
	s.Logger.Info("Downloading files, local path: %s", localPath)

	s.SFTPHandlerMutex.Lock()
	defer s.SFTPHandlerMutex.Unlock()

	if s.SFTPHandler == nil || !s.SFTPHandler.IsConnected {
		s.Logger.Error("Not connected to SFTP server")
		return resp.FailWithMsg("Not connected to SFTP server")
	}

	processedPaths := make([]string, len(remotePaths))
	for i, remotePath := range remotePaths {
		processedPath, err := s.SFTPHandler.ProcessPath(remotePath)
		if err != nil {
			s.Logger.Error("Failed to process path %s: %v", remotePath, err)
			return resp.FailWithMsg(err.Error())
		}
		processedPaths[i] = processedPath
		s.Logger.Info("Processed download path: %s -> %s", remotePath, processedPath)
	}

	totalFiles := len(processedPaths)
	completedFiles := 0

	for _, remotePath := range processedPaths {
		s.Logger.Info("Downloading file: %s -> %s", remotePath, localPath)

		fileName := sftp.GetFileName(remotePath)
		localFilePath := sftp.JoinPath(localPath, fileName)

		fileSize, err := s.SFTPHandler.GetRemoteFileSize(remotePath)
		if err != nil {
			s.Logger.Error("Failed to get remote file size: %v", err)
		}

		err = s.SFTPHandler.DownloadFile(remotePath, localFilePath, func(transferred, total int64) {
			progress := float64(transferred) / float64(total) * 100

			totalProgress := (float64(completedFiles*100) + progress) / float64(totalFiles)

			runtime.EventsEmit(s.AppContext.Context(), "transfer:progress", map[string]interface{}{
				"type":          "download",
				"fileName":      fileName,
				"fileSize":      fileSize,
				"progress":      progress,
				"totalFiles":    totalFiles,
				"completed":     completedFiles,
				"total":         total,
				"transferred":   transferred,
				"totalProgress": totalProgress,
			})
		})

		if err != nil {
			s.Logger.Error("Failed to download file: %v", err)
			return resp.FailWithMsg(err.Error())
		}

		completedFiles++
	}

	s.Logger.Info("Successfully downloaded %d files", len(processedPaths))
	return resp.OkWithCode(messages.DownloadSuccess)
}

func (s *FileTransferSrv) CreateRemoteFolder(path string) *resp.Resp {
	s.Logger.Info("Creating remote folder: %s", path)

	s.SFTPHandlerMutex.Lock()
	defer s.SFTPHandlerMutex.Unlock()

	if s.SFTPHandler == nil || !s.SFTPHandler.IsConnected {
		s.Logger.Error("Not connected to SFTP server")
		return resp.FailWithMsg("Not connected to SFTP server")
	}

	if path == "" {
		return resp.FailWithMsg("Invalid empty path")
	}

	processedPath, err := s.SFTPHandler.ProcessPath(path)
	if err != nil {
		s.Logger.Error("Failed to process path %s: %v", path, err)
		return resp.FailWithMsg(err.Error())
	}
	s.Logger.Info("Using processed path for folder creation: %s", processedPath)

	err = s.SFTPHandler.CreateRemoteFolder(processedPath)
	if err != nil {
		s.Logger.Error("Failed to create remote folder: %v", err)
		return resp.FailWithMsg(err.Error())
	}

	s.Logger.Info("Remote folder created successfully")
	return resp.OkWithCode(messages.CreateFolderSuccess)
}

func (s *FileTransferSrv) SelectDownloadDirectory(title string) *resp.Resp {
	s.Logger.Info("Opening directory selection dialog with title: %s", title)

	directory, err := runtime.OpenDirectoryDialog(s.AppContext.Context(), runtime.OpenDialogOptions{
		Title: title,
	})

	if err != nil {
		s.Logger.Error("Failed to open directory selection dialog: %v", err)
		return resp.FailWithMsg(err.Error())
	}

	if directory == "" {
		s.Logger.Info("User canceled directory selection")
		return resp.OkWithCode("file_transfer.user_canceled")
	}

	s.Logger.Info("User selected directory: %s", directory)
	return resp.OkWithData(directory)
}

func (s *FileTransferSrv) SelectUploadFiles(title string) *resp.Resp {
	s.Logger.Info("Opening file selection dialog with title: %s", title)

	files, err := runtime.OpenMultipleFilesDialog(s.AppContext.Context(), runtime.OpenDialogOptions{
		Title:                title,
		CanCreateDirectories: true,
		ResolvesAliases:      true,
		ShowHiddenFiles:      true,
	})

	if err != nil {
		s.Logger.Error("Failed to open file selection dialog: %v", err)
		return resp.FailWithMsg(err.Error())
	}

	if len(files) == 0 {
		s.Logger.Info("User canceled file selection")
		return resp.OkWithCode("file_transfer.user_canceled")
	}

	s.Logger.Info("User selected %d files: %v", len(files), files)
	return resp.OkWithData(files)
}

func (s *FileTransferSrv) GenFileTransferBinding(
	a *types.FileTransferItemInfo,
	b *types.FileList,
	c *types.FileTransferTask,
) {
}
