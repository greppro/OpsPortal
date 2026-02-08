package handlers

import (
	"encoding/json"
	"net/http"
	"ops-portal/config"
	"ops-portal/models"
	"os"
	"path/filepath"
	"time"

	"github.com/gin-gonic/gin"
)

// GetTools 获取工具列表
func GetTools(c *gin.Context) {
	var tools []models.Tool

	// 获取查询参数
	env := c.Query("env")
	project := c.Query("project")

	// 构建查询
	query := config.DB.Model(&models.Tool{})

	// 添加筛选条件
	if env != "" {
		query = query.Where("environment = ?", env)
	}
	if project != "" {
		query = query.Where("project = ?", project)
	}

	// 执行查询
	if err := query.Find(&tools).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取工具列表失败"})
		return
	}

	// #region agent log
	if len(tools) > 0 {
		sample := make([]map[string]interface{}, 0, 3)
		for i := 0; i < len(tools) && i < 3; i++ {
			t := tools[i]
			sample = append(sample, map[string]interface{}{"id": t.ID, "name": t.Name, "url": t.URL, "hasUrl": t.URL != ""})
		}
		line, _ := json.Marshal(map[string]interface{}{
			"location": "tool.go:GetTools", "message": "tools response sample", "data": map[string]interface{}{"count": len(tools), "sample": sample},
			"timestamp": time.Now().UnixMilli(), "hypothesisId": "H2,H4",
		})
		logPath := filepath.Join("..", ".cursor", "debug.log")
		if f, err := os.OpenFile(logPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644); err == nil {
			f.Write(append(line, '\n'))
			f.Close()
		}
	}
	// #endregion

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
		c.JSON(http.StatusNotFound, gin.H{"error": "工具不存���"})
		return
	}

	if err := config.DB.Delete(&tool).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "删除工具失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "删除成功"})
}
