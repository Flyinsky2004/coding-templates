/*
 * @Author: Flyinsky w2084151024@gmail.com
 * @Description: None
 */
package utility

import (
	"errors"
	"time"

	"example/config"

	"github.com/golang-jwt/jwt/v5"
)

/**
 * @author Flyinsky
 * @email w2084151024@gmail.com
 * @date 2024/12/24 11:32
 */

// 使用config中的密钥
var secretKey = []byte(config.Config.JWT.SecretKey)

type CustomClaims struct {
	UserID   int    `json:"user_id"`
	Username string `json:"username"`
	jwt.RegisteredClaims
}

// GenerateToken 生成JWT
func GenerateToken(userID int, username string) (string, error) {
	claims := CustomClaims{
		UserID:   userID,
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			// 使用config中的过期时间
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(config.Config.JWT.ExpireTime) * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(secretKey)
}

// ParseToken 解析JWT
func ParseToken(tokenString string) (*CustomClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return secretKey, nil
	})
	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*CustomClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("invalid token")
}
