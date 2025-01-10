package models

import "gorm.io/gorm"

type Notice struct {
	gorm.Model
	Content string `json:"content" gorm:"type:text;not null"`
	Active  bool   `json:"active" gorm:"default:true"`
}
