package handlers

import (
	"net/http"
	"ops-portal/config"
	"ops-portal/models"
	"strings"

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

type categoryRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Icon        string `json:"icon"`
}

func (r *categoryRequest) normalize() {
	r.Name = strings.TrimSpace(r.Name)
	r.Description = strings.TrimSpace(r.Description)
	r.Icon = strings.TrimSpace(r.Icon)
}

func (r *categoryRequest) validate() string {
	switch {
	case tooLong(r.Name, 50):
		return "分类名称不能超过 50 个字符"
	case tooLong(r.Description, 200):
		return "分类描述不能超过 200 个字符"
	case tooLong(r.Icon, 50):
		return "图标名称不能超过 50 个字符"
	}
	return ""
}

func categoryNameTaken(name string, excludeID uint) bool {
	var count int64
	q := config.DB.Model(&models.Category{}).Where("name = ?", name)
	if excludeID != 0 {
		q = q.Where("id <> ?", excludeID)
	}
	q.Count(&count)
	return count > 0
}

// CreateCategory 创建分类
func CreateCategory(c *gin.Context) {
	var req categoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	req.normalize()
	if req.Name == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "分类名称不能为空"})
		return
	}
	if msg := req.validate(); msg != "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": msg})
		return
	}
	if categoryNameTaken(req.Name, 0) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "分类名称已存在"})
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
	id, ok := parseID(c)
	if !ok {
		return
	}
	var cat models.Category
	if err := config.DB.First(&cat, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "分类不存在"})
		return
	}
	var req categoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	req.normalize()
	if msg := req.validate(); msg != "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": msg})
		return
	}
	if req.Name != "" && req.Name != cat.Name {
		if categoryNameTaken(req.Name, id) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "分类名称已存在"})
			return
		}
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
	id, ok := parseID(c)
	if !ok {
		return
	}
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
