package security

import (
	"coderhub/conf"
	"github.com/golang-jwt/jwt/v4"
	"strconv"
	"time"
)

// 定义一个密钥，用于签名和验证 Token
var jwtKey = []byte(conf.JWTSecret)

type CustomClaims struct {
	UserID int64 `json:"userId"`
	jwt.RegisteredClaims
}

func GenerateAuthorization(userID int64) (string, error) {
	claims := CustomClaims{
		UserID: userID, // 将 userID 放入自定义字段中
		RegisteredClaims: jwt.RegisteredClaims{
			ID:        strconv.FormatInt(userID, 10),                          // 唯一标识
			Issuer:    "CoderHub",                                             // 签发者
			IssuedAt:  jwt.NewNumericDate(time.Now()),                         // 签发时间
			Subject:   "auth",                                                 // 主题
			Audience:  jwt.ClaimStrings{"CoderHub"},                           // 接收者
			NotBefore: jwt.NewNumericDate(time.Now()),                         // 生效时间
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(7 * 24 * time.Hour)), // 过期时间
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString(jwtKey)
	if err != nil {
		return "", err
	}

	return signedToken, nil
}
