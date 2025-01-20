package types

type MessageType string

const (
	MessageTypeError     MessageType = "error"
	MessageTypeData      MessageType = "data"
	MessageTypeConnected MessageType = "connected"
)

type Message struct {
	Type      MessageType `json:"type"`
	Content   interface{} `json:"content,omitempty"`
	Error     string      `json:"error,omitempty"`
	Details   string      `json:"details,omitempty"`
	ErrorCode string      `json:"error_code,omitempty"`
}
