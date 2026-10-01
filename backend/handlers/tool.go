package handlers

import (
	"net/http"
	"ops-portal/config"
	"ops-portal/models"
	"ops-portal/probe"
	"strings"

	"github.com/gin-gonic/gin"
)

// toolRequest 是创建/更新网址时允许写入的字段（id、时间戳等不从请求体读取）。
type toolRequest struct {
	Name          string `json:"name"`
	Description   string `json:"description"`
	URL           string `json:"url"`
	Environment   string `json:"environment"`
	Project       string `json:"project"`
	Category      string `json:"category"`
	Icon          string `json:"icon"`
	ProbeDisabled bool   `json:"probe_disabled"`
}

func (r *toolRequest) normalize() {
	r.Name = strings.TrimSpace(r.Name)
	r.Description = strings.TrimSpace(r.Description)
	r.URL = strings.TrimSpace(r.URL)
	r.Environment = strings.TrimSpace(r.Environment)
	r.Project = strings.TrimSpace(r.Project)
	r.Category = strings.TrimSpace(r.Category)
	r.Icon = strings.TrimSpace(r.Icon)
}

// validate 返回第一条校验错误，没有错误时返回空字符串。
// checkURL 为 false 时不校验地址格式（更新时地址没变，兼容旧数据）。
func (r *toolRequest) validate(checkURL bool) string {
	switch {
	case r.Name == "":
		return "名称不能为空"
	case r.URL == "":
		return "地址不能为空"
	case r.Project == "":
		return "项目不能为空"
	case r.Environment == "":
		return "环境不能为空"
	case tooLong(r.Name, 100):
		return "名称不能超过 100 个字符"
	case tooLong(r.Description, 500):
		return "描述不能超过 500 个字符"
	case len(r.URL) > 2048:
		return "地址不能超过 2048 个字符"
	case tooLong(r.Category, 50):
		return "分类不能超过 50 个字符"
	case tooLong(r.Icon, 500):
		return "图标不能超过 500 个字符"
	}
	if checkURL {
		if _, err := probe.NormalizeURL(r.URL); err != nil {
			return "地址格式不正确，只支持 http 和 https 地址"
		}
	}
	return ""
}

func (r *toolRequest) applyTo(t *models.Tool) {
	t.Name = r.Name
	t.Description = r.Description
	t.URL = r.URL
	t.Environment = r.Environment
	t.Project = r.Project
	t.Category = r.Category
	t.Icon = r.Icon
	t.ProbeDisabled = r.ProbeDisabled
}

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

	c.JSON(http.StatusOK, tools)
}

// CreateTool 创建工具
func CreateTool(c *gin.Context) {
	var req toolRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	req.normalize()
	if msg := req.validate(true); msg != "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": msg})
		return
	}

	var tool models.Tool
	req.applyTo(&tool)
	if err := config.DB.Create(&tool).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "创建工具失败"})
		return
	}

	triggerProbe(tool)
	c.JSON(http.StatusOK, tool)
}

// UpdateTool 更新工具
func UpdateTool(c *gin.Context) {
	id, ok := parseID(c)
	if !ok {
		return
	}
	var tool models.Tool
	if err := config.DB.First(&tool, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "工具不存在"})
		return
	}

	var req toolRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	req.normalize()
	if msg := req.validate(req.URL != strings.TrimSpace(tool.URL)); msg != "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": msg})
		return
	}

	req.applyTo(&tool)
	if err := config.DB.Save(&tool).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "更新工具失败"})
		return
	}

	triggerProbe(tool)
	c.JSON(http.StatusOK, tool)
}

// DeleteTool 删除工具
func DeleteTool(c *gin.Context) {
	id, ok := parseID(c)
	if !ok {
		return
	}
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
