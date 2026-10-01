package handlers

import (
	"net/http"
	"ops-portal/config"
	"ops-portal/models"
	"regexp"
	"strings"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

var projectNamePattern = regexp.MustCompile(`^[a-z0-9-]+$`)

type projectRequest struct {
	Name      string `json:"name"`
	Label     string `json:"label"`
	IsDefault bool   `json:"is_default"`
}

func (r *projectRequest) normalize() {
	r.Name = strings.TrimSpace(r.Name)
	r.Label = strings.TrimSpace(r.Label)
}

func validateProjectName(name string) string {
	switch {
	case name == "":
		return "项目标识不能为空"
	case len(name) > 50:
		return "项目标识不能超过 50 个字符"
	case !projectNamePattern.MatchString(name):
		return "项目标识只能包含小写字母、数字和横线"
	}
	return ""
}

func validateProjectLabel(label string) string {
	switch {
	case label == "":
		return "项目名称不能为空"
	case tooLong(label, 50):
		return "项目名称不能超过 50 个字符"
	}
	return ""
}

func projectNameTaken(name string, excludeID uint) (bool, error) {
	var count int64
	q := config.DB.Model(&models.Project{}).Where("name = ?", name)
	if excludeID != 0 {
		q = q.Where("id <> ?", excludeID)
	}
	err := q.Count(&count).Error
	return count > 0, err
}

// GetProjects 获取项目列表
func GetProjects(c *gin.Context) {
	var projects []models.Project
	if err := config.DB.Find(&projects).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取项目列表失败"})
		return
	}

	c.JSON(http.StatusOK, projects)
}

// CreateProject 创建项目
func CreateProject(c *gin.Context) {
	var req projectRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	req.normalize()
	if msg := validateProjectName(req.Name); msg != "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": msg})
		return
	}
	if msg := validateProjectLabel(req.Label); msg != "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": msg})
		return
	}
	if taken, err := projectNameTaken(req.Name, 0); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "创建项目失败"})
		return
	} else if taken {
		c.JSON(http.StatusBadRequest, gin.H{"error": "项目标识已存在"})
		return
	}

	project := models.Project{Name: req.Name, Label: req.Label, IsDefault: req.IsDefault}
	err := config.DB.Transaction(func(tx *gorm.DB) error {
		// 如果设置为默认项目，取消其他默认项目
		if project.IsDefault {
			if err := tx.Model(&models.Project{}).Where("is_default = ?", true).Update("is_default", false).Error; err != nil {
				return err
			}
		}
		return tx.Create(&project).Error
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "创建项目失败"})
		return
	}

	c.JSON(http.StatusOK, project)
}

// UpdateProject 更新项目。修改项目标识时，环境和工具里按名称关联的字段一并更新。
func UpdateProject(c *gin.Context) {
	id, ok := parseID(c)
	if !ok {
		return
	}
	var project models.Project
	if err := config.DB.First(&project, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "项目不存在"})
		return
	}

	var req projectRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	req.normalize()
	oldName := project.Name
	nameChanged := req.Name != oldName
	if nameChanged {
		if msg := validateProjectName(req.Name); msg != "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": msg})
			return
		}
		if taken, err := projectNameTaken(req.Name, id); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "更新项目失败"})
			return
		} else if taken {
			c.JSON(http.StatusBadRequest, gin.H{"error": "项目标识已存在"})
			return
		}
	}
	if msg := validateProjectLabel(req.Label); msg != "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": msg})
		return
	}

	project.Name = req.Name
	project.Label = req.Label
	project.IsDefault = req.IsDefault
	err := config.DB.Transaction(func(tx *gorm.DB) error {
		// 如果设置为默认项目，取消其他默认项目
		if req.IsDefault {
			if err := tx.Model(&models.Project{}).Where("id <> ?", id).Update("is_default", false).Error; err != nil {
				return err
			}
		}
		if nameChanged {
			if err := tx.Model(&models.Environment{}).Where("project_id = ?", id).Update("project", req.Name).Error; err != nil {
				return err
			}
			if err := tx.Model(&models.Tool{}).Where("project = ?", oldName).Update("project", req.Name).Error; err != nil {
				return err
			}
		}
		return tx.Save(&project).Error
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "更新项目失败"})
		return
	}

	c.JSON(http.StatusOK, project)
}

// DeleteProject 删除项目
func DeleteProject(c *gin.Context) {
	id, ok := parseID(c)
	if !ok {
		return
	}
	var project models.Project

	if err := config.DB.First(&project, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "项目不存在"})
		return
	}

	// 检查是否有环境使用此项目
	var count int64
	if err := config.DB.Model(&models.Environment{}).Where("project_id = ?", id).Count(&count).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "检查项目使用情况失败"})
		return
	}
	if count > 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "该项目下存在环境，无法删除"})
		return
	}

	// 工具按项目标识关联，项目删掉后这些工具在首页就找不到了
	if err := config.DB.Model(&models.Tool{}).Where("project = ?", project.Name).Count(&count).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "检查项目使用情况失败"})
		return
	}
	if count > 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "该项目下存在工具，无法删除"})
		return
	}

	if err := config.DB.Delete(&project).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "删除项目失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "删除成功"})
}
