package config

import (
	"fmt"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"ops-portal/models"
)

var DB *gorm.DB

func InitDB() {
	// 数据库文件会保存在 data/ops_portal.db
	db, err := gorm.Open(sqlite.Open("data/ops_portal.db"), &gorm.Config{})
	if err != nil {
		panic(fmt.Sprintf("failed to connect database: %v", err))
	}

	// 自动迁移
	err = db.AutoMigrate(&models.Tool{})
	if err != nil {
		panic(fmt.Sprintf("failed to auto migrate: %v", err))
	}

	// 插入一些测试数据
	var count int64
	db.Model(&models.Tool{}).Count(&count)
	if count == 0 {
		tools := []models.Tool{
			{
				Name:        "Grafana Dev",
				Description: "Grafana 开发环境监控面板",
				URL:         "http://grafana-dev.example.com",
				Environment: "dev",
				Category:    "监控",
				Icon:        "Monitor",
			},
			{
				Name:        "Prometheus Dev",
				Description: "Prometheus 开发环境监控系统",
				URL:         "http://prometheus-dev.example.com",
				Environment: "dev",
				Category:    "监控",
				Icon:        "DataLine",
			},
			{
				Name:        "Grafana Prod",
				Description: "Grafana 生产环境监控面板",
				URL:         "http://grafana-prod.example.com",
				Environment: "prod",
				Category:    "监控",
				Icon:        "Monitor",
			},
			{
				Name:        "Prometheus Prod",
				Description: "Prometheus 生产环境监控系统",
				URL:         "http://prometheus-prod.example.com",
				Environment: "prod",
				Category:    "监控",
				Icon:        "DataLine",
			},
		}
		db.Create(&tools)
	}

	DB = db
} 