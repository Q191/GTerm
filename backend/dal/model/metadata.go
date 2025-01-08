package model

type Metadata struct {
	Common
	HostID uint   `json:"host_id" gorm:"not null"`
	OS     string `json:"os" gorm:"not null"`
}

func (m *Metadata) TableName() string {
	return "metadata"
}
