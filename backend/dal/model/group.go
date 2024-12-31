package model

import "gorm.io/gorm"

type Group struct {
	gorm.Model
	Name        string `json:"name" gorm:"uniqueIndex;not null"`
	Description string `json:"description"`
	Hosts       []*Host
}

func (c *Group) TableName() string {
	return "groups"
}
