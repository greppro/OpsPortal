package handlers

import (
	"fmt"
	"net/http"
	"ops-portal/config"
	"ops-portal/models"

	"github.com/gin-gonic/gin"
)

// GetProjects 获取项目列表
func GetProjects(c *gin.Context) {
	var projects []models.Project
	if err := config.DB.Find(&projects).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取项目列表失败"})
		return
	}

	// 打印调试信息
	fmt.Printf("Fetched projects: %+v\n", projects)

	c.JSON(http.StatusOK, projects)
}

// CreateProject 创建项目
func CreateProject(c *gin.Context) {
	var project models.Project
	if err := c.ShouldBindJSON(&project); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 如果设置为默认项目，取消其他默认项目
	if project.IsDefault {
		config.DB.Model(&models.Project{}).Where("is_default = ?", true).Update("is_default", false)
	}

	if err := config.DB.Create(&project).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "创建项目失败"})
		return
	}

	c.JSON(http.StatusOK, project)
}

// UpdateProject 更新项目
func UpdateProject(c *gin.Context) {
	id := c.Param("id")
	var project models.Project

	if err := config.DB.First(&project, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "项目不存在"})
		return
	}

	var req models.Project
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 如果设置为默认项目，取消其他默认项目
	if req.IsDefault {
		config.DB.Model(&models.Project{}).Where("id != ?", id).Update("is_default", false)
	}

	project.Name = req.Name
	project.Label = req.Label
	project.IsDefault = req.IsDefault

	if err := config.DB.Save(&project).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "更新项目失败"})
		return
	}

	c.JSON(http.StatusOK, project)
}

// DeleteProject 删除项目
func DeleteProject(c *gin.Context) {
	id := c.Param("id")
	var project models.Project

	if err := config.DB.First(&project, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "项目不存在"})
		return
	}

	// 检查是否有环境使用此项目
	var count int64
	if err := config.DB.Model(&models.Environment{}).Where("project_id = ?", id).Count(&count).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "检查项目使用情况失败"})
		return
	}

	if count > 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "该项目下存在环境，无法删除"})
		return
	}

	if err := config.DB.Delete(&project).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "删除项目失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "删除成功"})
}
