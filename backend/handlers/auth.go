package handlers

import (
	"net/http"
	"ops-portal/config"
	"ops-portal/models"
	"ops-portal/utils"
	"sync"
	"unicode/utf8"

	"github.com/gin-gonic/gin"
)

// LoginRequest 登录请求
type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

const minPasswordLength = 8

var (
	dummyHashOnce sync.Once
	dummyHash     string
)

// equalizeLoginTiming 用户不存在时也做一次 bcrypt 比较，避免通过响应时间判断用户名是否存在。
func equalizeLoginTiming(password string) {
	dummyHashOnce.Do(func() {
		dummyHash, _ = utils.HashPassword("opsportal-dummy-password")
	})
	utils.CheckPassword(dummyHash, password)
}

// @Summary 用户登录
// @Description 用户登录接口
// @Tags 认证
// @Accept json
// @Produce json
// @Param body body LoginRequest true "登录信息"
// @Success 200 {object} models.LoginResponse
// @Failure 400 {object} models.ErrorResponse
// @Failure 401 {object} models.ErrorResponse
// @Router /auth/login [post]
func Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Message: err.Error()})
		return
	}

	var user models.User
	if err := config.DB.Where("username = ?", req.Username).First(&user).Error; err != nil {
		equalizeLoginTiming(req.Password)
		c.JSON(http.StatusUnauthorized, models.ErrorResponse{Message: "用户名或密码错误"})
		return
	}
	if !utils.CheckPassword(user.Password, req.Password) {
		c.JSON(http.StatusUnauthorized, models.ErrorResponse{Message: "用户名或密码错误"})
		return
	}

	token, err := utils.GenerateToken(user.Username, user.TokenVersion)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Message: "生成token失败"})
		return
	}

	c.JSON(http.StatusOK, models.LoginResponse{
		Token: token,
		User: models.UserInfo{
			Username: user.Username,
		},
		DefaultPassword: req.Password == utils.LegacyDefaultPassword,
	})
}

// @Summary 修改密码
// @Description 修改当前登录用户的密码，成功后旧 token 全部失效，响应里返回新 token
// @Tags 认证
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param body body models.ChangePasswordRequest true "密码信息"
// @Success 200 {object} models.ChangePasswordResponse
// @Failure 400 {object} models.ErrorResponse
// @Failure 401 {object} models.ErrorResponse
// @Router /auth/change-password [post]
func ChangePassword(c *gin.Context) {
	var req models.ChangePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Message: err.Error()})
		return
	}

	if utf8.RuneCountInString(req.NewPassword) < minPasswordLength {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Message: "新密码至少 8 位"})
		return
	}
	if len(req.NewPassword) > utils.MaxPasswordBytes {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Message: "新密码不能超过 72 个字节"})
		return
	}
	if req.NewPassword == req.OldPassword {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Message: "新密码不能与旧密码相同"})
		return
	}

	var user models.User
	if err := config.DB.Where("username = ?", c.GetString("username")).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, models.ErrorResponse{Message: "用户不存在"})
		return
	}
	if !utils.CheckPassword(user.Password, req.OldPassword) {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Message: "旧密码错误"})
		return
	}

	hash, err := utils.HashPassword(req.NewPassword)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Message: "修改密码失败"})
		return
	}
	newVersion := user.TokenVersion + 1
	if err := config.DB.Model(&models.User{}).Where("id = ?", user.ID).Updates(map[string]interface{}{
		"password":      hash,
		"token_version": newVersion,
	}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Message: "修改密码失败"})
		return
	}

	token, err := utils.GenerateToken(user.Username, newVersion)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Message: "生成token失败"})
		return
	}
	c.JSON(http.StatusOK, models.ChangePasswordResponse{Message: "密码修改成功", Token: token})
}
