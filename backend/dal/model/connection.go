package model

import "github.com/MisakaTAT/GTerm/backend/enums"

type Connection struct {
	Common
	Label              string             `gorm:"uniqueIndex;not null" json:"label"`
	Host               string             `json:"host"`
	Port               uint               `json:"port"`
	SerialPort         string             `json:"serialPort"`
	ConnProtocol       enums.ConnProtocol `gorm:"not null" json:"connProtocol"`
	CredentialID       *uint              `gorm:"not null" json:"credentialID"`
	Credential         *Credential        `json:"credential"`
	IsCommonCredential bool               `gorm:"not null" json:"isCommonCredential"`
	Metadata           *Metadata          `json:"metadata"`
	GroupID            *uint              `json:"groupID"`
}

func (c *Connection) TableName() string {
	return "connections"
}
