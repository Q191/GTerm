package enums

type CredentialAuthType = int

const (
	CredentialAuthTypePassword CredentialAuthType = iota
	CredentialAuthTypePrivateKey
)

type CredentialType = int

const (
	CredentialTypePassword CredentialAuthType = iota
	CredentialTypePrivateKey
	CredentialTypeCommon
)
