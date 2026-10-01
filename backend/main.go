package main

import (
	"context"
	"log"
	"net"
	"ops-portal/config"
	"ops-portal/handlers"
	"ops-portal/middleware"
	"ops-portal/probe"
	"ops-portal/utils"
	"os"

	_ "ops-portal/docs"

	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title OpsPortal API
// @version 1.0
// @description OpsPortal运维导航平台的API文档
// @host localhost:3000
// @BasePath /api
// @schemes http https
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization

// 初始化必要的目录
func initDirectories() {
	dirs := []string{
		"uploads",
		"uploads/logos",
	}

	for _, dir := range dirs {
		err := os.MkdirAll(dir, 0755)
		if err != nil {
			log.Fatalf("Failed to create directory %s: %v", dir, err)
		}
	}
}

func main() {
	settings := config.LoadSettings()

	// 初始化目录
	initDirectories()

	// 初始化数据库、管理员账号和 JWT 密钥
	config.InitDB(settings)
	secret, err := config.InitJWTSecret(settings)
	if err != nil {
		log.Fatal("Failed to initialize JWT secret:", err)
	}
	utils.SetJWTSecret(secret)

	// 启动内置可用性检测
	p := probe.New(probe.Config{
		Enabled:            settings.Probe.Enabled,
		Interval:           settings.Probe.Interval,
		Timeout:            settings.Probe.Timeout,
		Concurrency:        settings.Probe.Concurrency,
		InsecureSkipVerify: settings.Probe.InsecureSkipVerify,
	}, handlers.ProbeTargets)
	handlers.SetProber(p)
	go p.Run(context.Background())

	r := setupRouter(settings)
	if err := r.Run(net.JoinHostPort(settings.Host, settings.Port)); err != nil {
		log.Fatal(err)
	}
}

// uploadSecurityHeaders 上传目录只放图片：禁止内容嗅探和脚本执行，
// 防止有人直接打开上传的 SVG 时执行其中嵌入的脚本。
func uploadSecurityHeaders() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("X-Content-Type-Options", "nosniff")
		c.Header("Content-Security-Policy", "default-src 'none'; img-src 'self' data:; style-src 'unsafe-inline'; sandbox")
		c.Next()
	}
}

func setupRouter(settings config.Settings) *gin.Engine {
	r := gin.Default()

	// 跨域默认关闭：前端通过同源反向代理访问后端。确实需要跨域时用 CORS_ALLOW_ORIGINS 指定来源
	if len(settings.CORSAllowOrigins) > 0 {
		corsConfig := cors.Config{
			AllowOrigins:  settings.CORSAllowOrigins,
			AllowMethods:  []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
			AllowHeaders:  []string{"Origin", "Content-Type", "Accept", "Authorization"},
			ExposeHeaders: []string{"Content-Length"},
			MaxAge:        12 * time.Hour,
		}
		if err := corsConfig.Validate(); err != nil {
			log.Fatalf("CORS_ALLOW_ORIGINS 配置无效：%v（示例：https://portal.example.com,https://other.example.com）", err)
		}
		r.Use(cors.New(corsConfig))
	}

	// Swagger 文档路由
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler,
		ginSwagger.URL("http://localhost:3000/swagger/doc.json"),
		ginSwagger.DefaultModelsExpandDepth(-1),
		ginSwagger.DocExpansion("none"),
		ginSwagger.InstanceName("swagger"),
	))

	// 公开路由组
	public := r.Group("")
	{
		// 健康检查接口
		public.GET("/api/health", func(c *gin.Context) {
			c.JSON(200, gin.H{
				"status":  "ok",
				"message": "OpsPortal backend is running",
			})
		})

		// 认证相关
		public.POST("/api/auth/login", handlers.Login)

		// 公开查询接口
		public.GET("/api/sites", handlers.GetTools)
		public.GET("/api/environments", handlers.GetEnvironmentsByProject)
		public.GET("/api/projects", handlers.GetProjects)

		// 公告相关的公开接口（只返回当前激活的公告）
		public.GET("/api/notices/active", handlers.GetActiveNotice)

		// Logo 相关的公开接口
		public.GET("/api/logo", handlers.GetLogo)

		// 系统配置（站点标题）
		public.GET("/api/config/site-title", handlers.GetSiteTitle)
		public.GET("/api/site-config", handlers.GetSiteConfig)

		// 分类列表（公开，供首页侧边栏与分类管理页使用）
		public.GET("/api/categories", handlers.GetCategories)

		// 可用性检测结果（首页状态点使用）
		public.GET("/api/probe/status", handlers.GetProbeStatus)
	}

	// 需要认证的路由组
	auth := r.Group("/api")
	auth.Use(middleware.AuthMiddleware())
	{
		// 用户相关
		auth.POST("/auth/change-password", handlers.ChangePassword)

		// 网址管理
		auth.POST("/sites", handlers.CreateTool)
		auth.PUT("/sites/:id", handlers.UpdateTool)
		auth.DELETE("/sites/:id", handlers.DeleteTool)

		// 环境管理
		auth.POST("/environments", handlers.CreateEnvironment)
		auth.PUT("/environments/:id", handlers.UpdateEnvironment)
		auth.DELETE("/environments/:id", handlers.DeleteEnvironment)

		// 项目管理
		auth.POST("/projects", handlers.CreateProject)
		auth.PUT("/projects/:id", handlers.UpdateProject)
		auth.DELETE("/projects/:id", handlers.DeleteProject)

		// 公告管理相关接口（列表含未激活的草稿，只对管理员开放）
		auth.GET("/notices", handlers.GetNotices)
		auth.POST("/notices", handlers.CreateNotice)
		auth.PUT("/notices/:id", handlers.UpdateNotice)
		auth.DELETE("/notices/:id", handlers.DeleteNotice)

		// 上传 Logo
		auth.POST("/upload/logo", handlers.UploadLogo)
		auth.DELETE("/logo", handlers.DeleteLogo)

		// 系统配置（站点标题）
		auth.PUT("/config/site-title", handlers.PutSiteTitle)
		auth.PUT("/site-config", handlers.PutSiteConfig)

		// 分类管理
		auth.POST("/categories", handlers.CreateCategory)
		auth.PUT("/categories/:id", handlers.UpdateCategory)
		auth.DELETE("/categories/:id", handlers.DeleteCategory)

		// 立即检测某个工具
		auth.POST("/probe/check/:id", handlers.CheckToolNow)
	}

	// Prometheus 对接：指标抓取和 http_sd 服务发现（设置 METRICS_TOKEN 后需要 Bearer Token）
	prom := r.Group("", handlers.MetricsAuth(settings.MetricsToken))
	{
		prom.GET("/metrics", handlers.Metrics)
		prom.GET("/prometheus/targets", handlers.PrometheusTargets)
	}

	// 上传文件（Logo）
	uploads := r.Group("/uploads", uploadSecurityHeaders())
	uploads.Static("/", "./uploads")

	return r
}
