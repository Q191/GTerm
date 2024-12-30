package model

import "gorm.io/gorm"

type Host struct {
	gorm.Model
	Label        string `gorm:"uniqueIndex;not null"`
	Address      string `gorm:"not null"`
	Port         uint32 `gorm:"not null"`
	Comment      string `gorm:"not null"`
	CredentialID uint   `gorm:"not null"`
	Credential   *Credential
	GroupID      *uint
}

func (c *Host) TableName() string {
	return "hosts"
}
