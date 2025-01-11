package models

type Logo struct {
	ID   uint   `json:"id" gorm:"primarykey"`
	URL  string `json:"url" gorm:"not null"`
	Name string `json:"name"`
} 