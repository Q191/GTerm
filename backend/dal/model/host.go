package model

type Host struct {
	Common
	Name         string      `gorm:"uniqueIndex;not null" json:"name"`
	Host         string      `gorm:"not null" json:"host"`
	Port         uint        `gorm:"not null;default:22" json:"port"`
	Description  string      `gorm:"not null" json:"description"`
	CredentialID uint        `gorm:"not null" json:"credential_id"`
	Credential   *Credential `json:"credential"`
	Metadata     *Metadata   `json:"metadata"`
	GroupID      *uint       `json:"group_id"`
}

func (h *Host) TableName() string {
	return "hosts"
}
