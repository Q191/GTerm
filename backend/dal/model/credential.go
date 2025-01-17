package model

import (
	"github.com/MisakaTAT/GTerm/backend/enums"
)

type Credential struct {
	Common
	Name               string           `json:"name" gorm:"uniqueIndex;not null"`
	Username           string           `json:"username"`
	Password           string           `json:"password"`
	AuthMethod         enums.AuthMethod `json:"authMethod"`
	PrivateKey         string           `json:"privateKey"`
	Passphrase         string           `json:"passphrase"`
	IsCommonCredential bool             `json:"isCommonCredential"`
}

func (c *Credential) TableName() string {
	return "credentials"
}
