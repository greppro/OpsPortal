package models

type User struct {
	ID       uint   `gorm:"primarykey"`
	Username string `gorm:"unique"`
	Password string
} 