package controllers

import (
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/golang-jwt/jwt/v5"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
	"math/rand"
	"neocognito-backend/models"
	"neocognito-backend/utils"
	"time"
)

func CreateTestSession(c *fiber.Ctx) error {
	type testSessionDTOModel struct {
		QuestionSet        string  `json:"question_set" validate:"required"`
		TestMode           string  `json:"test_mode" `
		RandomizeQuestions bool    `json:"randomize_questions"`
		NegativeMarks      float64 `json:"negative_marks"`
	}
	t := new(testSessionDTOModel)
	if err := c.BodyParser(&t); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Bad Request",
		})
	}
	validate := validator.New()
	err := validate.Struct(t)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": err.Error(),
		})
	}
	idObject, err := primitive.ObjectIDFromHex(t.QuestionSet)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "error", "error": err.Error()})
	}
	questionSetVar := new(models.QuestionSet)
	err = utils.Mg.Db.Collection("question_set").FindOne(c.Context(), bson.M{"_id": idObject}).Decode(&questionSetVar)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "error", "error": err.Error()})
	}
	testSession := new(models.TestSession)
	testSession.CurrentQuestionNum = 0
	testSession.QSetName = questionSetVar.Name

	tempMapVar := make(map[string]*models.QuestionAnswerData)
	for i, question := range questionSetVar.QIDList {
		questionAnswerData := new(models.QuestionAnswerData)
		questionAnswerData.Correct = []int{}
		questionAnswerData.Selected = []int{}
		questionAnswerData.QuestionsTotalMark = questionSetVar.MarkList[i]
		questionAnswerData.QuestionsScoredMark = 0
		questionAnswerData.Answered = false
		tempMapVar[question] = questionAnswerData
	}
	testSession.NegativeMarks = t.NegativeMarks
	testSession.QuestionAnswerData = tempMapVar
	if t.RandomizeQuestions {
		testSession.QuestionIDsOrdered = shuffleQuestionIDs(&questionSetVar.QIDList)
	} else {
		testSession.QuestionIDsOrdered = questionSetVar.QIDList
	}
	currentTime := time.Now()
	testSession.Mode = t.TestMode
	testSession.Finished = false
	user := c.Locals("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	testSession.TakenByID = claims["id"].(string)
	testSession.StartedTime = &currentTime
	testSession.QuestionSetID = questionSetVar.ID
	insertResult, err := utils.Mg.Db.Collection("test_session").InsertOne(c.Context(), testSession)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "error", "error": err.Error()})
	}
	testSessionID := insertResult.InsertedID.(primitive.ObjectID)
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"status":          "success",
		"test_session_id": testSessionID,
		"test_session":    testSession,
	})
}

func UpdateTestSession(c *fiber.Ctx) error {
	type answerDTO struct {
		QuestionAnswerData   map[string]interface{} `json:"question_answer_data"`
		CurrentQuestionIndex int                    `json:"current_question_index"`
		TotalMarksScored     float64                `json:"total_marks_scored"`
	}
	testSessionID := c.Params("test_session_id")
	if testSessionID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Test session ID is required",
		})
	}
	dto := new(answerDTO)
	if err := c.BodyParser(&dto); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Invalid input",
		})
	}
	validate := validator.New()
	if err := validate.Struct(dto); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": err.Error(),
		})
	}
	idObject, err := primitive.ObjectIDFromHex(testSessionID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status": "error",
			"error":  "Invalid test session ID",
		})
	}
	var testSession models.TestSession
	err = utils.Mg.Db.Collection("test_session").FindOne(c.Context(), bson.M{"_id": idObject}).Decode(&testSession)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"status":  "error",
			"message": "Test session not found",
		})
	}
	if testSession.Finished {
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"status":  "finished",
			"message": "Test session is already finished",
		})
	}
	user := c.Locals("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	userID := claims["id"].(string)

	if testSession.TakenByID != userID {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"status":  "error",
			"message": "Unauthorized",
		})
	}
	fullQuestionAnswerDataObj := make(map[string]*models.QuestionAnswerData)
	for key, value := range dto.QuestionAnswerData {
		qObj := new(models.QuestionAnswerData)
		mapOfValue := value.(map[string]interface{})

		for k, v := range mapOfValue {
			if k == "answered" {
				qObj.Answered = v.(bool)
			} else if k == "selected_answer_list" {
				selectedList := v.([]interface{})
				qObj.Selected = make([]int, len(selectedList))
				for i, val := range selectedList {
					qObj.Selected[i] = int(val.(float64)) // JSON numbers are float64
				}
			} else if k == "correct_answer_list" {
				correctList := v.([]interface{})
				qObj.Correct = make([]int, len(correctList))
				for i, val := range correctList {
					qObj.Correct[i] = int(val.(float64)) // JSON numbers are float64
				}
			} else if k == "questions_total_mark" {
				qObj.QuestionsTotalMark = v.(float64)
			} else if k == "questions_scored_mark" {
				qObj.QuestionsScoredMark = v.(float64)
			}
		}

		fullQuestionAnswerDataObj[key] = qObj
	}

	testSession.QuestionAnswerData = fullQuestionAnswerDataObj
	testSession.CurrentQuestionNum = dto.CurrentQuestionIndex
	testSession.ScoredMarks = dto.TotalMarksScored

	updateObject := bson.M{
		"$set": bson.M{
			"question_answer_data": testSession.QuestionAnswerData,
			"scored_marks":         testSession.ScoredMarks,
			"current_question_num": testSession.CurrentQuestionNum,
		},
	}

	_, err = utils.Mg.Db.Collection("test_session").UpdateOne(
		c.Context(),
		bson.M{"_id": idObject},
		updateObject,
	)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status": "error",
			"error":  "Failed to update the test session" + err.Error(),
		})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":       "success",
		"test_session": testSession,
	})
}

func GetTestSession(c *fiber.Ctx) error {
	testSessionID := c.Params("test_session_id")
	if testSessionID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Test session ID is required",
		})
	}

	idObject, err := primitive.ObjectIDFromHex(testSessionID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Invalid test session ID",
		})
	}

	var testSession models.TestSession
	err = utils.Mg.Db.Collection("test_session").FindOne(c.Context(), bson.M{"_id": idObject}).Decode(&testSession)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"status":  "error",
			"message": "Test session not found",
		})
	}

	user := c.Locals("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	userID := claims["id"].(string)

	if testSession.TakenByID != userID {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"status":  "error",
			"message": "Unauthorized",
		})
	}
	currentQuestionIndex := testSession.CurrentQuestionNum
	if currentQuestionIndex < len(testSession.QuestionIDsOrdered) {
		currentQuestionID := testSession.QuestionIDsOrdered[currentQuestionIndex]
		questionIDObject, err := primitive.ObjectIDFromHex(currentQuestionID)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"status":  "error",
				"message": "Invalid test session ID",
			})
		}

		var question models.Question
		err = utils.Mg.Db.Collection("questions").FindOne(c.Context(), bson.M{"_id": questionIDObject}).Decode(&question)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"status":  "error",
				"message": "Failed to fetch the current question",
			})
		}
	}
	var objectIDs []primitive.ObjectID
	for _, idStr := range testSession.QuestionIDsOrdered {
		id, err := primitive.ObjectIDFromHex(idStr)
		if err != nil {
			continue
		}
		objectIDs = append(objectIDs, id)
	}
	projection := bson.M{"id": 1, "correct_options": 1, "question": 1, "question_type": 1, "options": 1, "explanation": 1}
	findOptions := options.Find()
	findOptions.SetProjection(projection)

	filter := bson.M{"_id": bson.M{"$in": objectIDs}}
	cursor, err := utils.Mg.Db.Collection("questions").Find(c.Context(), filter, findOptions)
	var results []bson.M
	if err := cursor.All(c.Context(), &results); err != nil {
		log.Fatal(err)
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":                 "success",
		"test_session":           testSession,
		"current_question_index": testSession.CurrentQuestionNum,
		"current_question_id":    testSession.QuestionIDsOrdered[testSession.CurrentQuestionNum],
		"questions":              results,
	})
}

func FinishTestSession(c *fiber.Ctx) error {
	fmt.Println("test finish called")
	testSessionID := c.Params("test_session_id")
	if testSessionID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Test session ID is required",
		})
	}

	idObject, err := primitive.ObjectIDFromHex(testSessionID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Invalid test session ID",
		})
	}

	var testSession models.TestSession
	err = utils.Mg.Db.Collection("test_session").FindOne(c.Context(), bson.M{"_id": idObject}).Decode(&testSession)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"status":  "error",
			"message": "Test session not found",
		})
	}

	if testSession.Finished {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Test session is already finished",
		})
	}

	user := c.Locals("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	userID := claims["id"].(string)

	if testSession.TakenByID != userID {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"status":  "error",
			"message": "Unauthorized",
		})
	}

	var totalMarks, scoredMarks float64

	for _, quesAnsData := range testSession.QuestionAnswerData {
		totalMarks += quesAnsData.QuestionsTotalMark
		scoredMarks += quesAnsData.QuestionsScoredMark
	}

	testSession.FinishedTime = timePtr(time.Now())
	testSession.Finished = true
	testSession.TotalMarks = totalMarks
	testSession.ScoredMarks = scoredMarks

	_, err = utils.Mg.Db.Collection("test_session").UpdateOne(
		c.Context(),
		bson.M{"_id": idObject},
		bson.M{
			"$set": bson.M{
				"finished":      testSession.Finished,
				"finished_time": testSession.FinishedTime,
				"total_marks":   testSession.TotalMarks,
				"scored_marks":  testSession.ScoredMarks,
			},
		},
	)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status": "error",
			"error":  "Failed to finish the test session: " + err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":          "success",
		"started_time":    testSession.StartedTime,
		"finished_time":   testSession.FinishedTime,
		"total_marks":     testSession.TotalMarks,
		"scored_marks":    testSession.ScoredMarks,
		"test_session_id": testSession.ID,
		"finished_status": testSession.Finished,
	})
}

func timePtr(t time.Time) *time.Time {
	return &t
}

func CalculateScoredMarks(correctOptions, selectedOptions []int, maxMarks float64) float64 {
	correctOptionsMap := make(map[int]bool)
	for _, option := range correctOptions {
		correctOptionsMap[option] = true
	}
	for _, selected := range selectedOptions {
		if !correctOptionsMap[selected] {
			return 0
		}
	}
	numCorrectOptions := len(correctOptions)
	numSelectedOptions := len(selectedOptions)
	fraction := float64(numSelectedOptions) / float64(numCorrectOptions)
	return fraction * maxMarks
}

func shuffleQuestionIDs(input *[]string) []string {
	rand.Seed(time.Now().UnixNano())

	for i := range *input {
		j := rand.Intn(i + 1)
		(*input)[i], (*input)[j] = (*input)[j], (*input)[i]
	}
	return *input
}
func GetTestSessionByUserID(c *fiber.Ctx) error {
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status": "success",
	})
}
