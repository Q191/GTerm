package model

type Group struct {
	Common
	Name        string `json:"name" gorm:"uniqueIndex;not null"`
	Description string `json:"description" gorm:"not null"`
}

func (g *Group) TableName() string {
	return "groups"
}
