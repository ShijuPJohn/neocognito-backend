package utils

import (
	secretmanager "cloud.google.com/go/secretmanager/apiv1"
	"context"
	"github.com/golang-jwt/jwt/v5"
	"go.mongodb.org/mongo-driver/mongo"
	secretmanagerpb "google.golang.org/genproto/googleapis/cloud/secretmanager/v1"
	"time"
)

func JwtGenerate(email string, id mongo.InsertOneResult, role string) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)
	claims["email"] = email
	claims["role"] = role
	claims["id"] = id.InsertedID
	claims["exp"] = time.Now().Add(time.Hour * 72).Unix()
	t, err := token.SignedString([]byte(Secret))

}

func GetSecret() string {
	name := "projects/924880194744/secrets/dd2_jw_secret/versions/1"
	ctx := context.Background()
	client, err := secretmanager.NewClient(ctx)
	if err != nil {
		return ""
	}
	defer client.Close()

	req := &secretmanagerpb.AccessSecretVersionRequest{
		Name: name,
	}
	result, err := client.AccessSecretVersion(ctx, req)
	if err != nil {
		return ""
	}
	return string(result.Payload.Data)
}
