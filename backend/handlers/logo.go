package handlers

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"log"
	"mime/multipart"
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

const (
	uploadDir     = "uploads/logos"
	logoURLPrefix = "/uploads/logos/"
	maxLogoBytes  = 2 << 20 // 2MB
)

var allowedLogoExts = map[string]bool{".png": true, ".jpg": true, ".jpeg": true, ".svg": true}

var errLogoContent = errors.New("文件内容与扩展名不符")

// validateLogoContent 按文件头判断内容是否真的是对应格式的图片，不只看扩展名。
func validateLogoContent(fh *multipart.FileHeader, ext string) error {
	f, err := fh.Open()
	if err != nil {
		return err
	}
	defer f.Close()

	head := make([]byte, 4096)
	n, err := io.ReadFull(f, head)
	if err != nil && !errors.Is(err, io.ErrUnexpectedEOF) && !errors.Is(err, io.EOF) {
		return err
	}
	head = head[:n]

	switch ext {
	case ".png":
		if http.DetectContentType(head) != "image/png" {
			return errLogoContent
		}
	case ".jpg", ".jpeg":
		if http.DetectContentType(head) != "image/jpeg" {
			return errLogoContent
		}
	case ".svg":
		// SVG 是 XML 文本，开头部分应当出现 <svg 标签。
		// SVG 里的脚本由 /uploads 的 CSP 响应头禁止执行（见 main.go）。
		if !bytes.Contains(bytes.ToLower(head), []byte("<svg")) {
			return errLogoContent
		}
	}
	return nil
}

// removeLogoFile 删除 Logo 文件。只按文件名在上传目录里删，不信任数据库里的路径。
func removeLogoFile(logoURL string) {
	if !strings.HasPrefix(logoURL, logoURLPrefix) {
		return
	}
	name := filepath.Base(logoURL)
	if name == "." || name == "/" || name == ".." {
		return
	}
	if err := os.Remove(filepath.Join(uploadDir, name)); err != nil && !os.IsNotExist(err) {
		log.Printf("Failed to delete logo file: %v", err)
	}
}

// 上传 Logo
func UploadLogo(c *gin.Context) {
	// 限制请求体大小，超大文件在读取阶段就拒绝，不落临时文件
	c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, maxLogoBytes+64<<10)
	file, err := c.FormFile("logo")
	if err != nil {
		var tooLarge *http.MaxBytesError
		if errors.As(err, &tooLarge) {
			c.JSON(http.StatusRequestEntityTooLarge, gin.H{"error": "Logo 文件不能超过 2MB"})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{"error": "No file uploaded"})
		return
	}
	if file.Size > maxLogoBytes {
		c.JSON(http.StatusRequestEntityTooLarge, gin.H{"error": "Logo 文件不能超过 2MB"})
		return
	}

	ext := strings.ToLower(filepath.Ext(file.Filename))
	if !allowedLogoExts[ext] {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Only PNG, JPG, JPEG and SVG files are allowed"})
		return
	}
	if err := validateLogoContent(file, ext); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": errLogoContent.Error()})
		return
	}

	if err := os.MkdirAll(uploadDir, 0755); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save file"})
		return
	}
	filename := fmt.Sprintf("%s%s", uuid.New().String(), ext)
	dst := filepath.Join(uploadDir, filename)
	if err := c.SaveUploadedFile(file, dst); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save file"})
		return
	}

	// 新文件保存成功后再替换旧 Logo，上传失败时旧 Logo 保持不变
	var oldLogos []models.Logo
	config.DB.Find(&oldLogos)

	logo := models.Logo{
		URL:  logoURLPrefix + filename,
		Name: filepath.Base(file.Filename),
	}
	if err := config.DB.Create(&logo).Error; err != nil {
		_ = os.Remove(dst)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save logo info"})
		return
	}
	for _, old := range oldLogos {
		removeLogoFile(old.URL)
		if err := config.DB.Delete(&models.Logo{}, old.ID).Error; err != nil {
			log.Printf("Failed to delete old logo record: %v", err)
		}
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
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusOK, gin.H{"message": "No logo to delete"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to query logo"})
		return
	}

	// 删除数据库记录后再删文件
	if err := config.DB.Delete(&logo).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete logo"})
		return
	}
	removeLogoFile(logo.URL)

	c.JSON(http.StatusOK, gin.H{"message": "Logo deleted successfully"})
}
