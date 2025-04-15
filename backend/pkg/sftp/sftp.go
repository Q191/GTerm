package sftp

import (
	"errors"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/MisakaTAT/GTerm/backend/initialize"
	"github.com/MisakaTAT/GTerm/backend/pkg/exec"
	commonssh "github.com/MisakaTAT/GTerm/backend/pkg/ssh"
	"github.com/MisakaTAT/GTerm/backend/types"

	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
)

type Handler struct {
	Logger           initialize.Logger
	SSHClient        *ssh.Client
	SFTPClient       *sftp.Client
	IsConnected      bool
	LastErrorTime    time.Time
	HomeDir          string
	PermissionsCache *PermissionsCache
	execAdapter      *exec.Adapter
}

func NewSFTPHandler(logger initialize.Logger) *Handler {
	return &Handler{
		Logger:           logger,
		IsConnected:      false,
		PermissionsCache: NewPermissionsCache(logger),
	}
}

func (h *Handler) Connect(conf *commonssh.Config) error {
	if h.IsConnected {
		return errors.New("already connected to SFTP server")
	}

	h.Logger.Info("Connecting to SFTP server %s:%d", conf.Host, conf.Port)
	client, err := commonssh.NewSSHClient(conf, h.Logger)
	if err != nil {
		h.Logger.Error("SSH connection failed: %v", err)
		return err
	}
	h.SSHClient = client

	h.execAdapter = exec.New(h.SSHClient)
	// 预加载权限信息
	h.PermissionsCache.preloadPermissions(h.SSHClient, h.execAdapter)

	h.Logger.Info("Initializing SFTP client")
	sftpClient, err := sftp.NewClient(client)
	if err != nil {
		if err = h.SSHClient.Close(); err != nil {
			h.Logger.Error("Failed to close SFTP client: %v", err)
		}
		h.Logger.Error("SFTP client initialization failed: %v", err)
		return err
	}

	h.SFTPClient = sftpClient
	h.IsConnected = true
	h.Logger.Info("SFTP connection successful")

	homeDir, err := h.GetHomeDirectory()
	if err != nil {
		h.Logger.Warn("Failed to get home directory: %v, using / as default", err)
		h.HomeDir = "/"
	} else {
		h.HomeDir = homeDir
		h.Logger.Info("User home directory: %s", homeDir)
	}

	return nil
}

func (h *Handler) Close() error {
	if !h.IsConnected {
		return nil
	}

	var err error
	if h.SFTPClient != nil {
		if err = h.SFTPClient.Close(); err != nil {
			h.Logger.Error("Failed to close SFTP client: %v", err)
		}
		h.SFTPClient = nil
	}

	if h.SSHClient != nil {
		if err = h.SSHClient.Close(); err != nil {
			h.Logger.Error("Failed to close SSH client: %v", err)
		}
		h.SSHClient = nil
	}

	h.IsConnected = false
	h.Logger.Info("SFTP connection closed")
	return err
}

func (h *Handler) ListRemoteFiles(path string) ([]*types.FileTransferItemInfo, error) {
	if err := h.checkSFTPConnection(); err != nil {
		return nil, err
	}

	h.Logger.Info("Listing remote directory: %s", path)
	entries, err := h.SFTPClient.ReadDir(path)
	if err != nil {
		h.Logger.Error("Failed to list remote directory: %v", err)
		return nil, err
	}

	files := make([]*types.FileTransferItemInfo, 0, len(entries))
	for _, entry := range entries {
		info := &types.FileTransferItemInfo{
			Name:        entry.Name(),
			Size:        entry.Size(),
			IsDir:       entry.IsDir(),
			ModTime:     entry.ModTime().Format(time.RFC3339),
			Permissions: entry.Mode().String(),
			Owner:       "unknown",
			Group:       "unknown",
		}

		if stat, ok := entry.Sys().(*sftp.FileStat); ok {
			info.Owner = h.PermissionsCache.GetUsername(stat.UID)
			info.Group = h.PermissionsCache.GetGroupName(stat.GID)
		}
		files = append(files, info)
	}

	h.Logger.Info("Found %d files/directories in remote directory %s", len(files), path)
	return files, nil
}

func (h *Handler) UploadFile(localPath, remotePath string, progressCallback func(int64, int64)) error {
	if err := h.checkSFTPConnection(); err != nil {
		return err
	}

	localFile, err := os.Open(localPath)
	if err != nil {
		h.Logger.Error("Failed to open local file: %v", err)
		return err
	}
	defer func(localFile *os.File) {
		if err = localFile.Close(); err != nil {
			h.Logger.Error("Failed to close local file: %v", err)
		}
	}(localFile)

	localFileInfo, err := localFile.Stat()
	if err != nil {
		h.Logger.Error("Failed to get local file info: %v", err)
		return err
	}

	if localFileInfo.IsDir() {
		h.Logger.Error("Directory upload not supported")
		return errors.New("directory upload not supported")
	}

	totalSize := localFileInfo.Size()
	h.Logger.Info("Uploading file: %s -> %s, size: %d bytes", localPath, remotePath, totalSize)

	remoteDir := filepath.Dir(remotePath)
	if remoteDir != "." && remoteDir != "/" {
		err = h.SFTPClient.MkdirAll(remoteDir)
		if err != nil {
			h.Logger.Error("Failed to create remote directory: %v", err)
			return err
		}
	}

	remoteFile, err := h.SFTPClient.Create(remotePath)
	if err != nil {
		h.Logger.Error("Failed to create remote file: %v", err)
		return err
	}

	defer func(remoteFile *sftp.File) {
		if err = remoteFile.Close(); err != nil {
			h.Logger.Error("Failed to close remote file: %v", err)
		}
	}(remoteFile)

	_, err = io.Copy(localFile, &ProgressReader{
		Reader:           remoteFile,
		TotalSize:        totalSize,
		BytesRead:        0,
		ProgressCallback: progressCallback,
	})

	if err != nil {
		h.Logger.Error("File upload failed: %v", err)
		return err
	}

	h.Logger.Info("File upload successful: %s -> %s", localPath, remotePath)
	return nil
}

func (h *Handler) DownloadFile(remotePath, localPath string, progressCallback func(int64, int64)) error {
	if err := h.checkSFTPConnection(); err != nil {
		return err
	}

	remoteFileInfo, err := h.SFTPClient.Stat(remotePath)
	if err != nil {
		h.Logger.Error("Failed to get remote file info: %v", err)
		return err
	}

	if remoteFileInfo.IsDir() {
		h.Logger.Error("Directory download not supported")
		return errors.New("directory download not supported")
	}

	totalSize := remoteFileInfo.Size()
	h.Logger.Info("Downloading file: %s -> %s, size: %d bytes", remotePath, localPath, totalSize)

	remoteFile, err := h.SFTPClient.Open(remotePath)
	if err != nil {
		h.Logger.Error("Failed to open remote file: %v", err)
		return err
	}
	defer func(remoteFile *sftp.File) {
		if err = remoteFile.Close(); err != nil {
			h.Logger.Error("Failed to close remote file: %v", err)
		}
	}(remoteFile)

	localDir := filepath.Dir(localPath)
	if localDir != "." {
		err = os.MkdirAll(localDir, 0755)
		if err != nil {
			h.Logger.Error("Failed to create local directory: %v", err)
			return err
		}
	}

	localFile, err := os.Create(localPath)
	if err != nil {
		h.Logger.Error("Failed to create local file: %v", err)
		return err
	}
	defer func(localFile *os.File) {
		if err = localFile.Close(); err != nil {
			h.Logger.Error("Failed to close local file: %v", err)
		}
	}(localFile)

	_, err = io.Copy(localFile, &ProgressReader{
		Reader:           remoteFile,
		TotalSize:        totalSize,
		BytesRead:        0,
		ProgressCallback: progressCallback,
	})

	if err != nil {
		h.Logger.Error("File download failed: %v", err)
		return err
	}

	h.Logger.Info("File download successful: %s -> %s", remotePath, localPath)
	return nil
}

func (h *Handler) CreateRemoteFolder(path string) error {
	if err := h.checkSFTPConnection(); err != nil {
		return err
	}

	h.Logger.Info("Creating remote folder: %s", path)
	err := h.SFTPClient.MkdirAll(path)
	if err != nil {
		h.Logger.Error("Failed to create remote folder: %v", err)
		return err
	}

	h.Logger.Info("Remote folder created successfully: %s", path)
	return nil
}

func (h *Handler) GetRemoteFileSize(path string) (int64, error) {
	if err := h.checkSFTPConnection(); err != nil {
		return 0, err
	}

	fileInfo, err := h.SFTPClient.Stat(path)
	if err != nil {
		return 0, err
	}
	return fileInfo.Size(), nil
}

func (h *Handler) GetHomeDirectory() (string, error) {
	if err := h.checkSSHConnection(); err != nil {
		return "", err
	}

	result := h.execAdapter.Run("pwd")
	if !result.Success() {
		h.Logger.Error("Failed to execute pwd command: %v", result.Error())
	}

	homeDir := result.Unwrap()
	if homeDir == "" {
		h.Logger.Info("pwd command failed, trying $HOME environment variable")
		result = h.execAdapter.Run("echo $HOME")
		if !result.Success() {
			h.Logger.Error("Failed to get HOME environment variable: %v", result.Error())
		}
		homeDir = result.Unwrap()
	}

	if homeDir == "" {
		h.Logger.Warn("Failed to determine home directory, using / as default")
		return "/", nil
	}

	if !strings.HasPrefix(homeDir, "/") {
		homeDir = "/" + homeDir
	}

	if strings.HasSuffix(homeDir, "/") && len(homeDir) > 1 {
		homeDir = homeDir[:len(homeDir)-1]
	}

	h.Logger.Info("Home directory determined: %s", homeDir)
	return homeDir, nil
}

func (h *Handler) ProcessPath(path string) (string, error) {
	if err := h.checkSFTPConnection(); err != nil {
		return "", err
	}

	if path == "" {
		if h.HomeDir != "" {
			return h.HomeDir, nil
		}
		return "/", nil
	}
	if path == "/" {
		return "/", nil
	}
	if path[0] != '/' {
		path = "/" + path
	}

	return path, nil
}

func (h *Handler) JoinRemotePaths(base, relPath string) (string, error) {
	basePath, err := h.ProcessPath(base)
	if err != nil {
		return "", err
	}
	if relPath == "" {
		return basePath, nil
	}
	if relPath[0] == '/' {
		return relPath, nil
	}
	if basePath[len(basePath)-1] != '/' {
		basePath += "/"
	}

	return basePath + relPath, nil
}

func (h *Handler) isSFTPReady() bool {
	return h.IsConnected && h.SFTPClient != nil
}

func (h *Handler) isSSHReady() bool {
	return h.IsConnected && h.SSHClient != nil
}

func (h *Handler) checkSFTPConnection() error {
	if !h.isSFTPReady() {
		return errors.New("not connected to SFTP server")
	}
	return nil
}

func (h *Handler) checkSSHConnection() error {
	if !h.isSSHReady() {
		return errors.New("not connected to SSH server")
	}
	return nil
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
