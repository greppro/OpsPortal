package models

import "time"

type Tool struct {
	ID          uint   `json:"id" gorm:"primaryKey"`
	Name        string `json:"name" gorm:"not null"`
	Description string `json:"description"`
	URL         string `json:"url" gorm:"not null"`
	Environment string `json:"environment" gorm:"not null"`
	Project     string `json:"project" gorm:"not null"`
	Category    string `json:"category"`
	Icon        string `json:"icon"`
	// ProbeDisabled 为 true 时不做可用性检测，也不出现在 Prometheus 服务发现里（默认检测）
	ProbeDisabled bool      `json:"probe_disabled" gorm:"not null;default:false"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}
