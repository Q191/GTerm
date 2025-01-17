package model

import "github.com/MisakaTAT/GTerm/backend/enums"

type Host struct {
	Common
	Name               string             `gorm:"uniqueIndex;not null" json:"name"`
	Host               string             `gorm:"not null" json:"host"`
	Port               uint               `gorm:"not null;default:22" json:"port"`
	ConnProtocol       enums.ConnProtocol `gorm:"not null" json:"connProtocol"`
	CredentialID       *uint              `gorm:"not null" json:"credentialID"`
	Credential         *Credential        `json:"credential"`
	IsCommonCredential bool               `gorm:"not null" json:"isCommonCredential"`
	Metadata           *Metadata          `json:"metadata"`
	GroupID            *uint              `json:"groupID"`
}

func (h *Host) TableName() string {
	return "hosts"
}
