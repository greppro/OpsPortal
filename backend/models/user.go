package models

type User struct {
	ID       uint   `gorm:"primarykey"`
	Username string `gorm:"unique"`
	// Password 存 bcrypt 哈希；旧版本的明文密码会在启动时自动转成哈希
	Password string
	// TokenVersion 修改密码时加一，签发时间早于修改的 token 随之失效
	TokenVersion int `gorm:"not null;default:0"`
}
