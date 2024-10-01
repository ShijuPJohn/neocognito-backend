package controllers

import (
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
	"neocognito-backend/models"
	"neocognito-backend/utils"
	"strconv"
	"strings"
	"time"
)

func CreateQuestionSet(c *fiber.Ctx) error {
	type QuestionSetDT0Model struct {
		Questions   []string  `json:"questions" validate:"required"`
		Subject     string    `json:"subject" validate:"required"`
		Name        string    `json:"name" validate:"required"`
		Tags        []string  `json:"tags" validate:"required"`
		MarkList    []float64 `json:"mark_list"`
		Description string    `json:"description" validate:"required"`
		Exam        string    `json:"exam" validate:"required"`
	}
	var qSetDTO QuestionSetDT0Model
	user := c.Locals("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	q := new(models.QuestionSet)
	if err := c.BodyParser(&qSetDTO); err != nil {
		return c.Status(400).SendString(err.Error())
	}
	validate := validator.New()
	err := validate.Struct(qSetDTO)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": err.Error(),
		})
	}
	q.Name = qSetDTO.Name
	q.Description = qSetDTO.Description
	userID := claims["id"].(string)
	q.CreatedById = userID
	q.EditedByIds = []string{userID}
	currentTime := time.Now()
	q.CreatedAt = currentTime
	q.EditedAt = currentTime
	q.Subject = qSetDTO.Subject
	q.Tags = qSetDTO.Tags
	q.Exam = qSetDTO.Exam
	q.MarkList = qSetDTO.MarkList
	//var objectIDs []primitive.ObjectID
	//for _, idStr := range qSetDTO.Questions {
	//	id, err := primitive.ObjectIDFromHex(idStr)
	//	if err != nil {
	//		continue
	//	}
	//	objectIDs = append(objectIDs, id)
	//}
	//projection := bson.M{"id": 1, "correct_options": 1}
	//findOptions := options.Find()
	//findOptions.SetProjection(projection)
	//
	//filter := bson.M{"_id": bson.M{"$in": objectIDs}}
	//fmt.Println("before querying")
	//cursor, err := utils.Mg.Db.Collection("questions").Find(c.Context(), filter, findOptions)
	//fmt.Println("after querying")
	//var results []bson.M
	//if err := cursor.All(c.Context(), &results); err != nil {
	//	log.Fatal(err)
	//}
	//qIDs := make([]string, 0)
	//correctOptions := make([][]int, 0)
	//for _, result := range results {
	//	qIDString := result["_id"].(primitive.ObjectID).Hex()
	//	qIDs = append(qIDs, qIDString)
	//	//correctOptions = append(correctOptions, result["correct_options"].(string))
	//	fmt.Println(result["correct_options"])
	//	fmt.Println(reflect.TypeOf(result["correct_options"]))
	//	var intSlice []int
	//	for _, i := range result["correct_options"].(primitive.A) {
	//		intSlice = append(intSlice, int(i.(int32)))
	//	}
	//	correctOptions = append(correctOptions, intSlice)
	//}

	//q.CorrectAnswerList = correctOptions
	q.QIDList = qSetDTO.Questions
	if len(qSetDTO.MarkList) == 0 {
		tMarkList := make([]float64, len(qSetDTO.Questions))
		for i := range tMarkList {
			tMarkList[i] = 1
		}
		q.MarkList = tMarkList
	} else {
		q.MarkList = qSetDTO.MarkList
	}

	insertionResult, err := utils.Mg.Db.Collection("question_set").InsertOne(c.Context(), q)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": err.Error(),
		})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"status": "success", "question_set": q, "id": insertionResult.InsertedID})
}

func GetQuestionSets(c *fiber.Ctx) error {
	findOptions := options.Find()
	filter := bson.D{}

	if subject := c.Query("subject"); subject != "" {
		filter = append(filter, bson.E{Key: "subject", Value: subject})
	}
	if category := c.Query("category"); category != "" {
		filter = append(filter, bson.E{Key: "category", Value: category})
	}
	if exam := c.Query("exam"); exam != "" {
		filter = append(filter, bson.E{Key: "exam", Value: exam})
	}
	if name := c.Query("name"); name != "" {
		filter = append(filter, bson.E{Key: "name", Value: name})
	}
	if tags := c.Query("tags"); tags != "" {
		tagsList := strings.Split(tags, ",")
		filter = append(filter, bson.E{Key: "tags", Value: bson.M{"$in": tagsList}})
	}

	count, _ := strconv.Atoi(c.Query("count", "-1"))
	if count != -1 {
		findOptions.SetLimit(int64(count))
	}

	cursor, err := utils.Mg.Db.Collection("question_set").Find(c.Context(), filter, findOptions)
	if err != nil {
		return handleError(c, fiber.StatusInternalServerError, err.Error())
	}

	var questionSets []models.QuestionSet
	if err := cursor.All(c.Context(), &questionSets); err != nil {
		return handleError(c, fiber.StatusInternalServerError, err.Error())
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"question_sets": questionSets})
}

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
