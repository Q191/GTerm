package enums

import "strings"

type AuthMethod string

const (
	Password   AuthMethod = "Password"
	PrivateKey AuthMethod = "PrivateKey"
)

var AuthMethodEnums = []AuthMethod{Password, PrivateKey}

func (a AuthMethod) TSName() string {
	return strings.ToUpper(string(a))
}
