package utils

import (
	"github.com/dgrijalva/jwt-go"
	"strconv"
)

// ParseUserIDFromToken 从 JWT token 中解析用户 ID
func ParseUserIDFromToken(tokenStr string, secretKey string) (string, error) {
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		return []byte(secretKey), nil
	})
	if err != nil {
		return "", err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		// 尝试从 RegisteredClaims 的 ID 字段获取用户 ID
		if idStr, ok := claims["jti"].(string); ok {
			return idStr, nil
		}
		// 如果没有找到，再尝试从自定义的 UserID 字段获取
		if userID, ok := claims["UserID"].(float64); ok {
			return strconv.FormatInt(int64(userID), 10), nil
		}
	}

	return "", err
}
