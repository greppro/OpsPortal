package models

// LoginRequest 登录请求
type LoginRequest struct {
	Username string `json:"username" binding:"required" example:"admin"`
	Password string `json:"password" binding:"required" example:"your-password"`
}

// LoginResponse 登录响应
type LoginResponse struct {
	Token string   `json:"token" example:"eyJhbGciOiJIUzI1NiIs..."`
	User  UserInfo `json:"user"`
	// DefaultPassword 为 true 表示仍在使用旧版本的默认密码，前端会提示立即修改
	DefaultPassword bool `json:"default_password"`
}

// UserInfo 用户信息
type UserInfo struct {
	Username string `json:"username" example:"admin"`
}

// ChangePasswordRequest 修改密码请求
type ChangePasswordRequest struct {
	OldPassword string `json:"oldPassword" binding:"required" example:"oldpass123"`
	NewPassword string `json:"newPassword" binding:"required" example:"newpass123"`
}

// ChangePasswordResponse 修改密码响应，旧 token 失效，前端需要换成这里返回的新 token
type ChangePasswordResponse struct {
	Message string `json:"message" example:"密码修改成功"`
	Token   string `json:"token" example:"eyJhbGciOiJIUzI1NiIs..."`
}

// ErrorResponse 错误响应
type ErrorResponse struct {
	Message string `json:"message" example:"Invalid credentials"`
}

// SuccessResponse 成功响应
type SuccessResponse struct {
	Message string `json:"message" example:"Operation successful"`
}
