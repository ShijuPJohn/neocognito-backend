package utils

import (
	"github.com/golang-jwt/jwt/v5"
	"neocognito-backend/models"
	"time"
)

// JwtGenerate TODO parameter type mongo.InsertOneResult is a workaround
func JwtGenerate(user models.User) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["email"] = user.Email
	claims["role"] = user.Role
	claims["id"] = user.ID
	claims["exp"] = time.Now().Add(time.Hour * 72).Unix()
	claims["issued_at"] = time.Now().Unix()
	t, err := token.SignedString([]byte(Secret))
	return t, err

}
