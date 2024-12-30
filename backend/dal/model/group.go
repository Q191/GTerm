package model

import "gorm.io/gorm"

type Group struct {
	gorm.Model
	Name  string `gorm:"uniqueIndex;not null"`
	Hosts []*Host
}

func (c *Group) TableName() string {
	return "groups"
}
