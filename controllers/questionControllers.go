package controllers

import (
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
	"neocognito-backend/models"
	"neocognito-backend/utils"
	"strconv"
	"time"
)

func CreateQuestion(c *fiber.Ctx) error {
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
	fmt.Println(q)
	q.CreatedAt = time.Now()
	q.EditedAt = time.Now()
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
func CreateQuestionSet(c *fiber.Ctx) error {
	return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"sfdfd": "Sdfdsf"})
}
