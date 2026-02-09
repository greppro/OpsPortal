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
		&models.SystemConfig{},
		&models.Category{},
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

	// 初始化默认项目（4 个项目模拟企业多项目）
	var projCount int64
	DB.Model(&models.Project{}).Count(&projCount)
	if projCount == 0 {
		defaultProjects := []models.Project{
			{Name: "cloud", Label: "云平台", IsDefault: true},
			{Name: "monitor", Label: "监控平台"},
			{Name: "devops", Label: "DevOps工具"},
			{Name: "observability", Label: "可观测性"},
		}
		for _, proj := range defaultProjects {
			DB.Create(&proj)
		}
	}

	// 确保至少有一个项目为默认（兼容旧库或手动清空默认的情况）
	var defaultCount int64
	DB.Model(&models.Project{}).Where("is_default = ?", true).Count(&defaultCount)
	if defaultCount == 0 {
		var first models.Project
		if DB.First(&first).Error == nil {
			DB.Model(&first).Update("is_default", true)
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

	// 初始化分类（带图标），供工具归属与侧边栏展示
	seedCategories := []struct {
		name string
		desc string
		icon string
	}{
		{"可观测性", "链路、指标、日志等可观测性平台", "Monitor"},
		{"监控", "监控告警与可视化", "DataLine"},
		{"CI/CD", "持续集成与部署", "Operation"},
		{"云平台", "云控制台与资源管理", "ChromeFilled"},
		{"其它", "未归类的工具", "Grid"},
	}
	for _, c := range seedCategories {
		var n int64
		DB.Model(&models.Category{}).Where("name = ?", c.name).Count(&n)
		if n == 0 {
			DB.Create(&models.Category{Name: c.name, Description: c.desc, Icon: c.icon})
		}
	}

	// 初始化工具列表：常见运维开源平台与工具，模拟企业多项目多环境
	// 确保每个项目、每个环境、每个分类下都有卡片；部分工具仅单环境，部分多环境
	var toolCount int64
	DB.Model(&models.Tool{}).Count(&toolCount)
	if toolCount == 0 {
		projects := []string{"cloud", "monitor", "devops", "observability"}
		envs := []string{"dev", "test", "staging", "prod"}
		categories := []string{"可观测性", "监控", "CI/CD", "云平台", "其它"}

		// 开源工具池：(name, url, description)，按分类语义复用
		type toolDef struct {
			name string
			url  string
			desc string
		}
		byCat := map[string][]toolDef{
			"可观测性": {
				{"Grafana", "https://grafana.com", "度量分析与可视化，支持多种数据源"},
				{"Jaeger", "https://www.jaegertracing.io", "分布式链路追踪"},
				{"Zipkin", "https://zipkin.io", "分布式追踪系统"},
				{"Tempo", "https://grafana.com/oss/tempo/", "Grafana 分布式追踪后端"},
				{"Loki", "https://grafana.com/oss/loki/", "日志聚合系统"},
			},
			"监控": {
				{"Prometheus", "https://prometheus.io", "监控告警与时序数据库"},
				{"AlertManager", "https://prometheus.io/docs/alerting/latest/alertmanager/", "告警分组与路由"},
				{"Thanos", "https://thanos.io", "Prometheus 高可用与长期存储"},
				{"VictoriaMetrics", "https://victoriametrics.com", "高性能时序数据库"},
				{"Node Exporter", "https://github.com/prometheus/node_exporter", "主机指标采集"},
			},
			"CI/CD": {
				{"Jenkins", "https://www.jenkins.io", "持续集成与持续交付"},
				{"GitLab", "https://gitlab.com", "代码托管与 CI/CD 平台"},
				{"Argo CD", "https://argoproj.github.io/cd", "Kubernetes 声明式 GitOps 部署"},
				{"Tekton", "https://tekton.dev", "Kubernetes 原生 CI/CD"},
				{"Drone", "https://www.drone.io", "轻量级 CI 平台"},
			},
			"云平台": {
				{"阿里云控制台", "https://console.aliyun.com", "阿里云管理控制台"},
				{"腾讯云控制台", "https://console.cloud.tencent.com", "腾讯云管理控制台"},
				{"华为云控制台", "https://console.huaweicloud.com", "华为云管理控制台"},
				{"Kubernetes Dashboard", "https://kubernetes.io/docs/tasks/access-application-cluster/web-ui-dashboard/", "K8s 集群管理界面"},
				{"Rancher", "https://www.rancher.com", "多集群 K8s 管理"},
			},
			"其它": {
				{"Harbor", "https://goharbor.io", "企业级镜像仓库"},
				{"SonarQube", "https://www.sonarqube.org", "代码质量与安全分析"},
				{"Nexus", "https://www.sonatype.com/products/nexus-repository", "制品库管理"},
				{"Consul", "https://www.consul.io", "服务发现与配置"},
				{"Vault", "https://www.vaultproject.io", "密钥与敏感信息管理"},
			},
		}

		idx := 0
		for _, project := range projects {
			for _, env := range envs {
				for _, category := range categories {
					pool := byCat[category]
					if len(pool) == 0 {
						continue
					}
					tdef := pool[idx%len(pool)]
					idx++
					DB.Create(&models.Tool{
						Name:        tdef.name,
						URL:         tdef.url,
						Description: tdef.desc,
						Environment: env,
						Project:     project,
						Category:    category,
					})
				}
			}
		}
		// 再补充一批「多环境」同款工具：同一工具名在多个环境出现（部分 2/3/4 环境）
		multiEnvTools := []struct {
			name, url, desc, project, category string
			envs                               []string
		}{
			{"Grafana", "https://grafana.com", "度量分析与可视化", "monitor", "可观测性", []string{"dev", "test", "staging", "prod"}},
			{"Prometheus", "https://prometheus.io", "监控告警系统", "monitor", "监控", []string{"dev", "test", "prod"}},
			{"Jenkins", "https://www.jenkins.io", "CI/CD", "devops", "CI/CD", []string{"dev", "staging", "prod"}},
			{"Argo CD", "https://argoproj.github.io/cd", "GitOps 部署", "devops", "CI/CD", []string{"dev", "test", "staging", "prod"}},
			{"Harbor", "https://goharbor.io", "镜像仓库", "devops", "其它", []string{"test", "staging", "prod"}},
			{"Jaeger", "https://www.jaegertracing.io", "链路追踪", "observability", "可观测性", []string{"dev", "prod"}},
			{"阿里云控制台", "https://console.aliyun.com", "阿里云控制台", "cloud", "云平台", []string{"dev", "prod"}},
		}
		for _, m := range multiEnvTools {
			for _, e := range m.envs {
				DB.Create(&models.Tool{
					Name:        m.name,
					URL:         m.url,
					Description: m.desc,
					Environment: e,
					Project:     m.project,
					Category:    m.category,
				})
			}
		}
	}

	// 自动迁移 Notice 模型
	err = DB.AutoMigrate(&models.Notice{})
	if err != nil {
		log.Fatal("Failed to migrate Notice model:", err)
	}

}
