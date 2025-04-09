package messages

const (
	ConnectionClosed           = "websocket.error.connection_closed"
	ReadLimitExceeded          = "websocket.error.read_limit_exceeded"
	ConnectionTimeout          = "websocket.error.connection_timeout"
	ConnectionRefused          = "websocket.error.connection_refused"
	NoRoute                    = "websocket.error.no_route"
	AuthFailed                 = "websocket.error.auth_failed"
	UnknownError               = "websocket.error.unknown_error"
	InvalidCredentials         = "websocket.error.invalid_credentials"
	PermissionDenied           = "websocket.error.permission_denied"
	HostUnreachable            = "websocket.error.host_unreachable"
	SSHServiceDown             = "websocket.error.ssh_service_down"
	NetworkError               = "websocket.error.network_error"
	ProtocolError              = "websocket.error.protocol_error"
	ResourceExhausted          = "websocket.error.resource_exhausted"
	SessionEnded               = "websocket.info.session_ended"
	FailedToSendFingerprintMsg = "websocket.error.failed_to_send_fingerprint_msg"
	FailedToReadFingerprint    = "websocket.error.failed_to_read_fingerprint"
	FailedToParseFingerprint   = "websocket.error.failed_to_parse_fingerprint"
	FailedToAddFingerprint     = "websocket.error.failed_to_add_fingerprint"
	UserRejectedFingerprint    = "websocket.info.user_rejected_fingerprint"
)
