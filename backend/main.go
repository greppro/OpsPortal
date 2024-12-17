package main

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"ops-portal/config"
	"ops-portal/handlers"
)

func main() {
	// 初始化数据库
	config.InitDB()

	r := gin.Default()
	
	// 配置跨域
	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"http://localhost:5178", "http://127.0.0.1:5178"}
	config.AllowMethods = []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}
	config.AllowHeaders = []string{"Origin", "Content-Type", "Accept", "Authorization"}
	r.Use(cors.New(config))
	
	// 路由组
	api := r.Group("/api")
	{
		api.GET("/tools", handlers.GetTools)
		api.POST("/tools", handlers.CreateTool)
		api.PUT("/tools/:id", handlers.UpdateTool)
		api.DELETE("/tools/:id", handlers.DeleteTool)
	}
	
	r.Run(":8080")
} 