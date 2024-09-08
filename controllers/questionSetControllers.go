package controllers

//
//import (
//	"github.com/go-playground/validator/v10"
//	"github.com/gofiber/fiber/v2"
//	"github.com/golang-jwt/jwt/v4"
//	"go.mongodb.org/mongo-driver/bson"
//	"go.mongodb.org/mongo-driver/bson/primitive"
//	"go.mongodb.org/mongo-driver/mongo"
//	"go.mongodb.org/mongo-driver/mongo/options"
//	"neocognito-backend/models"
//	"neocognito-backend/utils"
//	"strconv"
//	"strings"
//)
//
//func CreateQuestionSet(c *fiber.Ctx) error {
//	user := c.Locals("user").(*jwt.Token)
//	claims := user.Claims.(jwt.MapClaims)
//	q := new(models.QuestionSet)
//	if err := c.BodyParser(&q); err != nil {
//		return c.Status(400).SendString(err.Error())
//	}
//	validate := validator.New()
//	err := validate.Struct(q)
//	if err != nil {
//		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
//			"status":  "error",
//			"message": err.Error(),
//		})
//	}
//	q.CreatedById = claims["id"].(string)
//	insertionResult, err := utils.Mg.Db.Collection("question_set").InsertOne(c.Context(), q)
//	if err != nil {
//		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
//			"status":  "error",
//			"message": err.Error(),
//		})
//	}
//	return c.Status(fiber.StatusOK).JSON(fiber.Map{"status": "success", "question_set": q, "id": insertionResult.InsertedID})
//}
//func GetQuestionSets(c *fiber.Ctx) error {
//	findOptions := options.Find()
//	filter := bson.D{}
//
//	// Parse query parameters
//	if subject := c.Query("subject"); subject != "" {
//		filter = append(filter, bson.E{Key: "subject", Value: subject})
//	}
//	if category := c.Query("category"); category != "" {
//		filter = append(filter, bson.E{Key: "category", Value: category})
//	}
//	if tags := c.Query("tags"); tags != "" {
//		tagsList := strings.Split(tags, ",")
//		filter = append(filter, bson.E{Key: "tags", Value: bson.M{"$in": tagsList}})
//	}
//	if language := c.Query("language"); language != "" {
//		filter = append(filter, bson.E{Key: "language", Value: language})
//	}
//
//	// Count of sets
//	count, _ := strconv.Atoi(c.Query("count", "-1"))
//	if count != -1 {
//		findOptions.SetLimit(int64(count))
//	}
//
//	cursor, err := utils.Mg.Db.Collection("question_set").Find(c.Context(), filter, findOptions)
//	if err != nil {
//		return handleError(c, fiber.StatusInternalServerError, err.Error())
//	}
//
//	var questionSets []models.QuestionSet
//	if err := cursor.All(c.Context(), &questionSets); err != nil {
//		return handleError(c, fiber.StatusInternalServerError, err.Error())
//	}
//
//	return c.Status(fiber.StatusOK).JSON(fiber.Map{"question_sets": questionSets})
//}
//func EditQuestionSet(c *fiber.Ctx) error {
//	idParam := c.Params("id")
//	qID, err := primitive.ObjectIDFromHex(idParam)
//	if err != nil {
//		return handleError(c, fiber.StatusBadRequest, err.Error())
//	}
//
//	update := bson.M{}
//	if err := c.BodyParser(&update); err != nil {
//		return handleError(c, fiber.StatusBadRequest, err.Error())
//	}
//
//	delete(update, "_id")
//
//	filter := bson.M{"_id": qID}
//	updateQuery := bson.M{"$set": update}
//
//	result, err := utils.Mg.Db.Collection("question_set").UpdateOne(c.Context(), filter, updateQuery)
//	if err != nil {
//		return handleError(c, fiber.StatusBadRequest, err.Error())
//	}
//
//	return c.Status(fiber.StatusOK).JSON(fiber.Map{"status": "success", "result": result})
//}
//
//func DeleteQuestionSet(c *fiber.Ctx) error {
//	idParam := c.Params("id")
//	qID, err := primitive.ObjectIDFromHex(idParam)
//	if err != nil {
//		return handleError(c, fiber.StatusBadRequest, err.Error())
//	}
//
//	result, err := utils.Mg.Db.Collection("question_set").DeleteOne(c.Context(), bson.M{"_id": qID})
//	if err != nil {
//		return handleError(c, fiber.StatusInternalServerError, err.Error())
//	}
//
//	if result.DeletedCount == 0 {
//		return handleError(c, fiber.StatusNotFound, "Question set not found")
//	}
//
//	return c.Status(fiber.StatusOK).JSON(fiber.Map{"status": "success", "message": "Question set deleted successfully"})
//}
//
//func GetQuestionSetByID(c *fiber.Ctx) error {
//	// Extract the ID from the request parameters
//	idParam := c.Params("id")
//
//	// Convert the ID string to ObjectID
//	qID, err := primitive.ObjectIDFromHex(idParam)
//	if err != nil {
//		return handleError(c, fiber.StatusBadRequest, err.Error())
//	}
//
//	// Define a filter to find the question set by ID
//	filter := bson.M{"_id": qID}
//
//	// Find the question set in the database
//	var questionSet models.QuestionSet
//	err = utils.Mg.Db.Collection("question_set").FindOne(c.Context(), filter).Decode(&questionSet)
//	if err != nil {
//		// If the question set is not found, return a 404 status code
//		if err == mongo.ErrNoDocuments {
//			return handleError(c, fiber.StatusNotFound, "Question set not found")
//		}
//		// If there's any other error, return a 500 status code
//		return handleError(c, fiber.StatusInternalServerError, err.Error())
//	}
//
//	// Return the question set as JSON response
//	return c.Status(fiber.StatusOK).JSON(fiber.Map{"question_set": questionSet})
//}
