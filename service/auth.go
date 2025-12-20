package service

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"webssh/model"
)

var SecretKey = []byte("your-secret-key-change-it")

type Claims struct {
	Username string `json:"username"`
	jwt.RegisteredClaims
}

// Login 验证账号密码并签发Token
func Login(username, password string) (string, error) {
	var user model.User
	result := model.DB.Where("username = ? AND password = ?", username, password).First(&user)
	if result.Error != nil {
		return "", errors.New("用户名或密码错误")
	}

	expirationTime := time.Now().Add(24 * time.Hour)
	claims := &Claims{
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(SecretKey)
	return tokenString, err
}
