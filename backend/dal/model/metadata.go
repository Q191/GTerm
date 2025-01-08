package model

type Metadata struct {
	Common
	HostID     uint `json:"host_id" gorm:"not null"`
	Processors uint `json:"processors" gorm:"not null"`
	MemTotal   uint `json:"mem_total" gorm:"not null"`
}

func (m *Metadata) TableName() string {
	return "metadata"
}
