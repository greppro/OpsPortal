package handlers

import (
	"net/http"
	"ops-portal/config"
	"ops-portal/models"

	"github.com/gin-gonic/gin"
)

// GetTools 获取工具列表
func GetTools(c *gin.Context) {
	var tools []models.Tool
	if err := config.DB.Find(&tools).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取工具列表失败"})
		return
	}
	c.JSON(http.StatusOK, tools)
}

// CreateTool 创建工具
func CreateTool(c *gin.Context) {
	var tool models.Tool
	if err := c.ShouldBindJSON(&tool); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := config.DB.Create(&tool).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "创建工具失败"})
		return
	}

	c.JSON(http.StatusOK, tool)
}

// UpdateTool 更新工具
func UpdateTool(c *gin.Context) {
	id := c.Param("id")
	var tool models.Tool

	if err := config.DB.First(&tool, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "工具不存在"})
		return
	}

	if err := c.ShouldBindJSON(&tool); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := config.DB.Save(&tool).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "更新工具失败"})
		return
	}

	c.JSON(http.StatusOK, tool)
}

// DeleteTool 删除工具
func DeleteTool(c *gin.Context) {
	id := c.Param("id")
	var tool models.Tool

	if err := config.DB.First(&tool, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "工具不存在"})
		return
	}

	if err := config.DB.Delete(&tool).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "删除工具失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "删除成功"})
}
