package handlers

import (
	"net/http"
	"ops-portal/config"
	"ops-portal/models"

	"github.com/gin-gonic/gin"
)

// GetCategories 获取分类列表（供分类管理页与侧边栏使用）
func GetCategories(c *gin.Context) {
	var list []models.Category
	if err := config.DB.Order("name").Find(&list).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取分类列表失败"})
		return
	}
	c.JSON(http.StatusOK, list)
}

// CreateCategory 创建分类
func CreateCategory(c *gin.Context) {
	var req struct {
		Name        string `json:"name" binding:"required"`
		Description string `json:"description"`
		Icon        string `json:"icon"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	cat := models.Category{
		Name:        req.Name,
		Description: req.Description,
		Icon:        req.Icon,
	}
	if err := config.DB.Create(&cat).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "创建分类失败"})
		return
	}
	c.JSON(http.StatusOK, cat)
}

// UpdateCategory 更新分类
func UpdateCategory(c *gin.Context) {
	id := c.Param("id")
	var cat models.Category
	if err := config.DB.First(&cat, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "分类不存在"})
		return
	}
	var req struct {
		Name        string `json:"name"`
		Description string `json:"description"`
		Icon        string `json:"icon"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if req.Name != "" {
		cat.Name = req.Name
	}
	cat.Description = req.Description
	cat.Icon = req.Icon
	if err := config.DB.Save(&cat).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "更新分类失败"})
		return
	}
	c.JSON(http.StatusOK, cat)
}

// DeleteCategory 删除分类（仅删 DB 记录，不修改工具的 category 字段）
func DeleteCategory(c *gin.Context) {
	id := c.Param("id")
	var cat models.Category
	if err := config.DB.First(&cat, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "分类不存在"})
		return
	}
	if err := config.DB.Delete(&cat).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "删除分类失败"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "ok"})
}
