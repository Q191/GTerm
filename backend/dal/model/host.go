package model

import "gorm.io/gorm"

type Host struct {
	gorm.Model
	Label        string `gorm:"uniqueIndex;not null"`
	Address      string
	Port         uint32
	Comment      string
	CredentialID uint
	Credential   *Credential
	GroupID      *uint
}

func (c *Host) TableName() string {
	return "hosts"
}
