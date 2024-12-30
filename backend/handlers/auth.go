package handlers

import (
	"net/http"
	"ops-portal/models"
	"ops-portal/utils"

	"github.com/gin-gonic/gin"
)

// LoginRequest 登录请求
type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
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

	// 验证用户名和密码
	if req.Username != "admin" || req.Password != "admin123" {
		c.JSON(http.StatusUnauthorized, models.ErrorResponse{Message: "用户名或密码错误"})
		return
	}

	// 生成 token
	token, err := utils.GenerateToken(req.Username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Message: "生成token失败"})
		return
	}

	c.JSON(http.StatusOK, models.LoginResponse{Token: token})
}

// @Summary 修改密码
// @Description 修改用户密码
// @Tags 认证
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param body body models.ChangePasswordRequest true "密码信息"
// @Success 200 {object} models.SuccessResponse
// @Failure 400 {object} models.ErrorResponse
// @Failure 401 {object} models.ErrorResponse
// @Router /auth/change-password [post]
func ChangePassword(c *gin.Context) {
	var req struct {
		OldPassword string `json:"oldPassword" binding:"required"`
		NewPassword string `json:"newPassword" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Message: err.Error()})
		return
	}

	// 验证旧密码
	if req.OldPassword != "admin123" {
		c.JSON(http.StatusUnauthorized, models.ErrorResponse{Message: "旧密码错误"})
		return
	}

	// TODO: 更新密码
	c.JSON(http.StatusOK, models.SuccessResponse{Message: "密码修改成功"})
}
