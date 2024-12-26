package models

type Environment struct {
	ID        uint   `json:"id" gorm:"primarykey"`
	Name      string `json:"name" gorm:"not null;uniqueIndex:idx_project_env,priority:2"`
	Label     string `json:"label" gorm:"not null"`
	ProjectID uint   `json:"project_id" gorm:"not null;uniqueIndex:idx_project_env,priority:1"`
	Project   string `json:"project" gorm:"not null"`
}
