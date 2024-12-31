package model

type Group struct {
	Common
	Name        string `json:"name" gorm:"uniqueIndex;not null"`
	Description string `json:"description" gorm:"not null"`
	Hosts       []*Host
}

func (c *Group) TableName() string {
	return "groups"
}
