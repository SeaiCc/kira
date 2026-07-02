package utils

import (
	"errors"
	"time"

	"kira-go/internal/config"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

// CreateToken 创建 JWT token
// data 示例：map[string]any{"sub": "username", "admin": true}
func CreateToken(data map[string]any) (string, error) {
	expireTime := time.Now().Add(time.Duration(config.C.JWTExpireHours) * time.Hour)

	// 使用 jwt.MapClaims（类似 Python 的 dict）
	claims := jwt.MapClaims{
		"sub":   data["sub"],   // username（对应 Python 的 sub 字段）
		"admin": data["admin"], // is_admin（对应 Python 的 admin 字段）
		"exp":   expireTime.Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(config.C.JWTSecret))
}

// DecodeToken 解析 JWT token
// 返回 jwt.MapClaims
func DecodeToken(tokenString string) (jwt.MapClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, jwt.MapClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(config.C.JWTSecret), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, nil
	}
	return nil, errors.New("invalid token")
}

// HashPassword 密码哈希
func HashPassword(password string) (string, error) {
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedBytes), nil
}

// VerifyPassword 验证密码
func VerifyPassword(password, hashedPassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err == nil
}
