package model

import "gorm.io/gorm"

type Connection struct {
	gorm.Model
	Label             string `gorm:"uniqueIndex;not null"`
	Address           string `gorm:"not null"`
	Port              uint32 `gorm:"not null"`
	Comment           string `gorm:"not null"`
	CredentialID      uint   `gorm:"not null"`
	Credential        *Credential
	ConnectionGroupID *uint
}

func (c *Connection) TableName() string {
	return "connections"
}
