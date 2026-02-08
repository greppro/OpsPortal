package main

import (
	"log"
	"ops-portal/config"
	"ops-portal/handlers"
	"ops-portal/middleware"
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
		"data",
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
	// 初始化目录
	initDirectories()

	// 初始化数据库
	config.InitDB()

	r := gin.Default()

	// 配置跨域
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

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
				"status": "ok",
				"message": "OpsPortal backend is running",
			})
		})

		// 认证相关
		public.POST("/api/auth/login", handlers.Login)

		// 公开查询接口
		public.GET("/api/sites", handlers.GetTools)
		public.GET("/api/environments", handlers.GetEnvironmentsByProject)
		public.GET("/api/projects", handlers.GetProjects)

		// 公告相关的公开接口
		public.GET("/api/notices/active", handlers.GetActiveNotice) // 获取当前激活的公告
		public.GET("/api/notices", handlers.GetNotices)             // 获取公告列表

		// Logo 相关的公开接口
		public.GET("/api/logo", handlers.GetLogo)

		// 系统配置（站点标题）
		public.GET("/api/config/site-title", handlers.GetSiteTitle)

		// 分类列表（公开，供首页侧边栏与分类管理页使用）
		public.GET("/api/categories", handlers.GetCategories)
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

		// 公告管理相关接口
		auth.POST("/notices", handlers.CreateNotice)
		auth.PUT("/notices/:id", handlers.UpdateNotice)
		auth.DELETE("/notices/:id", handlers.DeleteNotice)

		// 上传 Logo
		auth.POST("/upload/logo", handlers.UploadLogo)
		auth.DELETE("/logo", handlers.DeleteLogo)

		// 系统配置（站点标题）
		auth.PUT("/config/site-title", handlers.PutSiteTitle)

		// 分类管理
		auth.POST("/categories", handlers.CreateCategory)
		auth.PUT("/categories/:id", handlers.UpdateCategory)
		auth.DELETE("/categories/:id", handlers.DeleteCategory)
	}

	// 添加静态文件服务
	r.Static("/uploads", "./uploads")

	r.Run(":8080")
}
