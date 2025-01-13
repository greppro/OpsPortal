package models

import (
	"time"

	"gorm.io/gorm"
)

type Notice struct {
	gorm.Model
	ID        uint      `json:"ID" gorm:"primarykey"`
	Content   string    `json:"content" gorm:"type:text;not null"`
	Active    bool      `json:"active" gorm:"default:true"`
	UpdatedAt time.Time `json:"UpdatedAt"`
}
