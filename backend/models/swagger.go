package models

// LoginRequest 登录请求
type LoginRequest struct {
	Username string `json:"username" binding:"required" example:"admin"`
	Password string `json:"password" binding:"required" example:"admin123"`
}

// LoginResponse 登录响应
type LoginResponse struct {
	Token string `json:"token" example:"eyJhbGciOiJIUzI1NiIs..."`
}

// ChangePasswordRequest 修改密码请求
type ChangePasswordRequest struct {
	OldPassword string `json:"oldPassword" binding:"required" example:"oldpass123"`
	NewPassword string `json:"newPassword" binding:"required" example:"newpass123"`
}

// ErrorResponse 错误响应
type ErrorResponse struct {
	Message string `json:"message" example:"Invalid credentials"`
}

// SuccessResponse 成功响应
type SuccessResponse struct {
	Message string `json:"message" example:"Operation successful"`
}
