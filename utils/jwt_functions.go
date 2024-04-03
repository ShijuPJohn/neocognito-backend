package utils

import (
	"github.com/golang-jwt/jwt/v5"
	"go.mongodb.org/mongo-driver/mongo"
	"neocognito-backend/models"
	"time"
)

// JwtGenerate TODO parameter type mongo.InsertOneResult is a workaround
func JwtGenerate(user models.User, id mongo.InsertOneResult) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["email"] = user.Email
	claims["role"] = user.Role
	claims["id"] = user.ID
	claims["exp"] = time.Now().Add(time.Hour * 72).Unix()
	t, err := token.SignedString([]byte(Secret))
	return t, err

}
