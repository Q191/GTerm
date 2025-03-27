package services

import (
	"errors"
	"fmt"
	"github.com/MisakaTAT/GTerm/backend/consts"
	"github.com/MisakaTAT/GTerm/backend/enums"
	"github.com/MisakaTAT/GTerm/backend/initialize"
	"github.com/MisakaTAT/GTerm/backend/pkg/terminal"
	"github.com/MisakaTAT/GTerm/backend/pkg/terminal/adapter"
	"github.com/MisakaTAT/GTerm/backend/types"
	"github.com/MisakaTAT/GTerm/backend/utils/resp"
	"github.com/google/wire"
	"github.com/gorilla/websocket"
	"go.bug.st/serial"
	"go.uber.org/zap"
	"time"
)

var TerminalSrvSet = wire.NewSet(wire.Struct(new(TerminalSrv), "*"))

type TerminalSrv struct {
	Logger           *zap.Logger
	ConnectionSrv    *ConnectionSrv
	MetadataSrv      *MetadataSrv
	HTTPListenerPort *initialize.HTTPListenerPort
}

func (s *TerminalSrv) SSH(ws *websocket.Conn, hostID uint) error {
	conn, err := s.ConnectionSrv.FindByID(hostID)
	if err != nil {
		return fmt.Errorf("failed to find host: %v", err)
	}

	if conn.Metadata == nil {
		go s.MetadataSrv.UpdateByConnection(conn)
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

	ssh, err := adapter.NewSSH(sshConf, ws, s.Logger).Connect()
	if err != nil {
		return err
	}

	// 发送连接成功消息
	if err = ws.WriteJSON(&types.Message{Type: types.MessageTypeConnected}); err != nil {
		s.Logger.Error("Failed to send connected message",
			zap.Error(err),
		)
		return err
	}

	term := terminal.NewTerminal(ws, ssh, s.SessionEnded)
	term.Start()

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
	ports, err := serial.GetPortsList()
	if err != nil {
		return resp.FailWithMsg(err.Error())
	}
	return resp.OkWithData(ports)
}

func (s *TerminalSrv) CloseSession(ws *websocket.Conn, reason string) {
	data := websocket.FormatCloseMessage(websocket.CloseNormalClosure, reason)
	err := ws.WriteControl(websocket.CloseMessage, data, time.Now().Add(consts.WebSocketWriteWait))
	if err != nil && !errors.Is(err, websocket.ErrCloseSent) {
		// If close message could not be sent, then close without the handshake.
		_ = ws.Close()
	}
}

func (s *TerminalSrv) SessionEnded(ws *websocket.Conn) {
	s.CloseSession(ws, "Session ended")
}

func (s *TerminalSrv) WebsocketPort() int {
	return int(*s.HTTPListenerPort)
}
