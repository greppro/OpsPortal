package models

type Project struct {
	ID           uint          `json:"id" gorm:"primarykey"`
	Name         string        `json:"name"`
	Label        string        `json:"label"`
	IsDefault    bool          `json:"is_default" gorm:"default:false"`
	Environments []Environment `json:"environments" gorm:"foreignKey:ProjectID"`
}
