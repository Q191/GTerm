package model

import (
	"github.com/OpenToolkitLab/GTerm/backend/enums"
	"gorm.io/gorm"
)

type Credential struct {
	gorm.Model
	Name                   string `gorm:"uniqueIndex;not null"`
	Username               string `gorm:"not null"`
	Password               string `gorm:"not null"`
	PrivateKey             string `gorm:"not null"`
	KeyPassword            string `gorm:"not null"`
	Comment                string `gorm:"not null"`
	PasswordLoginPreferred bool   `gorm:"not null"`
	IsCommonCredential     bool   `gorm:"not null"`
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
