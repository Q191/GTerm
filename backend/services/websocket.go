package services

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"

	"github.com/MisakaTAT/GTerm/backend/consts/messages"
	"github.com/MisakaTAT/GTerm/backend/enums"
	"github.com/MisakaTAT/GTerm/backend/initialize"
	"github.com/MisakaTAT/GTerm/backend/types"
	"github.com/google/wire"
	"github.com/gorilla/websocket"
)

var WebsocketSrvSet = wire.NewSet(wire.Struct(new(WebsocketSrv), "*"))

type WebsocketSrv struct {
	TerminalSrv *TerminalSrv
	Logger      initialize.Logger
}

var ug = websocket.Upgrader{
	ReadBufferSize:    1024,
	WriteBufferSize:   1024 * 1024 * 10,
	EnableCompression: true,
	CheckOrigin: func(r *http.Request) bool {
		origin := r.Header.Get("Origin")
		if origin == "" {
			return false
		}
		allowedOrigins := []string{
			"wails://wails",
		}
		for _, allowed := range allowedOrigins {
			if strings.HasPrefix(origin, allowed) {
				return true
			}
		}
		return false
	},
}

var errorMappings = map[string]string{
	"i/o timeout":          messages.ConnectionTimeout,
	"connection refused":   messages.ConnectionRefused,
	"no route":             messages.NoRoute,
	"auth":                 messages.AuthFailed,
	"invalid credentials":  messages.InvalidCredentials,
	"incorrect password":   messages.InvalidCredentials,
	"wrong username":       messages.InvalidCredentials,
	"permission denied":    messages.PermissionDenied,
	"access denied":        messages.PermissionDenied,
	"host unreachable":     messages.HostUnreachable,
	"no such host":         messages.HostUnreachable,
	"ssh service":          messages.SSHServiceDown,
	"port 22":              messages.SSHServiceDown,
	"network":              messages.NetworkError,
	"connection reset":     messages.NetworkError,
	"broken pipe":          messages.NetworkError,
	"protocol":             messages.ProtocolError,
	"incompatible":         messages.ProtocolError,
	"version":              messages.ProtocolError,
	"resource":             messages.ResourceExhausted,
	"too many connections": messages.ResourceExhausted,
	"connection limit":     messages.ResourceExhausted,
}

func (s *WebsocketSrv) formatError(err error) *types.Message {
	switch {
	case errors.Is(err, websocket.ErrCloseSent):
		return &types.Message{
			Type:    enums.TerminalTypeError,
			Message: messages.CodeMapping[messages.ConnectionClosed],
			Code:    messages.ConnectionClosed,
			Details: err.Error(),
		}
	case errors.Is(err, websocket.ErrReadLimit):
		return &types.Message{
			Type:    enums.TerminalTypeError,
			Message: messages.CodeMapping[messages.ReadLimitExceeded],
			Code:    messages.ReadLimitExceeded,
			Details: err.Error(),
		}
	default:
		errStr := err.Error()
		for pattern, code := range errorMappings {
			if strings.Contains(errStr, pattern) {
				return &types.Message{
					Type:    enums.TerminalTypeError,
					Message: messages.CodeMapping[code],
					Code:    code,
					Details: errStr,
				}
			}
		}

		s.Logger.Error("Unhandled error type: %s", errStr)
		return &types.Message{
			Type:    enums.TerminalTypeError,
			Message: messages.CodeMapping[messages.UnknownError],
			Code:    messages.UnknownError,
			Details: errStr,
		}
	}
}

func (s *WebsocketSrv) handleError(ws *websocket.Conn, err error) {
	if err != nil && ws != nil {
		e := s.formatError(err)
		s.Logger.Error("Terminal connection error: %s, code: %s, details: %s",
			e.Message,
			e.Code,
			e.Details,
		)
		if err = ws.WriteJSON(e); err != nil {
			s.Logger.Error("Failed to write error message: %v", err)
		}
	}
}

func (s *WebsocketSrv) TerminalHandle(w http.ResponseWriter, r *http.Request) {
	s.Logger.Info("Received terminal connection request: %s", r.RemoteAddr)
	hostIDStr := r.URL.Query().Get("hostId")
	if hostIDStr == "" {
		s.Logger.Error("Request missing hostId parameter, remote_addr: %s", r.RemoteAddr)
		http.Error(w, "missing host id", http.StatusBadRequest)
		return
	}
	hostID, err := strconv.ParseUint(hostIDStr, 10, 32)
	if err != nil {
		s.Logger.Error("Invalid hostId parameter: %s, remote_addr: %s", hostIDStr, r.RemoteAddr)
		http.Error(w, "invalid host id", http.StatusBadRequest)
		return
	}

	s.Logger.Info("Upgrading WebSocket connection, hostId: %d, remote_addr: %s", hostID, r.RemoteAddr)
	ws, err := ug.Upgrade(w, r, nil)
	if err != nil {
		s.Logger.Error("Failed to upgrade WebSocket connection: %v, remote_addr: %s",
			err,
			r.RemoteAddr,
		)
		return
	}
	s.Logger.Info("WebSocket connection upgraded successfully, hostId: %d, remote_addr: %s", hostID, r.RemoteAddr)

	// 建立SSH连接
	s.Logger.Info("Starting SSH connection, hostId: %d", hostID)
	err = s.TerminalSrv.SSH(ws, uint(hostID))

	var fingerprintErr *types.FingerprintError
	if errors.As(err, &fingerprintErr) {
		// 发送主机指纹确认消息
		s.Logger.Info("Host fingerprint confirmation needed, host: %s, fingerprint: %s",
			fingerprintErr.Host,
			fingerprintErr.Fingerprint,
		)
		if err = ws.WriteJSON(&types.Message{
			Type:        enums.TerminalTypeFingerprintConfirm,
			Host:        fingerprintErr.Host,
			Fingerprint: fingerprintErr.Fingerprint,
		}); err != nil {
			s.Logger.Error("Failed to send host fingerprint confirmation message: %v", err)
			s.TerminalSrv.CloseSession(ws, messages.FailedToSendFingerprintMsg)
			return
		}
		s.Logger.Info("Host fingerprint confirmation message sent, waiting for client confirmation")

		// 等待客户端确认
		_, data, err := ws.ReadMessage()
		if err != nil {
			s.Logger.Error("Failed to read host fingerprint confirmation response: %v", err)
			s.TerminalSrv.CloseSession(ws, messages.FailedToReadFingerprint)
			return
		}
		s.Logger.Info("Received client fingerprint confirmation response")

		var fg types.Fingerprint
		if err = json.Unmarshal(data, &fg); err != nil {
			s.Logger.Error("Failed to parse host fingerprint confirmation response: %v, data: %s", err, string(data))
			s.TerminalSrv.CloseSession(ws, messages.FailedToParseFingerprint)
			return
		}

		if fg.Type == enums.TerminalTypeFingerprintConfirm && fg.Accept {
			s.Logger.Info("Client accepted host fingerprint, adding to known_hosts")
			// 添加主机指纹并重新连接
			if err = s.TerminalSrv.AddFingerprint(uint(hostID), fingerprintErr.Host, fingerprintErr.Fingerprint); err != nil {
				s.Logger.Error("Failed to add host fingerprint: %v", err)
				s.handleError(ws, err)
				s.TerminalSrv.CloseSession(ws, messages.FailedToAddFingerprint)
				return
			}
			s.Logger.Info("Host fingerprint added, retrying connection")

			// 重新尝试连接
			if err = s.TerminalSrv.SSH(ws, uint(hostID)); err != nil {
				s.Logger.Error("Failed to reconnect SSH: %v", err)
				s.handleError(ws, err)
				s.TerminalSrv.CloseSession(ws, s.formatError(err).Message)
				return
			}
			s.Logger.Info("SSH reconnection successful")
		} else {
			// 用户拒绝添加主机指纹
			s.Logger.Info("Client rejected host fingerprint, closing connection")
			s.TerminalSrv.CloseSession(ws, messages.UserRejectedFingerprint)
			return
		}
	} else if err != nil {
		s.Logger.Error("SSH connection failed: %v", err)
		s.handleError(ws, err)
		s.TerminalSrv.CloseSession(ws, s.formatError(err).Message)
		return
	}

	s.Logger.Info("Terminal session ended, closing WebSocket connection")
	s.TerminalSrv.CloseSession(ws, messages.SessionEnded)
}
