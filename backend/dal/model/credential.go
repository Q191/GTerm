package model

import (
	"github.com/OpenToolkitLab/GTerm/backend/enums"
	"gorm.io/gorm"
)

type Credential struct {
	gorm.Model
	Name                   string `gorm:"uniqueIndex;not null"`
	Username               string
	Password               string
	PrivateKey             string
	KeyPassword            string
	Comment                string
	PasswordLoginPreferred bool
	IsCommonCredential     bool
}

func (c *Credential) TableName() string {
	return "credentials"
}

func (c *Credential) AuthType() enums.AuthType {
	if !c.PasswordLoginPreferred && c.PrivateKey != "" {
		return enums.PrivateKey
	}
	return enums.Password
}
