package model

import "gorm.io/gorm"

type ConnectionGroup struct {
	gorm.Model
	Label       string `gorm:"uniqueIndex;not null"`
	Connections []*Connection
}

func (c *ConnectionGroup) TableName() string {
	return "connection_groups"
}
