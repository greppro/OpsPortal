package config

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"log"
	"ops-portal/models"
	"ops-portal/utils"
	"os"
	"path/filepath"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB

const (
	defaultAdminUsername = "admin"
	keyJWTSecret         = "jwt_secret"
)

func InitDB(s Settings) {
	var err error
	if dir := filepath.Dir(s.DBPath); dir != "" {
		if err := os.MkdirAll(dir, 0755); err != nil {
			log.Fatalf("Failed to create database directory %s: %v", dir, err)
		}
	}

	// busy_timeout：后台探测在读库时，写操作等待锁而不是直接报 database is locked
	DB, err = gorm.Open(sqlite.Open(s.DBPath+"?_busy_timeout=5000"), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// 自动迁移数据库表（只增加缺少的表和列，不删除已有数据）
	err = DB.AutoMigrate(
		&models.Tool{},
		&models.Project{},
		&models.Environment{},
		&models.Notice{},
		&models.Logo{},
		&models.SystemConfig{},
		&models.Category{},
		&models.User{},
	)
	if err != nil {
		log.Fatal("Failed to migrate database:", err)
	}

	if err := initAdminUser(s.AdminPassword); err != nil {
		log.Fatal("Failed to initialize admin user:", err)
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

	var titleConfig models.SystemConfig
	if err := DB.Where("key = ?", "site_title").First(&titleConfig).Error; err != nil {
		DB.Create(&models.SystemConfig{Key: "site_title", Value: "OpsPortal运维导航"})
	} else if titleConfig.Value == "" {
		DB.Model(&titleConfig).Update("value", "OpsPortal运维导航")
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
}

// initAdminUser 首次启动时创建管理员；已有用户时把旧版本留下的明文密码转成 bcrypt 哈希。
func initAdminUser(adminPassword string) error {
	var users []models.User
	if err := DB.Find(&users).Error; err != nil {
		return err
	}

	if len(users) == 0 {
		password := adminPassword
		generated := password == ""
		if generated {
			var err error
			if password, err = utils.RandomString(16); err != nil {
				return err
			}
		}
		if len(password) > utils.MaxPasswordBytes {
			return fmt.Errorf("ADMIN_PASSWORD 不能超过 %d 字节", utils.MaxPasswordBytes)
		}
		hash, err := utils.HashPassword(password)
		if err != nil {
			return err
		}
		if err := DB.Create(&models.User{Username: defaultAdminUsername, Password: hash}).Error; err != nil {
			return err
		}
		if generated {
			log.Printf("==================================================================")
			log.Printf("已创建管理员账号 %s，初始密码：%s", defaultAdminUsername, password)
			log.Printf("该密码只在首次启动时显示一次，请登录后台后立即修改。")
			log.Printf("也可以在首次启动前通过环境变量 ADMIN_PASSWORD 指定初始密码。")
			log.Printf("==================================================================")
		} else {
			log.Printf("已使用 ADMIN_PASSWORD 创建管理员账号 %s", defaultAdminUsername)
		}
		return nil
	}

	for _, u := range users {
		if utils.IsPasswordHash(u.Password) {
			continue
		}
		plain := u.Password
		if len(plain) > utils.MaxPasswordBytes {
			plain = plain[:utils.MaxPasswordBytes]
		}
		hash, err := utils.HashPassword(plain)
		if err != nil {
			return err
		}
		if err := DB.Model(&models.User{}).Where("id = ?", u.ID).Update("password", hash).Error; err != nil {
			return err
		}
		log.Printf("用户 %s 的密码已从明文迁移为 bcrypt 哈希", u.Username)
		if plain == utils.LegacyDefaultPassword {
			log.Printf("安全提醒：用户 %s 仍在使用默认密码，请登录后台后立即修改", u.Username)
		}
	}
	return nil
}

// InitJWTSecret 返回 JWT 签名密钥：优先用 JWT_SECRET；未设置时读取数据库里保存的密钥，
// 没有则随机生成一个并保存，保证重启后已登录的 token 仍然有效。
func InitJWTSecret(s Settings) ([]byte, error) {
	if s.JWTSecret != "" {
		return []byte(s.JWTSecret), nil
	}

	var rows []models.SystemConfig
	if err := DB.Where("key = ?", keyJWTSecret).Limit(1).Find(&rows).Error; err != nil {
		return nil, err
	}
	if len(rows) == 1 && len(rows[0].Value) >= 32 {
		return []byte(rows[0].Value), nil
	}

	buf := make([]byte, 32)
	if _, err := rand.Read(buf); err != nil {
		return nil, err
	}
	secret := hex.EncodeToString(buf)
	if err := DB.Save(&models.SystemConfig{Key: keyJWTSecret, Value: secret}).Error; err != nil {
		return nil, err
	}
	log.Printf("未设置 JWT_SECRET，已生成随机密钥并保存在数据库中")
	return []byte(secret), nil
}
