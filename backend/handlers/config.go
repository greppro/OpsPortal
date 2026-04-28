package handlers

import (
	"net/http"
	"ops-portal/config"
	"ops-portal/models"
	"strings"

	"github.com/gin-gonic/gin"
)

const (
	keySiteTitle = "site_title"
	defaultTitle = "OpsPortal运维导航"
)

type siteConfigResponse struct {
	SystemName string `json:"system_name"`
}

// GetSiteTitle 获取站点标题（公开）
// @Summary 获取站点标题
// @Tags config
// @Produce json
// @Success 200 {object} map[string]string
// @Router /api/config/site-title [get]
func GetSiteTitle(c *gin.Context) {
	title := getSiteTitleValue()
	c.JSON(http.StatusOK, gin.H{"site_title": title})
}

func getSiteTitleValue() string {
	var row models.SystemConfig
	err := config.DB.Where("key = ?", keySiteTitle).First(&row).Error
	if err != nil || row.Value == "" {
		return defaultTitle
	}
	return row.Value
}

// GetSiteConfig 获取站点配置（公开）
func GetSiteConfig(c *gin.Context) {
	c.JSON(http.StatusOK, siteConfigResponse{SystemName: getSiteTitleValue()})
}

// SiteTitleRequest 站点标题请求体
type SiteTitleRequest struct {
	SiteTitle string `json:"site_title" binding:"required"`
}

type SiteConfigRequest struct {
	SystemName string `json:"system_name" binding:"required"`
}

func validateSystemName(systemName string) (string, bool) {
	name := strings.TrimSpace(systemName)
	return name, len([]rune(name)) >= 2 && len([]rune(name)) <= 40
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
	title, ok := validateSystemName(req.SiteTitle)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "系统名称长度需为 2 到 40 个字符"})
		return
	}
	row := models.SystemConfig{Key: keySiteTitle, Value: title}
	err := config.DB.Save(&row).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "保存失败"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"site_title": row.Value})
}

// PutSiteConfig 更新站点配置（需认证）
func PutSiteConfig(c *gin.Context) {
	var req SiteConfigRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请提供 system_name"})
		return
	}
	name, ok := validateSystemName(req.SystemName)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "系统名称长度需为 2 到 40 个字符"})
		return
	}
	row := models.SystemConfig{Key: keySiteTitle, Value: name}
	if err := config.DB.Save(&row).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "保存失败"})
		return
	}
	c.JSON(http.StatusOK, siteConfigResponse{SystemName: row.Value})
}
