package main

import (
	"ops-portal/config"
	"ops-portal/handlers"
	"ops-portal/middleware"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	// 初始化数据库
	config.InitDB()

	r := gin.Default()

	// 配置跨域
	corsConfig := cors.DefaultConfig()
	corsConfig.AllowOrigins = []string{"http://localhost:3000", "http://127.0.0.1:3000"}
	corsConfig.AllowMethods = []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}
	corsConfig.AllowHeaders = []string{"Origin", "Content-Type", "Accept", "Authorization", "X-User"}
	r.Use(cors.New(corsConfig))

	// 公开路由组
	public := r.Group("/api")
	{
		// 认证相关
		public.POST("/auth/login", handlers.Login)

		// 公开查询接口
		public.GET("/sites", handlers.GetTools)
		public.GET("/environments", handlers.GetEnvironmentsByProject)
		public.GET("/projects", handlers.GetProjects)
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
	}

	r.Run(":8080")
}
