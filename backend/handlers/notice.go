package handlers

import (
	"net/http"
	"ops-portal/config"
	"ops-portal/models"

	"github.com/gin-gonic/gin"
)

// GetActiveNotice 获取当前激活的公告
func GetActiveNotice(c *gin.Context) {
	var notice models.Notice
	result := config.DB.Where("active = ?", true).First(&notice)
	if result.Error != nil {
		c.JSON(http.StatusOK, gin.H{"content": ""})
		return
	}
	c.JSON(http.StatusOK, notice)
}

// GetNotices 获取所有公告
func GetNotices(c *gin.Context) {
	var notices []models.Notice
	config.DB.Find(&notices)
	c.JSON(http.StatusOK, notices)
}

// CreateNotice 创建公告
func CreateNotice(c *gin.Context) {
	var notice models.Notice
	if err := c.ShouldBindJSON(&notice); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if notice.Active {
		config.DB.Model(&models.Notice{}).Where("active = ?", true).Update("active", false)
	}

	if err := config.DB.Create(&notice).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "创建公告失败"})
		return
	}
	c.JSON(http.StatusOK, notice)
}

// UpdateNotice 更新公告
func UpdateNotice(c *gin.Context) {
	id := c.Param("id")
	if id == "" || id == "undefined" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的公告ID"})
		return
	}

	var notice models.Notice
	if err := config.DB.First(&notice, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "公告不存在"})
		return
	}

	// 绑定新的数据
	var updateData struct {
		Content string `json:"content"`
		Active  bool   `json:"active"`
	}

	if err := c.ShouldBindJSON(&updateData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 更新字段
	notice.Content = updateData.Content
	notice.Active = updateData.Active

	// 如果更新为激活状态，则停用其他公告
	if notice.Active {
		config.DB.Model(&models.Notice{}).Where("id != ? AND active = ?", id, true).Update("active", false)
	}

	// 保存更新
	if err := config.DB.Save(&notice).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "更新公告失败"})
		return
	}

	c.JSON(http.StatusOK, notice)
}

// DeleteNotice 删除公告
func DeleteNotice(c *gin.Context) {
	id := c.Param("id")
	if id == "" || id == "undefined" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的公告ID"})
		return
	}

	if err := config.DB.Delete(&models.Notice{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "删除公告失败"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "删除成功"})
}
