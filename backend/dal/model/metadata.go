package model

type Metadata struct {
	Common
	ConnectionID uint   `json:"connectionId" gorm:"not null"`
	OS           string `json:"os" gorm:"not null"`
}

func (m *Metadata) TableName() string {
	return "metadata"
}
