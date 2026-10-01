package utils

import (
	"crypto/rand"
	"math/big"

	"golang.org/x/crypto/bcrypt"
)

// LegacyDefaultPassword 是旧版本写死的默认密码，仍在使用时提示管理员修改。
const LegacyDefaultPassword = "admin123"

// MaxPasswordBytes 是 bcrypt 能处理的最大长度。
const MaxPasswordBytes = 72

// HashPassword 使用 bcrypt 生成密码哈希。
func HashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

// CheckPassword 校验明文密码与哈希是否匹配。
func CheckPassword(hash, password string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)) == nil
}

// IsPasswordHash 判断字符串是否已经是 bcrypt 哈希（用于迁移旧库里的明文密码）。
func IsPasswordHash(s string) bool {
	_, err := bcrypt.Cost([]byte(s))
	return err == nil
}

const randomAlphabet = "ABCDEFGHJKLMNPQRSTUVWXYZabcdefghijkmnopqrstuvwxyz23456789"

// RandomString 生成指定长度的随机字符串（去掉了易混淆的字符）。
func RandomString(n int) (string, error) {
	buf := make([]byte, n)
	max := big.NewInt(int64(len(randomAlphabet)))
	for i := range buf {
		idx, err := rand.Int(rand.Reader, max)
		if err != nil {
			return "", err
		}
		buf[i] = randomAlphabet[idx.Int64()]
	}
	return string(buf), nil
}
