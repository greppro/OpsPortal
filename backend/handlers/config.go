package handlers

import (
	"net/http"
	"ops-portal/config"
	"ops-portal/models"

	"github.com/gin-gonic/gin"
)

const (
	keySiteTitle   = "site_title"
	defaultTitle   = "OpsPortal运维导航"
)

// GetSiteTitle 获取站点标题（公开）
// @Summary 获取站点标题
// @Tags config
// @Produce json
// @Success 200 {object} map[string]string
// @Router /api/config/site-title [get]
func GetSiteTitle(c *gin.Context) {
	var row models.SystemConfig
	err := config.DB.Where("key = ?", keySiteTitle).First(&row).Error
	if err != nil || row.Value == "" {
		c.JSON(http.StatusOK, gin.H{"site_title": defaultTitle})
		return
	}
	c.JSON(http.StatusOK, gin.H{"site_title": row.Value})
}

// SiteTitleRequest 站点标题请求体
type SiteTitleRequest struct {
	SiteTitle string `json:"site_title" binding:"required"`
}

// PutSiteTitle 更新站点标题（需认证）
// @Summary 更新站点标题
// @Tags config
// @Accept json
// @Produce json
// @Param body body SiteTitleRequest true "站点标题"
// @Success 200 {object} map[string]string
// @Router /api/config/site-title [put]
func PutSiteTitle(c *gin.Context) {
	var req SiteTitleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请提供 site_title"})
		return
	}
	row := models.SystemConfig{Key: keySiteTitle, Value: req.SiteTitle}
	err := config.DB.Save(&row).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "保存失败"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"site_title": row.Value})
}
