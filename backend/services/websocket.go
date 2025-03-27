package services

import (
	"errors"
	"github.com/MisakaTAT/GTerm/backend/types"
	"github.com/google/wire"
	"github.com/gorilla/websocket"
	"go.uber.org/zap"
	"net/http"
	"strconv"
	"strings"
)

var WebsocketSrvSet = wire.NewSet(wire.Struct(new(WebsocketSrv), "*"))

type WebsocketSrv struct {
	TerminalSrv *TerminalSrv
	Logger      *zap.Logger
}

var ug = websocket.Upgrader{
	ReadBufferSize:    1024,
	WriteBufferSize:   1024 * 1024 * 10,
	EnableCompression: true,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func (s *WebsocketSrv) formatError(err error) *types.Message {
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

func (s *WebsocketSrv) handleError(ws *websocket.Conn, err error) {
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

func (s *WebsocketSrv) TerminalHandle(w http.ResponseWriter, r *http.Request) {
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

	if err = s.TerminalSrv.SSH(ws, uint(hostID)); err != nil {
		s.handleError(ws, err)
		s.TerminalSrv.CloseSession(ws, s.formatError(err).Error)
		return
	}

	s.TerminalSrv.CloseSession(ws, "会话已结束")
}
