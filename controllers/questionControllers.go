package controllers

import (
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"neocognito-backend/models"
	"neocognito-backend/utils"
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
