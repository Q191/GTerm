package services

import (
	"errors"
	"fmt"
	"github.com/MisakaTAT/GTerm/backend/consts"
	"github.com/MisakaTAT/GTerm/backend/enums"
	"github.com/MisakaTAT/GTerm/backend/initialize"
	"github.com/MisakaTAT/GTerm/backend/pkg/adapter"
	"github.com/MisakaTAT/GTerm/backend/pkg/terminal"
	"github.com/MisakaTAT/GTerm/backend/pkg/types"
	"github.com/google/wire"
	"github.com/gorilla/websocket"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"net/http"
	"strconv"
	"strings"
	"time"
)

var TerminalSrvSet = wire.NewSet(wire.Struct(new(TerminalSrv), "*"))

type TerminalSrv struct {
	DB               *gorm.DB
	HTTPListenerPort *initialize.HTTPListenerPort
	Logger           *zap.Logger
	HostSrv          *HostSrv
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

func (s *TerminalSrv) closeWs(ws *websocket.Conn, reason string) {
	data := websocket.FormatCloseMessage(websocket.CloseNormalClosure, reason)
	err := ws.WriteControl(websocket.CloseMessage, data, time.Now().Add(consts.WebSocketWriteWait))
	if err != nil && !errors.Is(err, websocket.ErrCloseSent) {
		// If close message could not be sent, then close without the handshake.
		_ = ws.Close()
	}
}

func (s *TerminalSrv) closeWsWrapper(ws *websocket.Conn) {
	s.closeWs(ws, "Session ended")
}

func (s *TerminalSrv) Startup(w http.ResponseWriter, r *http.Request) {
	hostIDStr := r.URL.Query().Get("hostId")
	if hostIDStr == "" {
		http.Error(w, "missing host id", http.StatusBadRequest)
		return
	}
	hostID, err := strconv.ParseUint(hostIDStr, 10, 32)
	if err != nil {
		http.Error(w, "invalid host id", http.StatusBadRequest)
		return
	}

	ws, err := ug.Upgrade(w, r, nil)
	if err != nil {
		s.Logger.Error("Failed to upgrade connection",
			zap.Error(err),
			zap.String("remote_addr", r.RemoteAddr),
		)
		return
	}

	if err = s.SSH(ws, uint(hostID)); err != nil {
		s.handleError(ws, err)
		s.closeWs(ws, s.formatError(err).Error)
		return
	}

	s.closeWs(ws, "会话已结束")
}

func (s *TerminalSrv) formatError(err error) *types.Message {
	switch {
	case errors.Is(err, websocket.ErrCloseSent):
		return &types.Message{
			Type:      types.MessageTypeError,
			Error:     "连接已关闭",
			ErrorCode: "CONNECTION_CLOSED",
			Details:   err.Error(),
		}
	case errors.Is(err, websocket.ErrReadLimit):
		return &types.Message{
			Type:      types.MessageTypeError,
			Error:     "连接数据超出限制",
			ErrorCode: "READ_LIMIT_EXCEEDED",
			Details:   err.Error(),
		}
	default:
		// 处理SSH相关错误
		errStr := err.Error()
		switch {
		case strings.Contains(errStr, "i/o timeout"):
			return &types.Message{
				Type:      types.MessageTypeError,
				Error:     "连接超时，请检查网络连接和主机状态",
				ErrorCode: "CONNECTION_TIMEOUT",
				Details:   errStr,
			}
		case strings.Contains(errStr, "connection refused"):
			return &types.Message{
				Type:      types.MessageTypeError,
				Error:     "连接被拒绝，请检查主机是否开启SSH服务",
				ErrorCode: "CONNECTION_REFUSED",
				Details:   errStr,
			}
		case strings.Contains(errStr, "no route"):
			return &types.Message{
				Type:      types.MessageTypeError,
				Error:     "无法连接到主机，请检查网络连接和主机地址",
				ErrorCode: "NO_ROUTE",
				Details:   errStr,
			}
		case strings.Contains(errStr, "auth"):
			return &types.Message{
				Type:      types.MessageTypeError,
				Error:     "认证失败，请检查用户名和密码",
				ErrorCode: "AUTH_FAILED",
				Details:   errStr,
			}
		default:
			return &types.Message{
				Type:      types.MessageTypeError,
				Error:     "连接失败，请稍后重试",
				ErrorCode: "UNKNOWN_ERROR",
				Details:   errStr,
			}
		}
	}
}

func (s *TerminalSrv) handleError(ws *websocket.Conn, err error) {
	if err != nil && ws != nil {
		msg := s.formatError(err)
		s.Logger.Error("Terminal connection error",
			zap.String("error_code", msg.ErrorCode),
			zap.String("error", msg.Error),
			zap.String("details", msg.Details),
		)
		if err = ws.WriteJSON(msg); err != nil {
			s.Logger.Error("Failed to write error message",
				zap.Error(err),
			)
		}
	}
}

func (s *TerminalSrv) SSH(ws *websocket.Conn, hostID uint) error {
	host, err := s.HostSrv.FindByID(hostID)
	if err != nil {
		return fmt.Errorf("failed to find host: %v", err)
	}
	sshConf := &adapter.SSHConfig{
		Host:     host.Host,
		Port:     host.Port,
		User:     host.Credential.Username,
		AuthType: host.Credential.AuthType,
	}

	switch host.Credential.AuthType {
	case enums.Password:
		sshConf.Password = host.Credential.Password
	case enums.PrivateKey:
		sshConf.PrivateKey = host.Credential.PrivateKey
		sshConf.KeyPassword = host.Credential.KeyPassword
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

	term := terminal.NewTerminal(ws, ssh, s.closeWsWrapper)
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

	term := terminal.NewTerminal(ws, serial, s.closeWsWrapper)
	term.Start()

	return nil
}
