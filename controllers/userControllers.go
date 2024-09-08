package controllers

import (
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
	"neocognito-backend/models"
	"neocognito-backend/utils"
	"time"
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
	validate := validator.New()
	err := validate.Struct(u)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": err.Error(),
		})
	}
	u.Role = "user"
	currentTime := time.Now()
	u.CreatedAt = &currentTime
	u.UpdatedAt = &currentTime
	u.PasswordChangedAt = &currentTime
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
	_, ok := insertionResult.InsertedID.(primitive.ObjectID)
	if !ok {
		fmt.Println("Failed to convert inserted ID to ObjectID")
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}
	token, err := utils.JwtGenerate(*u)
	if err != nil {
		fmt.Println(err.Error())
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"status":  "success",
		"message": "User Created",
		"token":   token,
		"_id":     insertionResult.InsertedID,
		"user":    u})
}
func LoginUser(c *fiber.Ctx) error {
	type loginModel struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	loginObject := new(loginModel)
	err := c.BodyParser(&loginObject)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "error", "error": err.Error()})
	}
	var userFromDB models.User
	err = utils.Mg.Db.Collection("users").FindOne(c.Context(), bson.M{"email": loginObject.Email}).Decode(&userFromDB)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "error", "error": err.Error()})
	}
	err = bcrypt.CompareHashAndPassword([]byte(userFromDB.Password), []byte(loginObject.Password))
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"status": "error", "message": "invalid credentials"})
	}
	//TODO parameter type mongo.InsertOneResult is a workaround
	token, err := utils.JwtGenerate(userFromDB)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"status": "error", "message": "invalid credentials"})
	}
	if token == "" {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "error", "message": "internal server error"})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"status": "success", "token": token})

}
func GetUserDetails(c *fiber.Ctx) error {
	fmt.Println("Authorization Passed")
	fmt.Println(c.Params("user"))
	var userFromDB models.User
	idObject, err := primitive.ObjectIDFromHex(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "error", "error": err.Error()})
	}
	err = utils.Mg.Db.Collection("users").FindOne(c.Context(), bson.M{"_id": idObject}).Decode(&userFromDB)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "error", "error": err.Error()})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"status": "success", "user": userFromDB})

}
