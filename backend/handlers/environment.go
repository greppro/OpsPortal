package handlers

import (
	"fmt"
	"net/http"
	"ops-portal/config"
	"ops-portal/models"

	"github.com/gin-gonic/gin"
)

// GetEnvironmentsByProject 获取环境列表
func GetEnvironmentsByProject(c *gin.Context) {
	var environments []models.Environment
	projectID := c.Query("project_id")

	query := config.DB.Model(&models.Environment{})

	if projectID != "" {
		query = query.Where("project_id = ?", projectID)
	}

	if err := query.Find(&environments).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取环境列表失败"})
		return
	}

	// 打印调试信息
	fmt.Printf("Fetched environments: %+v\n", environments)

	c.JSON(http.StatusOK, environments)
}

// CreateEnvironment 创建环境
func CreateEnvironment(c *gin.Context) {
	var env models.Environment
	if err := c.ShouldBindJSON(&env); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 获取项目信息
	var project models.Project
	if err := config.DB.First(&project, env.ProjectID).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "项目不存在"})
		return
	}
	env.Project = project.Name

	// 检查同一项目下环境标识是否重复
	var count int64
	if err := config.DB.Model(&models.Environment{}).
		Where("project_id = ? AND name = ?", env.ProjectID, env.Name).
		Count(&count).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "创建环境失败"})
		return
	}

	if count > 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "同一项目下环境标识不能重复"})
		return
	}

	if err := config.DB.Create(&env).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "创建环境失败"})
		return
	}

	c.JSON(http.StatusOK, env)
}

// UpdateEnvironment 更新环境
func UpdateEnvironment(c *gin.Context) {
	id := c.Param("id")
	var env models.Environment

	if err := config.DB.First(&env, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "环境不存在"})
		return
	}

	var newEnv models.Environment
	if err := c.ShouldBindJSON(&newEnv); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 如果更改了项目
	if env.ProjectID != newEnv.ProjectID {
		var project models.Project
		if err := config.DB.First(&project, newEnv.ProjectID).Error; err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "项目不存在"})
			return
		}
		newEnv.Project = project.Name
	} else {
		newEnv.Project = env.Project
	}

	// 检查同一项目下环境标识是否重复
	var count int64
	if err := config.DB.Model(&models.Environment{}).
		Where("project_id = ? AND name = ? AND id != ?", newEnv.ProjectID, newEnv.Name, id).
		Count(&count).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "更新环境失败"})
		return
	}

	if count > 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "同一项目下环境标识不能重复"})
		return
	}

	if err := config.DB.Model(&env).Updates(newEnv).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "更新环境失败"})
		return
	}

	c.JSON(http.StatusOK, env)
}

// DeleteEnvironment 删除环境
func DeleteEnvironment(c *gin.Context) {
	id := c.Param("id")
	var env models.Environment

	if err := config.DB.First(&env, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "环境不存在"})
		return
	}

	// 检查是否有工具使用此环境
	var count int64
	if err := config.DB.Model(&models.Tool{}).Where("environment = ?", env.Name).Count(&count).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "检查环境使用情况失败"})
		return
	}

	if count > 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "该环境下存在工具，无法删除"})
		return
	}

	if err := config.DB.Delete(&env).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "删除环境失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "删除成功"})
}
