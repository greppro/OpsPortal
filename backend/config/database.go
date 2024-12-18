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
	err = db.AutoMigrate(&models.Tool{}, &models.User{})
	if err != nil {
		panic(fmt.Sprintf("failed to auto migrate: %v", err))
	}

	// 插入测试数据
	var count int64
	// 检查是否存在默认用户
	db.Model(&models.User{}).Count(&count)
	if count == 0 {
		db.Create(&models.User{
			Username: "admin",
			Password: "admin123",
		})
	}

	// 检查是否存在工具数据
	db.Model(&models.Tool{}).Count(&count)
	if count == 0 {
		tools := []models.Tool{
			{
				Name:        "Grafana Dev",
				Description: "Grafana 开发环境监控面板",
				URL:         "http://grafana-dev.example.com",
				Environment: "dev",
				Category:    "监控",
			},
			{
				Name:        "Prometheus Dev",
				Description: "Prometheus 开发环境监控系统",
				URL:         "http://prometheus-dev.example.com",
				Environment: "dev",
				Category:    "监控",
			},
			{
				Name:        "Jenkins Dev",
				Description: "Jenkins 开发环境持续集成平台",
				URL:         "http://jenkins-dev.example.com",
				Environment: "dev",
				Category:    "部署",
			},
			{
				Name:        "GitLab Dev",
				Description: "GitLab 开发环境代码仓库",
				URL:         "http://gitlab-dev.example.com",
				Environment: "dev",
				Category:    "代码",
			},
		}
		db.Create(&tools)
	}

	DB = db
} 