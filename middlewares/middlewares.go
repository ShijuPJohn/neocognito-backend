package middlewares

import (
	jwtware "github.com/gofiber/contrib/jwt"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"neocognito-backend/models"
	"neocognito-backend/utils"
)

func NotFound(c *fiber.Ctx) error {
	return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
		"Status":  "Error",
		"Message": "Not Found",
	}) // => 404 "Not Found"
}

func Protected() func(*fiber.Ctx) error {
	return jwtware.New(jwtware.Config{
		SigningKey:     jwtware.SigningKey{Key: []byte(utils.Secret)},
		ErrorHandler:   jwtError,
		SuccessHandler: handleSuccess,
		ContextKey:     "user",
	})
}

func jwtError(c *fiber.Ctx, err error) error {
	if err.Error() == "Missing or malformed JWT" {
		c.Status(fiber.StatusBadRequest)
		return c.JSON(fiber.Map{"status": "error", "message": "Missing or malformed JWT", "data": nil})

	} else {
		c.Status(fiber.StatusUnauthorized)
		return c.JSON(fiber.Map{"status": "error", "message": "Invalid or expired JWT", "data": nil})
	}
}
func handleSuccess(c *fiber.Ctx) error {
	user := c.Locals("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	issuedAt := int64(claims["issued_at"].(float64))
	userId := claims["id"].(string)
	var userFromDB models.User
	idObject, err := primitive.ObjectIDFromHex(userId)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "error", "error": err.Error()})
	}
	err = utils.Mg.Db.Collection("users").FindOne(c.Context(), bson.M{"_id": idObject}).Decode(&userFromDB)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "error", "error": "User not found"})
	}
	if userFromDB.PasswordChangedAt != nil && issuedAt < userFromDB.PasswordChangedAt.Unix() {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Token is no longer valid due to password change",
		})
	}
	return c.Next()
}
