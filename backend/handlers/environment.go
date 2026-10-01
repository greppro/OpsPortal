package handlers

import (
	"net/http"
	"ops-portal/config"
	"ops-portal/models"
	"strings"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type environmentRequest struct {
	Name      string `json:"name"`
	Label     string `json:"label"`
	ProjectID uint   `json:"project_id"`
}

func (r *environmentRequest) normalize() {
	r.Name = strings.TrimSpace(r.Name)
	r.Label = strings.TrimSpace(r.Label)
}

func (r *environmentRequest) validate() string {
	switch {
	case r.ProjectID == 0:
		return "请选择项目"
	case r.Name == "":
		return "环境标识不能为空"
	case strings.ContainsAny(r.Name, " \t\r\n"):
		return "环境标识不能包含空白字符"
	case tooLong(r.Name, 50):
		return "环境标识不能超过 50 个字符"
	case r.Label == "":
		return "环境名称不能为空"
	case tooLong(r.Label, 50):
		return "环境名称不能超过 50 个字符"
	}
	return ""
}

// countToolsInEnvironment 统计某个项目某个环境下的工具数（工具按项目标识 + 环境标识关联）。
func countToolsInEnvironment(project, env string) (int64, error) {
	var count int64
	err := config.DB.Model(&models.Tool{}).
		Where("project = ? AND environment = ?", project, env).
		Count(&count).Error
	return count, err
}

// GetEnvironmentsByProject 获取环境列表
func GetEnvironmentsByProject(c *gin.Context) {
	var environments []models.Environment
	projectID := c.Query("project_id")

	query := config.DB.Model(&models.Environment{})

	if projectID != "" {
		query = query.Where("project_id = ?", projectID)
	}

	if err := query.Find(&environments).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取环境列表失败"})
		return
	}

	c.JSON(http.StatusOK, environments)
}

// CreateEnvironment 创建环境
func CreateEnvironment(c *gin.Context) {
	var req environmentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	req.normalize()
	if msg := req.validate(); msg != "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": msg})
		return
	}

	// 获取项目信息
	var project models.Project
	if err := config.DB.First(&project, req.ProjectID).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "项目不存在"})
		return
	}

	// 检查同一项目下环境标识是否重复
	var count int64
	if err := config.DB.Model(&models.Environment{}).
		Where("project_id = ? AND name = ?", req.ProjectID, req.Name).
		Count(&count).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "创建环境失败"})
		return
	}
	if count > 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "同一项目下环境标识不能重复"})
		return
	}

	env := models.Environment{
		Name:      req.Name,
		Label:     req.Label,
		ProjectID: project.ID,
		Project:   project.Name,
	}
	if err := config.DB.Create(&env).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "创建环境失败"})
		return
	}

	c.JSON(http.StatusOK, env)
}

// UpdateEnvironment 更新环境。修改环境标识时，本项目下挂在这个环境上的工具一并改名。
func UpdateEnvironment(c *gin.Context) {
	id, ok := parseID(c)
	if !ok {
		return
	}
	var env models.Environment
	if err := config.DB.First(&env, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "环境不存在"})
		return
	}

	var req environmentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	req.normalize()
	if msg := req.validate(); msg != "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": msg})
		return
	}

	var project models.Project
	if err := config.DB.First(&project, req.ProjectID).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "项目不存在"})
		return
	}

	projectChanged := env.ProjectID != req.ProjectID
	nameChanged := env.Name != req.Name
	if projectChanged {
		count, err := countToolsInEnvironment(env.Project, env.Name)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "更新环境失败"})
			return
		}
		if count > 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "该环境下还有工具，不能移动到其他项目"})
			return
		}
	}

	// 检查同一项目下环境标识是否重复
	var count int64
	if err := config.DB.Model(&models.Environment{}).
		Where("project_id = ? AND name = ? AND id <> ?", req.ProjectID, req.Name, id).
		Count(&count).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "更新环境失败"})
		return
	}
	if count > 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "同一项目下环境标识不能重复"})
		return
	}

	oldProject, oldName := env.Project, env.Name
	env.Name = req.Name
	env.Label = req.Label
	env.ProjectID = project.ID
	env.Project = project.Name
	err := config.DB.Transaction(func(tx *gorm.DB) error {
		if nameChanged && !projectChanged {
			if err := tx.Model(&models.Tool{}).
				Where("project = ? AND environment = ?", oldProject, oldName).
				Update("environment", req.Name).Error; err != nil {
				return err
			}
		}
		return tx.Save(&env).Error
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "更新环境失败"})
		return
	}

	c.JSON(http.StatusOK, env)
}

// DeleteEnvironment 删除环境
func DeleteEnvironment(c *gin.Context) {
	id, ok := parseID(c)
	if !ok {
		return
	}
	var env models.Environment

	if err := config.DB.First(&env, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "环境不存在"})
		return
	}

	// 只检查本项目下是否有工具使用此环境（其他项目的同名环境不影响）
	count, err := countToolsInEnvironment(env.Project, env.Name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "检查环境使用情况失败"})
		return
	}
	if count > 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "该环境下存在工具，无法删除"})
		return
	}

	if err := config.DB.Delete(&env).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "删除环境失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "删除成功"})
}
