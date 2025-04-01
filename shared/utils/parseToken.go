package utils

import "github.com/dgrijalva/jwt-go"

// ParseUserIDFromToken 从 JWT token 中解析用户 ID
func ParseUserIDFromToken(tokenStr string, secretKey string) (string, error) {
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		return []byte(secretKey), nil
	})
	if err != nil {
		return "", err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		userID, ok := claims["user_id"].(string)
		if !ok {
			return "", err
		}
		return userID, nil
	}

	return "", err
}
