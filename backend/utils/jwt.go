package utils

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// TokenTTL 是登录 token 的有效期。
const TokenTTL = 24 * time.Hour

// jwtSecret 在启动时由 SetJWTSecret 注入（来自 JWT_SECRET 或数据库里自动生成的密钥）。
var jwtSecret []byte

var errSecretNotSet = errors.New("jwt secret is not initialized")

// SetJWTSecret 设置签名密钥。
func SetJWTSecret(secret []byte) {
	jwtSecret = secret
}

// Claims 是 token 的载荷。TokenVersion 与用户表里的版本号对应，修改密码后旧 token 失效。
type Claims struct {
	Username     string `json:"username"`
	TokenVersion int    `json:"ver"`
	jwt.RegisteredClaims
}

// GenerateToken 生成 JWT token。
func GenerateToken(username string, tokenVersion int) (string, error) {
	if len(jwtSecret) == 0 {
		return "", errSecretNotSet
	}
	now := time.Now()
	claims := &Claims{
		Username:     username,
		TokenVersion: tokenVersion,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(TokenTTL)),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}

// ParseToken 解析并校验 JWT token，只接受 HS256 签名。
func ParseToken(tokenString string) (*Claims, error) {
	if len(jwtSecret) == 0 {
		return nil, errSecretNotSet
	}
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	}, jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}))
	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, jwt.ErrSignatureInvalid
}
