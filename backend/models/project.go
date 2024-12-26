package models

type Project struct {
	ID           uint          `json:"id" gorm:"primarykey"`
	Name         string        `json:"name"`
	Label        string        `json:"label"`
	Environments []Environment `json:"environments" gorm:"foreignKey:ProjectID"`
}
