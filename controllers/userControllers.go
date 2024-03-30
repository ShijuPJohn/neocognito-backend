package controllers

import (
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
	"neocognito-backend/models"
	"neocognito-backend/utils"
)

func Index(c *fiber.Ctx) error {
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"status": "success", "page": "index page"})
}

func GetAllUsers(c *fiber.Ctx) error {
	query := bson.D{{}}
	cursor, err := utils.Mg.Db.Collection("users").Find(c.Context(), query)
	if err != nil {
		return c.Status(500).SendString(err.Error())
	}

	var users = make([]models.User, 0)

	if err := cursor.All(c.Context(), &users); err != nil {
		return c.Status(500).SendString(err.Error())

	}
	return c.JSON(users)
}
func CreateUser(c *fiber.Ctx) error {
	u := new(models.User)

	if err := c.BodyParser(u); err != nil {
		return c.Status(400).SendString(err.Error())
	}
	type ConfirmPasswordStruct struct {
		ConfirmPassword string `json:"confirmPassword" validate:"required"`
	}
	confirmPassword := new(ConfirmPasswordStruct)
	if err := c.BodyParser(confirmPassword); err != nil {
		return c.Status(400).SendString(err.Error())
	}
	if confirmPassword.ConfirmPassword == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Confirm Password not present",
		})
	}
	if u.Password != confirmPassword.ConfirmPassword {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Passwords don't match",
		})
	}
	validate := validator.New()
	err := validate.Struct(u)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": err.Error(),
		})
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "Error",
			"message": err.Error(),
		})
	}
	u.Password = string(hash)
	insertionResult, err := utils.Mg.Db.Collection("users").InsertOne(c.Context(), u)
	if err != nil {
		if _, ok := err.(mongo.WriteException); ok {
			return c.Status(fiber.StatusCreated).JSON(fiber.Map{
				"status":  "error",
				"message": "Email already registered"})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Server Error"})
	}
	if err != nil {
		return c.Status(500).SendString(err.Error())
	}
	token := utils.JwtGenerate(u.Email, *insertionResult, u.Role)
	fmt.Println(token)
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"status":  "success",
		"message": "User Created",
		"_id":     insertionResult.InsertedID,
		"user":    u})
}
