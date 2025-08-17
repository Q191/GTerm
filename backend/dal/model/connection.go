package model

import (
	"github.com/Q191/GTerm/backend/enums"
	"go.bug.st/serial"
)

type Connection struct {
	Common
	Label                  string             `gorm:"uniqueIndex;not null" json:"label"`
	Host                   string             `json:"host"`
	Port                   uint               `json:"port"`
	SerialPort             string             `json:"serialPort"`
	ConnProtocol           enums.ConnProtocol `gorm:"not null" json:"connProtocol"`
	CredentialID           *uint              `gorm:"not null" json:"credentialID"`
	Credential             *Credential        `json:"credential"`
	UseCommonCredential    bool               `gorm:"not null" json:"useCommonCredential"`
	Metadata               *Metadata          `json:"metadata"`
	GroupID                *uint              `json:"groupID"`
	BaudRate               int                `json:"baudRate"`
	DataBits               int                `json:"dataBits"`
	StopBits               serial.StopBits    `json:"stopBits"`
	Parity                 serial.Parity      `json:"parity"`
	Theme                  string             `json:"theme" gorm:"not null,default:'Default'"`
	SSHKeyExchanges        []string           `json:"sshKeyExchanges" gorm:"type:json;serializer:json"`
	SSHCiphers             []string           `json:"sshCiphers" gorm:"type:json;serializer:json"`
	SSHMACs                []string           `json:"sshMACs" gorm:"type:json;serializer:json"`
	SSHPublicKeyAlgorithms []string           `json:"sshPublicKeyAlgorithms" gorm:"type:json;serializer:json"`
	SSHHostKeyAlgorithms   []string           `json:"sshHostKeyAlgorithms" gorm:"type:json;serializer:json"`
	SSHCharset             string             `json:"sshCharset" gorm:"default:'UTF-8'"`
}

func (c *Connection) TableName() string {
	return "connections"
}
