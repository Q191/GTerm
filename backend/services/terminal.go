package services

import (
	"errors"
	"fmt"
	"time"

	"github.com/MisakaTAT/GTerm/backend/consts"
	"github.com/MisakaTAT/GTerm/backend/consts/messages"
	"github.com/MisakaTAT/GTerm/backend/enums"
	"github.com/MisakaTAT/GTerm/backend/initialize"
	"github.com/MisakaTAT/GTerm/backend/pkg/terminal"
	"github.com/MisakaTAT/GTerm/backend/pkg/terminal/adapter"
	"github.com/MisakaTAT/GTerm/backend/types"
	"github.com/MisakaTAT/GTerm/backend/utils/resp"
	"github.com/google/wire"
	"github.com/gorilla/websocket"
	"go.bug.st/serial"
)

var TerminalSrvSet = wire.NewSet(wire.Struct(new(TerminalSrv), "*"))

type TerminalSrv struct {
	Logger           initialize.Logger
	ConnectionSrv    *ConnectionSrv
	MetadataSrv      *MetadataSrv
	HTTPListenerPort *initialize.HTTPListenerPort
}

func (s *TerminalSrv) SSH(ws *websocket.Conn, hostID uint) error {
	s.Logger.Info("Starting SSH connection, hostID: %d", hostID)
	conn, err := s.ConnectionSrv.FindByID(hostID)
	if err != nil {
		s.Logger.Error("Failed to find host information: %v, hostID: %d", err, hostID)
		return fmt.Errorf("failed to find host: %v", err)
	}

	s.Logger.Info("Found host information, host: %s, port: %d", conn.Host, conn.Port)
	if conn.Metadata == nil {
		s.Logger.Info("Host metadata is empty, starting metadata update")
		go s.MetadataSrv.UpdateByConnection(conn)
	}

	sshConf := &adapter.SSHConfig{
		Host:       conn.Host,
		Port:       conn.Port,
		User:       conn.Credential.Username,
		AuthMethod: conn.Credential.AuthMethod,
	}

	s.Logger.Info("SSH configuration ready, host: %s, user: %s, auth method: %s",
		conn.Host,
		conn.Credential.Username,
		conn.Credential.AuthMethod)

	switch conn.Credential.AuthMethod {
	case enums.Password:
		sshConf.Password = conn.Credential.Password
		s.Logger.Debug("Using password authentication")
	case enums.PrivateKey:
		sshConf.PrivateKey = conn.Credential.PrivateKey
		sshConf.KeyPassword = conn.Credential.Passphrase
		s.Logger.Debug("Using private key authentication")
	}

	s.Logger.Info("Connecting to SSH server, host: %s, port: %d", conn.Host, conn.Port)
	ssh, err := adapter.NewSSH(sshConf, ws, s.Logger).Connect()
	if err != nil {
		s.Logger.Error("SSH connection failed: %v, host: %s, port: %d", err, conn.Host, conn.Port)
		return err
	}
	s.Logger.Info("SSH connection successful, host: %s, port: %d", conn.Host, conn.Port)

	// 发送连接成功消息
	if err = ws.WriteJSON(&types.Message{Type: enums.TerminalTypeConnected}); err != nil {
		s.Logger.Error("Failed to send connection success message: %v", err)
		return err
	}
	s.Logger.Info("Connection success message sent")

	term := terminal.NewTerminal(ws, ssh, s.SessionEnded, s.Logger)
	s.Logger.Info("Starting terminal session, host: %s, port: %d", conn.Host, conn.Port)
	term.Start()

	return nil
}

func (s *TerminalSrv) AddFingerprint(hostID uint, hostAddress string, fingerprint string) error {
	s.Logger.Info("Adding host fingerprint, hostID: %d, address: %s, fingerprint: %s", hostID, hostAddress, fingerprint)
	conn, err := s.ConnectionSrv.FindByID(hostID)
	if err != nil {
		s.Logger.Error("Failed to find host information: %v, hostID: %d", err, hostID)
		return fmt.Errorf("failed to find host: %v", err)
	}

	sshConf := &adapter.SSHConfig{
		Host:       conn.Host,
		Port:       conn.Port,
		User:       conn.Credential.Username,
		AuthMethod: conn.Credential.AuthMethod,
	}

	switch conn.Credential.AuthMethod {
	case enums.Password:
		sshConf.Password = conn.Credential.Password
	case enums.PrivateKey:
		sshConf.PrivateKey = conn.Credential.PrivateKey
		sshConf.KeyPassword = conn.Credential.Passphrase
	}

	// 创建SSH配置
	ssh := adapter.NewSSH(sshConf, nil, s.Logger)

	// 添加主机指纹
	if err = ssh.AddFingerprint(hostAddress, fingerprint); err != nil {
		s.Logger.Error("Failed to add host fingerprint: %v, address: %s", err, hostAddress)
		return fmt.Errorf("failed to add host fingerprint: %v", err)
	}
	s.Logger.Info("Successfully added host fingerprint, address: %s", hostAddress)

	return nil
}

// func (s *TerminalSrv) Serial(ws *websocket.Conn) error {
// 	serial := adapter.NewSerial(ws, s.Logger)
//
// 	// test code
// 	serialPort := "/dev/cu.usbserial-2130"
//
// 	err := serial.Open(serialPort)
// 	if err != nil {
// 		return fmt.Errorf("failed to open serial port: %v", err)
// 	}
//
// 	term := terminal.NewTerminal(ws, serial, s.closeWsWrapper)
// 	term.Start()
//
// 	return nil
// }

func (s *TerminalSrv) SerialPorts() *resp.Resp {
	s.Logger.Info("Getting available serial ports")
	ports, err := serial.GetPortsList()
	if err != nil {
		s.Logger.Error("Failed to get serial port list: %v", err)
		return resp.FailWithMsg(err.Error())
	}
	s.Logger.Info("Found %d available serial ports", len(ports))
	return resp.OkWithData(ports)
}

func (s *TerminalSrv) CloseSession(ws *websocket.Conn, reason string) {
	s.Logger.Info("Closing session, reason: %s", reason)
	data := websocket.FormatCloseMessage(websocket.CloseNormalClosure, reason)
	err := ws.WriteControl(websocket.CloseMessage, data, time.Now().Add(consts.WebSocketWriteWait))
	if err != nil && !errors.Is(err, websocket.ErrCloseSent) {
		s.Logger.Error("Failed to close session: %v, forcibly closing connection", err)
		// If close message could not be sent, then close without the handshake.
		_ = ws.Close()
	}
}

func (s *TerminalSrv) SessionEnded(ws *websocket.Conn) {
	s.CloseSession(ws, messages.SessionEnded)
}

func (s *TerminalSrv) WebsocketPort() int {
	port := int(*s.HTTPListenerPort)
	s.Logger.Debug("WebSocket service port: %d", port)
	return port
}
