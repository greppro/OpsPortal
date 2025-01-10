package config

import (
	"log"
	"ops-portal/models"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDB() {
	var err error
	DB, err = gorm.Open(sqlite.Open("data/opsportal.db"), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// 自动迁移表结构
	DB.AutoMigrate(&models.User{})
	DB.AutoMigrate(&models.Tool{})

	// 修改 Environment 模型的迁移
	if err := DB.Migrator().DropTable(&models.Environment{}); err != nil {
		log.Printf("Failed to drop environments table: %v", err)
	}

	// 重新创建环境表
	if err := DB.AutoMigrate(&models.Environment{}); err != nil {
		log.Printf("Failed to migrate environments table: %v", err)
	}

	DB.AutoMigrate(&models.Project{})

	// 初始化默认用户
	var userCount int64
	DB.Model(&models.User{}).Count(&userCount)
	if userCount == 0 {
		DB.Create(&models.User{
			Username: "admin",
			Password: "admin123",
		})
	}

	// 初始化默认项目
	var projCount int64
	DB.Model(&models.Project{}).Count(&projCount)
	if projCount == 0 {
		defaultProjects := []models.Project{
			{Name: "default", Label: "默认项目"},
			{Name: "project-a", Label: "项目A"},
			{Name: "project-b", Label: "项目B"},
		}
		for _, proj := range defaultProjects {
			DB.Create(&proj)
		}
	}

	// 初始化默认环境
	var envCount int64
	DB.Model(&models.Environment{}).Count(&envCount)
	if envCount == 0 {
		// 获取默认项目
		var defaultProject models.Project
		DB.Where("name = ?", "default").First(&defaultProject)

		defaultEnvs := []models.Environment{
			{
				Name:      "dev",
				Label:     "开发环境",
				ProjectID: defaultProject.ID,
				Project:   defaultProject.Name,
			},
			{
				Name:      "prod",
				Label:     "生产环境",
				ProjectID: defaultProject.ID,
				Project:   defaultProject.Name,
			},
		}
		for _, env := range defaultEnvs {
			DB.Create(&env)
		}
	}

	// 自动迁移 Notice 模型
	err = DB.AutoMigrate(&models.Notice{})
	if err != nil {
		log.Fatal("Failed to migrate Notice model:", err)
	}
}
