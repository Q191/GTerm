package enums

type AuthType = int

const (
	Password AuthType = iota
	PrivateKey
)
