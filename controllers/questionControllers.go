package controllers

import (
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
	"neocognito-backend/models"
	"neocognito-backend/utils"
	"strconv"
	"strings"
	"time"
)

func CreateQuestion(c *fiber.Ctx) error {
	user := c.Locals("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	q := new(models.Question)
	if err := c.BodyParser(q); err != nil {
		return c.Status(400).SendString(err.Error())
	}
	validate := validator.New()
	err := validate.Struct(q)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": err.Error(),
		})
	}
	q.CreatedAt = time.Now()
	q.EditedAt = time.Now()
	q.CreatedById = claims["id"].(string)
	insertionResult, err := utils.Mg.Db.Collection("questions").InsertOne(c.Context(), q)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": err.Error(),
		})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"status": "success", "question": q, "id": insertionResult.InsertedID})

}
func GetQuestions(c *fiber.Ctx) error {
	findOptions := options.Find()
	mapOfQuery := bson.D{{}}
	projection := bson.M{}
	page, _ := strconv.Atoi(c.Query("page", "1"))
	numberOfQuestions, _ := strconv.Atoi(c.Query("noQs", "1"))
	findOptions.SetLimit(int64(numberOfQuestions))
	findOptions.SetSkip((int64(page) - 1) * int64(numberOfQuestions))
	if c.Query("category") != "" {
		mapOfQuery = append(mapOfQuery, bson.E{Key: "category", Value: c.Query("category")})
	}
	if c.Query("subject") != "" {
		mapOfQuery = append(mapOfQuery, bson.E{Key: "subject", Value: c.Query("subject")})
	}
	if c.Query("exam") != "" {
		mapOfQuery = append(mapOfQuery, bson.E{Key: "exam", Value: c.Query("exam")})
	}
	if c.Query("sort") == "asc" {
		findOptions.SetSort(bson.D{{Key: "createdAt", Value: 1}})
	}
	if c.Query("sort") == "desc" {
		findOptions.SetSort(bson.D{{Key: "createdAt", Value: -1}})
	}
	if c.Query("fields") != "" {
		fields := strings.Split(c.Query("fields"), ",")
		for _, field := range fields {
			projection[field] = 1
		}
	}
	findOptions.SetProjection(projection)
	cursor, err := utils.Mg.Db.Collection("questions").Find(c.Context(), mapOfQuery, findOptions)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": err.Error()})
	}
	var arrayOfQuestions []models.Question
	err = cursor.All(c.Context(), &arrayOfQuestions)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": err.Error()})
	}

	//count, err := utils.Mg.Db.Collection("questions").CountDocuments(c.Context(), mapOfQuery)
	//if err != nil {
	//	return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": err.Error()})
	//}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"questions": arrayOfQuestions})
}
func EditQuestion(c *fiber.Ctx) error {
	idParam := c.Params("id")
	qID, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		return handleError(c, fiber.StatusBadRequest, err.Error())
	}

	update := bson.M{}
	if err := c.BodyParser(&update); err != nil {
		return handleError(c, fiber.StatusBadRequest, err.Error())
	}

	// Remove the _id field from the update to prevent accidentally changing it
	delete(update, "_id")

	// Construct the update query
	filter := bson.M{"_id": qID}
	updateQuery := bson.M{"$set": update}

	result, err := utils.Mg.Db.Collection("questions").UpdateOne(c.Context(), filter, updateQuery)
	if err != nil {
		return handleError(c, fiber.StatusBadRequest, err.Error())
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"status": "success", "result": result})
}
func DeleteQuestion(c *fiber.Ctx) error {
	idParam := c.Params("id")
	qID, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		return handleError(c, fiber.StatusBadRequest, err.Error())
	}

	result, err := utils.Mg.Db.Collection("questions").DeleteOne(c.Context(), bson.M{"_id": qID})
	if err != nil {
		return handleError(c, fiber.StatusInternalServerError, err.Error())
	}

	if result.DeletedCount == 0 {
		return handleError(c, fiber.StatusNotFound, "Question not found")
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"status": "success", "message": "Question deleted successfully"})
}

func handleError(c *fiber.Ctx, statusCode int, errorMessage string) error {
	return c.Status(statusCode).JSON(fiber.Map{
		"status":  "error",
		"message": errorMessage,
	})
}
