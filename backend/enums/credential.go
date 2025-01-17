package enums

type AuthMethod = int

const (
	Password AuthMethod = iota
	PrivateKey
)
