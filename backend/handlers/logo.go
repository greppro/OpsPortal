package handlers

import (
	"fmt"
	"log"
	"net/http"
	"ops-portal/config"
	"ops-portal/models"
	"os"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

const uploadDir = "uploads/logos"

func init() {
	if err := os.MkdirAll(uploadDir, 0755); err != nil {
			panic(err)
	}
}

// 上传 Logo
func UploadLogo(c *gin.Context) {
	file, err := c.FormFile("logo")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No file uploaded"})
		return
	}

	ext := filepath.Ext(file.Filename)
	if ext != ".png" && ext != ".jpg" && ext != ".jpeg" && ext != ".svg" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Only PNG, JPG, JPEG and SVG files are allowed"})
		return
	}

	// 删除旧的 logo 文件和记录
	var oldLogo models.Logo
	if err := config.DB.Last(&oldLogo).Error; err == nil {
		// 如果找到旧的 logo
		if oldLogo.URL != "" {
			oldFilePath := strings.TrimPrefix(oldLogo.URL, "/")
			if err := os.Remove(oldFilePath); err != nil {
				if !os.IsNotExist(err) {
					log.Printf("Failed to delete old logo file: %v", err)
				}
			}
		}
		// 删除旧的数据库记录
		if err := config.DB.Delete(&oldLogo).Error; err != nil {
			log.Printf("Failed to delete old logo record: %v", err)
		}
	}

	filename := fmt.Sprintf("%s%s", uuid.New().String(), ext)
	filepath := filepath.Join(uploadDir, filename)

	if err := c.SaveUploadedFile(file, filepath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save file"})
		return
	}

	// 保存到数据库
	logo := models.Logo{
		URL:  fmt.Sprintf("/uploads/logos/%s", filename),
		Name: file.Filename,
	}

	if err := config.DB.Create(&logo).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save logo info"})
		return
	}

	c.JSON(http.StatusOK, logo)
}

// 获取当前 Logo
func GetLogo(c *gin.Context) {
	var logo models.Logo
	if err := config.DB.Last(&logo).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{"url": ""})
		return
	}
	c.JSON(http.StatusOK, logo)
}

// 添加删除 Logo 的处理器
func DeleteLogo(c *gin.Context) {
	// 获取当前 logo
	var logo models.Logo
	if err := config.DB.Last(&logo).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusOK, gin.H{"message": "No logo to delete"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to query logo"})
		return
	}

	// 删除文件
	if logo.URL != "" {
		filePath := strings.TrimPrefix(logo.URL, "/")
		if err := os.Remove(filePath); err != nil {
			if !os.IsNotExist(err) {
				log.Printf("Failed to delete logo file: %v", err)
			}
		}
	}

	// 删除数据库记录
	if err := config.DB.Delete(&logo).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete logo"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Logo deleted successfully"})
} 