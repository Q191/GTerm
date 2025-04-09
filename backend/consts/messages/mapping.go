package messages

var CodeMapping = map[string]string{
	CreateSuccess: "Creation successful",
	UpdateSuccess: "Update successful",
	DeleteSuccess: "Deletion successful",

	ConnectionClosed:           "Connection closed",
	ReadLimitExceeded:          "Connection data exceeded limit",
	ConnectionTimeout:          "Connection timeout",
	ConnectionRefused:          "Connection refused",
	NoRoute:                    "No route to host",
	AuthFailed:                 "Authentication failed",
	UnknownError:               "Connection failed",
	InvalidCredentials:         "Invalid credentials",
	PermissionDenied:           "Permission denied",
	HostUnreachable:            "Host unreachable",
	SSHServiceDown:             "SSH service down",
	NetworkError:               "Network error",
	ProtocolError:              "Protocol error",
	ResourceExhausted:          "Server resources exhausted",
	SessionEnded:               "Session ended",
	FailedToSendFingerprintMsg: "Failed to send fingerprint message",
	FailedToReadFingerprint:    "Failed to read fingerprint confirmation",
	FailedToParseFingerprint:   "Failed to parse fingerprint confirmation",
	FailedToAddFingerprint:     "Failed to add host fingerprint",
	UserRejectedFingerprint:    "User rejected host fingerprint",
}
