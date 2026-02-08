package models

// SystemConfig 系统配置键值对
type SystemConfig struct {
	Key   string `json:"key" gorm:"primaryKey;size:64;not null"`
	Value string `json:"value" gorm:"type:text"`
}
