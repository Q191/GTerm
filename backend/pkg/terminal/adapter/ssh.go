package adapter

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/MisakaTAT/GTerm/backend/enums"
	"github.com/MisakaTAT/GTerm/backend/initialize"
	"github.com/MisakaTAT/GTerm/backend/pkg/terminal"
	"github.com/MisakaTAT/GTerm/backend/types"
	"github.com/gorilla/websocket"
	"github.com/skeema/knownhosts"
	"golang.org/x/crypto/ssh"
)

type SSHConfig struct {
	AuthMethod       enums.AuthMethod
	Port             uint
	Host             string
	User             string
	Password         string
	PrivateKey       string
	KeyPassword      string
	TrustUnknownHost bool
}

type SSH struct {
	conf      *SSHConfig
	ws        *websocket.Conn
	session   *ssh.Session
	stdinPipe io.WriteCloser
	writer    *writer
	logger    initialize.Logger
}

type writer struct {
	buffer bytes.Buffer
	mu     sync.Mutex
}

func (w *writer) Write(p []byte) (int, error) {
	w.mu.Lock()
	defer w.mu.Unlock()
	return w.buffer.Write(p)
}

func (w *writer) Bytes() []byte {
	w.mu.Lock()
	defer w.mu.Unlock()
	return w.buffer.Bytes()
}

func (w *writer) String() string {
	w.mu.Lock()
	defer w.mu.Unlock()
	return w.buffer.String()
}

func (w *writer) Reset() {
	w.mu.Lock()
	defer w.mu.Unlock()
	w.buffer.Reset()
}

func NewSSH(conf *SSHConfig, ws *websocket.Conn, logger initialize.Logger) *SSH {
	return &SSH{
		conf:   conf,
		ws:     ws,
		logger: logger,
		writer: new(writer),
	}
}

func (s *SSH) Connect() (*SSH, error) {
	hostPort := fmt.Sprintf("%s:%d", s.conf.Host, s.conf.Port)
	s.logger.Info("Attempting to connect SSH, host: %s, port: %d", s.conf.Host, s.conf.Port)

	var auth []ssh.AuthMethod

	switch s.conf.AuthMethod {
	case enums.Password:
		auth = append(auth, ssh.Password(s.conf.Password))
		auth = append(auth, ssh.KeyboardInteractive(func(user, instruction string, questions []string, echos []bool) ([]string, error) {
			answers := make([]string, len(questions))
			if len(questions) == 1 {
				answers[0] = s.conf.Password
			}
			return answers, nil
		}))
		s.logger.Info("Using password authentication")
	case enums.PrivateKey:
		signer, err := s.signer()
		if err != nil {
			s.logger.Error("Failed to parse private key: %v", err)
			return s, err
		}
		auth = append(auth, ssh.PublicKeys(signer))
		s.logger.Info("Using private key authentication")
	default:
		s.logger.Error("Invalid authentication type: %v", s.conf.AuthMethod)
		return s, errors.New("invalid authentication type provided")
	}

	var hostKeyCallback ssh.HostKeyCallback
	var hostKeyAlgorithms []string
	knownHostsFile := filepath.Join(os.Getenv("HOME"), ".ssh", "known_hosts")
	s.logger.Debug("Using known_hosts file: %s", knownHostsFile)

	if s.conf.TrustUnknownHost {
		hostKeyCallback = ssh.InsecureIgnoreHostKey()
		s.logger.Warn("Configured to trust all unknown hosts, this may pose a security risk")
	} else {
		db, err := knownhosts.NewDB(knownHostsFile)
		if err != nil {
			s.logger.Error("Failed to load known_hosts database: %v", err)
			return s, err
		}

		hostKeyCallback = db.HostKeyCallback()
		hostKeyAlgorithms = db.HostKeyAlgorithms(hostPort)
		s.logger.Info("Using known_hosts for host key verification")
	}

	config := &ssh.ClientConfig{
		User:            s.conf.User,
		Auth:            auth,
		Timeout:         10 * time.Second,
		HostKeyCallback: hostKeyCallback,
	}

	// 优先从 known_hosts 获取算法列表
	if len(hostKeyAlgorithms) > 0 {
		config.HostKeyAlgorithms = hostKeyAlgorithms
		s.logger.Info("Using host key algorithms from known_hosts: %v", hostKeyAlgorithms)
	} else {
		// 默认算法列表
		config.HostKeyAlgorithms = []string{
			ssh.KeyAlgoRSASHA512,
			ssh.KeyAlgoRSASHA256,
			ssh.KeyAlgoRSA,
			ssh.KeyAlgoECDSA256,
			ssh.KeyAlgoED25519,
		}
		s.logger.Info("Using default host key algorithm list")
	}

	s.logger.Info("Starting SSH connection to server, %s@%s", s.conf.User, hostPort)
	client, err := ssh.Dial("tcp", hostPort, config)
	if err != nil {
		if knownhosts.IsHostUnknown(err) {
			s.logger.Info("Unknown host, attempting to get host key: %s", hostPort)
			key, keyErr := s.getHostKey(hostPort)
			if keyErr != nil {
				s.logger.Error("Failed to get host key: %v", keyErr)
				return s, keyErr
			}
			fingerprint := ssh.FingerprintSHA256(key)
			s.logger.Info("Successfully obtained unknown host key, fingerprint: %s", fingerprint)
			return s, &types.FingerprintError{
				Host:        hostPort,
				Fingerprint: fingerprint,
			}
		}
		// 如果主机密钥已更改（可能存在中间人攻击）
		if knownhosts.IsHostKeyChanged(err) {
			s.logger.Warn("Host key has changed! This may indicate a MitM attack, host: %s", hostPort)
		}
		s.logger.Error("SSH connection failed: %v", err)
		return s, err
	}
	s.logger.Info("SSH connection successful, %s@%s", s.conf.User, hostPort)

	s.logger.Info("Creating SSH session")
	session, err := client.NewSession()
	if err != nil {
		s.logger.Error("Failed to create SSH session: %v", err)
		return s, err
	}
	s.session = session

	s.logger.Debug("Getting session stdin pipe")
	s.stdinPipe, err = s.session.StdinPipe()
	if err != nil {
		s.logger.Error("Failed to get stdin pipe: %v", err)
		return nil, err
	}

	s.session.Stdout = s.writer
	s.session.Stderr = s.writer
	s.logger.Debug("Stdout and stderr configured")

	modes := ssh.TerminalModes{
		ssh.ECHO:          1,
		ssh.TTY_OP_ISPEED: 14400,
		ssh.TTY_OP_OSPEED: 14400,
	}

	s.logger.Debug("Requesting PTY terminal, type: xterm")
	if err = s.session.RequestPty("xterm", 0, 0, modes); err != nil {
		s.logger.Error("Failed to request PTY terminal: %v", err)
		return nil, err
	}

	s.logger.Debug("Starting shell")
	if err = s.session.Shell(); err != nil {
		s.logger.Error("Failed to start shell: %v", err)
		return nil, err
	}

	s.logger.Info("SSH session ready")
	return s, nil
}

func (s *SSH) createHostKeyCallback(knownHostsFile string) (ssh.HostKeyCallback, error) {
	dir := filepath.Dir(knownHostsFile)
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		if err = os.MkdirAll(dir, 0700); err != nil {
			return nil, err
		}
	}
	if _, err := os.Stat(knownHostsFile); os.IsNotExist(err) {
		if _, err = os.Create(knownHostsFile); err != nil {
			return nil, err
		}
	}
	db, err := knownhosts.NewDB(knownHostsFile)
	if err != nil {
		return nil, err
	}
	return db.HostKeyCallback(), nil
}

func (s *SSH) getHostKey(hostPort string) (ssh.PublicKey, error) {
	var hostKey ssh.PublicKey
	keyCallback := func(hostname string, remote net.Addr, key ssh.PublicKey) error {
		hostKey = key
		return nil
	}
	config := &ssh.ClientConfig{
		User:            s.conf.User,
		Auth:            []ssh.AuthMethod{ssh.Password(s.conf.Password)},
		HostKeyCallback: keyCallback,
		Timeout:         5 * time.Second,
	}
	conn, err := ssh.Dial("tcp", hostPort, config)
	if err != nil {
		if hostKey != nil {
			s.logger.Info("Successfully obtained host key, fingerprint: %s", ssh.FingerprintSHA256(hostKey))
			return hostKey, nil
		}
		return nil, fmt.Errorf("无法获取主机密钥")
	}
	defer conn.Close()

	if hostKey == nil {
		return nil, fmt.Errorf("无法获取主机密钥")
	}

	s.logger.Info("Successfully obtained host key, fingerprint: %s", ssh.FingerprintSHA256(hostKey))
	return hostKey, nil
}

func (s *SSH) AddFingerprint(hostPort, fingerprint string) error {
	s.logger.Info("Adding host fingerprint, host: %s, fingerprint: %s", hostPort, fingerprint)
	knownHostsFile := filepath.Join(os.Getenv("HOME"), ".ssh", "known_hosts")

	key, err := s.getHostKey(hostPort)
	if err != nil {
		return err
	}

	// 验证指纹是否匹配
	actualFingerprint := ssh.FingerprintSHA256(key)
	s.logger.Info("Comparing fingerprints, expected: %s, actual: %s", fingerprint, actualFingerprint)
	if actualFingerprint != fingerprint {
		return fmt.Errorf("指纹不匹配: 期望 %s, 实际 %s", fingerprint, actualFingerprint)
	}

	// 添加到 known_hosts 文件
	f, err := os.OpenFile(knownHostsFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0600)
	if err != nil {
		s.logger.Error("Failed to open known_hosts file: %v", err)
		return err
	}
	defer func(f *os.File) {
		_ = f.Close()
	}(f)

	hostname := strings.Split(hostPort, ":")[0]
	port := int(s.conf.Port)

	var remoteAddr net.Addr
	if ip := net.ParseIP(hostname); ip != nil {
		remoteAddr = &net.TCPAddr{IP: ip, Port: port}
	} else {
		// 对于非IP地址的主机名，使用伪造的地址
		remoteAddr = &net.TCPAddr{IP: net.IPv4(127, 0, 0, 1), Port: port}
	}

	err = knownhosts.WriteKnownHost(f, hostPort, remoteAddr, key)
	if err != nil {
		s.logger.Error("Failed to write to known_hosts file: %v", err)
		return err
	}

	s.logger.Info("Successfully added fingerprint to known_hosts, fingerprint: %s", fingerprint)
	return nil
}

func (s *SSH) signer() (ssh.Signer, error) {
	s.logger.Debug("Parsing SSH private key")
	if s.conf.KeyPassword == "" {
		return ssh.ParsePrivateKey([]byte(s.conf.PrivateKey))
	} else {
		return ssh.ParsePrivateKeyWithPassphrase([]byte(s.conf.PrivateKey), []byte(s.conf.KeyPassword))
	}
}

func (s *SSH) flushWriter() {
	if len(s.writer.String()) != 0 {
		if err := s.ws.WriteJSON(&types.Message{
			Type:    enums.TerminalTypeData,
			Content: s.writer.String(),
		}); err != nil {
			s.logger.Error("failed write data to websocket: %v", err)
		}
		s.writer.Reset()
	}
}

func (s *SSH) Input(quitSignal chan bool) {
	s.logger.Info("Starting WebSocket input monitoring")
	defer s.setQuit(quitSignal)

	for {
		select {
		case <-quitSignal:
			return
		default:
			_, data, err := s.ws.ReadMessage()
			if err != nil {
				return
			}
			msg := &terminal.Payload{}
			_ = json.Unmarshal(data, &msg)

			switch msg.Type {
			case enums.TerminalTypeResize:
				if msg.Cols > 0 && msg.Rows > 0 {
					if err = s.session.WindowChange(msg.Rows, msg.Cols); err != nil {
						s.logger.Error("failed change ssh pty window size: %v", err)
					}
				}
			case enums.TerminalTypeCMD:
				if _, err = s.stdinPipe.Write([]byte(msg.Cmd)); err != nil {
					s.logger.Error("failed write command to stdin pipe: %v", err)
				}
			}
		}
	}
}

func (s *SSH) Output(quitSignal chan bool) {
	s.logger.Info("Starting WebSocket output")
	defer s.setQuit(quitSignal)
	tick := time.NewTicker(time.Millisecond * time.Duration(5))
	defer tick.Stop()
	for {
		select {
		case <-quitSignal:
			s.flushWriter()
			return
		case <-tick.C:
			s.flushWriter()
		}
	}
}

func (s *SSH) close() {
	s.logger.Info("Closing SSH session")
	if s.session != nil {
		_ = s.session.Close()
	}
}

func (s *SSH) Wait(quitSignal chan bool) {
	defer s.close()
	defer s.setQuit(quitSignal)
	_ = s.session.Wait()
}

func (s *SSH) setQuit(ch chan bool) {
	ch <- true
}
