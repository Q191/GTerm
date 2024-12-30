package adapter

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/OpenToolkitLab/GTerm/backend/enums"
	"github.com/OpenToolkitLab/GTerm/backend/pkg/terminal"
	"github.com/gorilla/websocket"
	"go.uber.org/zap"
	"golang.org/x/crypto/ssh"
	"io"
	"net"
	"sync"
	"time"
)

type SSHConfig struct {
	AuthType    enums.AuthType
	Port        uint
	Host        string
	User        string
	Password    string
	PrivateKey  string
	KeyPassword string
}

type SSH struct {
	conf      *SSHConfig
	ws        *websocket.Conn
	session   *ssh.Session
	stdinPipe io.WriteCloser
	writer    *writer
	logger    *zap.Logger
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
func (w *writer) Reset() {
	w.mu.Lock()
	defer w.mu.Unlock()
	w.buffer.Reset()
}

func NewSSH(conf *SSHConfig, ws *websocket.Conn, logger *zap.Logger) *SSH {
	return &SSH{
		conf:   conf,
		ws:     ws,
		logger: logger,
		writer: new(writer),
	}
}

func (s *SSH) Connect() (*SSH, error) {
	var auth []ssh.AuthMethod

	switch s.conf.AuthType {
	case enums.Password:
		auth = append(auth, ssh.Password(s.conf.Password))
		auth = append(auth, ssh.KeyboardInteractive(func(user, instruction string, questions []string, echos []bool) ([]string, error) {
			answers := make([]string, len(questions))
			if len(questions) == 1 {
				answers[0] = s.conf.Password
			}
			return answers, nil
		}))
	case enums.PrivateKey:
		signer, err := s.signer()
		if err != nil {
			return s, err
		}
		auth = append(auth, ssh.PublicKeys(signer))
	default:
		return s, errors.New("invalid authentication type provided")
	}

	config := &ssh.ClientConfig{
		User:    s.conf.User,
		Auth:    auth,
		Timeout: 10 * time.Second,
		HostKeyCallback: func(hostname string, remote net.Addr, key ssh.PublicKey) error {
			return nil
		},
	}

	client, err := ssh.Dial("tcp", fmt.Sprintf("%s:%d", s.conf.Host, s.conf.Port), config)
	if err != nil {
		return s, err
	}
	session, err := client.NewSession()
	if err != nil {
		return s, err
	}
	s.session = session

	s.stdinPipe, err = s.session.StdinPipe()
	if err != nil {
		return nil, err
	}

	s.session.Stdout = s.writer
	s.session.Stderr = s.writer

	modes := ssh.TerminalModes{
		ssh.ECHO:          1,
		ssh.TTY_OP_ISPEED: 14400,
		ssh.TTY_OP_OSPEED: 14400,
	}

	if err = s.session.RequestPty("xterm", 0, 0, modes); err != nil {
		return nil, err
	}

	if err = s.session.Shell(); err != nil {
		return nil, err
	}

	return s, nil
}

func (s *SSH) signer() (ssh.Signer, error) {
	if s.conf.KeyPassword == "" {
		return ssh.ParsePrivateKey([]byte(s.conf.PrivateKey))
	} else {
		return ssh.ParsePrivateKeyWithPassphrase([]byte(s.conf.PrivateKey), []byte(s.conf.KeyPassword))
	}
}

func (s *SSH) flushWriter() {
	if len(s.writer.Bytes()) != 0 {
		if err := s.ws.WriteMessage(websocket.BinaryMessage, s.writer.Bytes()); err != nil {
			s.logger.Error("failed write data to websocket", zap.Error(err))
		}
		s.writer.Reset()
	}
}

func (s *SSH) Input(quitSignal chan bool) {
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
			case terminal.Resize:
				if msg.Cols > 0 && msg.Rows > 0 {
					if err = s.session.WindowChange(msg.Rows, msg.Cols); err != nil {
						s.logger.Error("failed change ssh pty window size", zap.Error(err))
					}
				}
			case terminal.Command:
				if _, err = s.stdinPipe.Write([]byte(msg.Cmd)); err != nil {
					s.logger.Error("failed write command to stdin pipe", zap.Error(err))
				}
			}
		}
	}
}

func (s *SSH) Output(quitSignal chan bool) {
	defer s.setQuit(quitSignal)
	tick := time.NewTicker(time.Millisecond * time.Duration(60))
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
