package adapter

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"sync"
	"time"

	"github.com/Q191/GTerm/backend/enums"
	"github.com/Q191/GTerm/backend/initialize"
	commonssh "github.com/Q191/GTerm/backend/pkg/ssh"
	"github.com/Q191/GTerm/backend/pkg/terminal"
	"github.com/Q191/GTerm/backend/types"
	"github.com/gorilla/websocket"
	"golang.org/x/crypto/ssh"
)

type SSH struct {
	conf      *commonssh.Config
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

func NewSSH(conf *commonssh.Config, ws *websocket.Conn, logger initialize.Logger) *SSH {
	return &SSH{
		conf:   conf,
		ws:     ws,
		logger: logger,
		writer: new(writer),
	}
}

func (s *SSH) Connect() (*SSH, error) {
	host := fmt.Sprintf("%s:%d", s.conf.Host, s.conf.Port)
	s.logger.Info("Attempting to connect SSH, host: %s, port: %d", s.conf.Host, s.conf.Port)
	client, err := commonssh.NewSSHClient(s.conf, s.logger)
	if err != nil {
		var fingerprintErr *types.FingerprintError
		if errors.As(err, &fingerprintErr) {
			return s, &types.FingerprintError{
				Host:        fingerprintErr.Host,
				Fingerprint: fingerprintErr.Fingerprint,
			}
		}
		s.logger.Error("SSH connection failed: %v", err)
		return s, err
	}
	s.logger.Info("SSH connection successful, %s@%s", s.conf.User, host)

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

	// TODO: 支持自定义终端类型
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
