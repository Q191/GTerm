package model

type Metadata struct {
	Common
	ConnectionID uint   `json:"connectionId" gorm:"not null"`
	Vendor       string `json:"vendor" gorm:"not null"`
	Type         string `json:"type" gorm:"not null"`
}

func (m *Metadata) TableName() string {
	return "metadata"
}
