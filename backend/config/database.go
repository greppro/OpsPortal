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

	// 自动迁移数据库表
	err = DB.AutoMigrate(
		&models.Tool{},
		&models.Project{},
		&models.Environment{},
		&models.Notice{},
		&models.Logo{},
	)
	if err != nil {
		log.Fatal("Failed to migrate database:", err)
	}

	// 自动迁移用户表
	DB.AutoMigrate(&models.User{})

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
			{Name: "cloud", Label: "云平台"},
			{Name: "monitor", Label: "监控平台"},
			{Name: "devops", Label: "DevOps工具"},
		}
		for _, proj := range defaultProjects {
			DB.Create(&proj)
		}
	}

	// 初始化默认环境
	var envCount int64
	DB.Model(&models.Environment{}).Count(&envCount)
	if envCount == 0 {
		// 为每个项目创建环境
		var projects []models.Project
		DB.Find(&projects)

		environments := []struct {
			Name  string
			Label string
		}{
			{"dev", "开发环境"},
			{"test", "测试环境"},
			{"staging", "预发环境"},
			{"prod", "生产环境"},
		}

		for _, project := range projects {
			for _, env := range environments {
				DB.Create(&models.Environment{
					Name:      env.Name,
					Label:     env.Label,
					ProjectID: project.ID,
					Project:   project.Name,
				})
			}
		}
	}

	// 初始化工具列表
	var toolCount int64
	DB.Model(&models.Tool{}).Count(&toolCount)
	if toolCount == 0 {
		// 获取项目ID
		var cloudProject, monitorProject, devopsProject models.Project
		DB.Where("name = ?", "cloud").First(&cloudProject)
		DB.Where("name = ?", "monitor").First(&monitorProject)
		DB.Where("name = ?", "devops").First(&devopsProject)

		tools := []models.Tool{
			// 云平台工具
			{
				Name:        "阿里云控制台",
				URL:         "https://console.aliyun.com",
				Description: "阿里云管理控制台，提供云服务器、数据库、存储等云服务",
				Environment: "prod",
				Project:    "cloud",
			},
			{
				Name:        "腾讯云控制台",
				URL:         "https://console.cloud.tencent.com",
				Description: "腾讯云管理控制台，提供云服务器、数据库、CDN等云服务",
				Environment: "prod",
				Project:    "cloud",
			},
			{
				Name:        "华为云控制台",
				URL:         "https://console.huaweicloud.com",
				Description: "华为云管理控制台，提供云计算、人工智能、大数据等服务",
				Environment: "prod",
				Project:    "cloud",
			},
			// 监控平台工具
			{
				Name:        "Grafana",
				URL:         "https://grafana.com",
				Description: "开源的度量分析和可视化工具，支持多种数据源",
				Environment: "prod",
				Project:    "monitor",
			},
			{
				Name:        "Prometheus",
				URL:         "https://prometheus.io",
				Description: "开源的监控告警系统，提供强大的时序数据存储和查询功能",
				Environment: "prod",
				Project:    "monitor",
			},
			{
				Name:        "AlertManager",
				URL:         "https://prometheus.io/docs/alerting/latest/alertmanager/",
				Description: "Prometheus告警管理器，处理告警分发和分组",
				Environment: "prod",
				Project:    "monitor",
			},
			// DevOps工具
			{
				Name:        "Jenkins",
				URL:         "https://www.jenkins.io",
				Description: "开源的持续集成和持续交付工具，自动化构建、测试和部署",
				Environment: "prod",
				Project:    "devops",
			},
			{
				Name:        "ArgoCD",
				URL:         "https://argoproj.github.io/cd",
				Description: "Kubernetes的声明式持续部署工具",
				Environment: "prod",
				Project:    "devops",
			},
			{
				Name:        "Harbor",
				URL:         "https://goharbor.io",
				Description: "企业级容器镜像仓库，提供镜像存储、签名、扫描等功能",
				Environment: "prod",
				Project:    "devops",
			},
			{
				Name:        "GitLab",
				URL:         "https://gitlab.com",
				Description: "完整的 DevOps 平台，提供代码托管、CI/CD、制品库等功能",
				Environment: "prod",
				Project:    "devops",
			},
			{
				Name:        "SonarQube",
				URL:         "https://www.sonarqube.org",
				Description: "代码质量管理平台，提供代码检查、漏洞分析等功能",
				Environment: "prod",
				Project:    "devops",
			},
		}

		for _, tool := range tools {
			DB.Create(&tool)
		}
	}

	// 自动迁移 Notice 模型
	err = DB.AutoMigrate(&models.Notice{})
	if err != nil {
		log.Fatal("Failed to migrate Notice model:", err)
	}
}
