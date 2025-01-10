package main

import (
	"ops-portal/config"
	"ops-portal/handlers"
	"ops-portal/middleware"

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

func main() {
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
		// 认证相关
		public.POST("/api/auth/login", handlers.Login)

		// 公开查询接口
		public.GET("/api/sites", handlers.GetTools)
		public.GET("/api/environments", handlers.GetEnvironmentsByProject)
		public.GET("/api/projects", handlers.GetProjects)

		// 公告相关的公开接口
		public.GET("/api/notices/active", handlers.GetActiveNotice) // 获取当前激活的公告
		public.GET("/api/notices", handlers.GetNotices)             // 获取公告列表
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
	}

	r.Run(":8080")
}
