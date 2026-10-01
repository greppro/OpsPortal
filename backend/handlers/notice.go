package handlers

import (
	"net/http"
	"ops-portal/config"
	"ops-portal/models"
	"strings"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type noticeRequest struct {
	Content string `json:"content"`
	Active  bool   `json:"active"`
}

func (r *noticeRequest) validate() string {
	r.Content = strings.TrimSpace(r.Content)
	switch {
	case r.Content == "":
		return "公告内容不能为空"
	case tooLong(r.Content, 1000):
		return "公告内容不能超过 1000 个字符"
	}
	return ""
}

// GetActiveNotice 获取当前激活的公告（公开）
func GetActiveNotice(c *gin.Context) {
	var notice models.Notice
	result := config.DB.Where("active = ?", true).First(&notice)
	if result.Error != nil {
		c.JSON(http.StatusOK, gin.H{"content": ""})
		return
	}
	c.JSON(http.StatusOK, notice)
}

// GetNotices 获取所有公告（含未激活的草稿，需要登录）
func GetNotices(c *gin.Context) {
	var notices []models.Notice
	if err := config.DB.Find(&notices).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取公告列表失败"})
		return
	}
	c.JSON(http.StatusOK, notices)
}

// CreateNotice 创建公告
func CreateNotice(c *gin.Context) {
	var req noticeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if msg := req.validate(); msg != "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": msg})
		return
	}

	notice := models.Notice{Content: req.Content, Active: req.Active}
	err := config.DB.Transaction(func(tx *gorm.DB) error {
		if notice.Active {
			if err := tx.Model(&models.Notice{}).Where("active = ?", true).Update("active", false).Error; err != nil {
				return err
			}
		}
		if err := tx.Create(&notice).Error; err != nil {
			return err
		}
		// active 字段带 default:true，GORM 创建时会把 false 换成默认值，这里显式改回去
		if !req.Active {
			notice.Active = false
			return tx.Model(&notice).Update("active", false).Error
		}
		return nil
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "创建公告失败"})
		return
	}
	c.JSON(http.StatusOK, notice)
}

// UpdateNotice 更新公告
func UpdateNotice(c *gin.Context) {
	id, ok := parseID(c)
	if !ok {
		return
	}

	var notice models.Notice
	if err := config.DB.First(&notice, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "公告不存在"})
		return
	}

	var req noticeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if msg := req.validate(); msg != "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": msg})
		return
	}

	notice.Content = req.Content
	notice.Active = req.Active
	err := config.DB.Transaction(func(tx *gorm.DB) error {
		// 如果更新为激活状态，则停用其他公告
		if notice.Active {
			if err := tx.Model(&models.Notice{}).Where("id <> ? AND active = ?", id, true).Update("active", false).Error; err != nil {
				return err
			}
		}
		return tx.Save(&notice).Error
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "更新公告失败"})
		return
	}

	c.JSON(http.StatusOK, notice)
}

// DeleteNotice 删除公告
func DeleteNotice(c *gin.Context) {
	id, ok := parseID(c)
	if !ok {
		return
	}

	if err := config.DB.Delete(&models.Notice{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "删除公告失败"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "删除成功"})
}
