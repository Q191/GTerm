package services

import (
	"errors"
	"fmt"
	"github.com/OpenToolkitLab/GTerm/backend/consts"
	"github.com/OpenToolkitLab/GTerm/backend/enums"
	"github.com/OpenToolkitLab/GTerm/backend/initialize"
	"github.com/OpenToolkitLab/GTerm/backend/pkg/adapter"
	"github.com/OpenToolkitLab/GTerm/backend/pkg/terminal"
	"github.com/google/wire"
	"github.com/gorilla/websocket"
	"go.uber.org/zap"
	"net/http"
	"time"
)

var TerminalSrvSet = wire.NewSet(wire.Struct(new(TerminalSrv), "*"))

type TerminalSrv struct {
	HTTPListenerPort *initialize.HTTPListenerPort
	Logger           *zap.Logger
}

var ug = websocket.Upgrader{
	ReadBufferSize:    1024,
	WriteBufferSize:   1024 * 1024 * 10,
	EnableCompression: true,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func (s *TerminalSrv) WebSocketPort() int {
	return int(*s.HTTPListenerPort)
}

func (s *TerminalSrv) closeWs(ws *websocket.Conn) {
	data := websocket.FormatCloseMessage(websocket.CloseNormalClosure, "Server closed")
	err := ws.WriteControl(websocket.CloseMessage, data, time.Now().Add(consts.WebSocketWriteWait))
	if err != nil && !errors.Is(err, websocket.ErrCloseSent) {
		// If close message could not be sent, then close without the handshake.
		_ = ws.Close()
	}
}

func (s *TerminalSrv) Startup(w http.ResponseWriter, r *http.Request) {
	ws, err := ug.Upgrade(w, r, nil)
	if err != nil {
		s.Logger.Error("Failed to upgrade connection", zap.Error(err))
		return
	}
	defer s.closeWs(ws)
	s.handleError(ws, s.SSH(ws))
}

func (s *TerminalSrv) handleError(ws *websocket.Conn, err error) {
	if err != nil && ws != nil {
		s.Logger.Error("Handle connect err", zap.Error(err))
		errMsg := fmt.Sprintf("\x1B[31m[ERR] %s\x1B[0m", err.Error())
		if err = ws.WriteMessage(websocket.BinaryMessage, []byte(errMsg)); err != nil {
			s.Logger.Error("Failed write err message", zap.Error(err))
		}
	}
}

func (s *TerminalSrv) SSH(ws *websocket.Conn) error {
	// test code
	sshConf := &adapter.SSHConfig{
		Host:     "192.168.100.77",
		Port:     22,
		User:     "root",
		AuthType: enums.Password,
		Password: "Admin@123",
	}

	termConf := &terminal.Config{
		WebSocket: ws,
	}

	ssh, err := adapter.NewSSH(sshConf, ws, s.Logger).Connect()
	if err != nil {
		return err
	}

	term := terminal.NewTerminal(termConf, ssh, s.closeWs)
	term.Start()

	return nil
}

func (s *TerminalSrv) Serial(ws *websocket.Conn) error {
	serial := adapter.NewSerial(ws, s.Logger)

	ports := serial.SerialPorts()
	if len(ports) == 0 {
		return errors.New("no serial ports available")
	}

	// test code
	serialPort := "/dev/cu.usbserial-2130"

	err := serial.Open(serialPort)
	if err != nil {
		return fmt.Errorf("failed to open serial port: %v", err)
	}

	termConf := &terminal.Config{
		WebSocket: ws,
	}

	term := terminal.NewTerminal(termConf, serial, s.closeWs)
	term.Start()

	return nil
}
