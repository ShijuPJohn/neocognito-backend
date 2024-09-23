package controllers

import (
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"math/rand"
	"neocognito-backend/models"
	"neocognito-backend/utils"
	"time"
)

func CreateTestSession(c *fiber.Ctx) error {
	type testSessionDTOModel struct {
		QuestionSet        string `json:"question_set" validate:"required"`
		TestMode           string `json:"test_mode" `
		RandomizeQuestions bool   `json:"randomize_questions"`
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
	testSession := new(models.TestSession)
	testSession.CurrentQuestionNum = 0
	tempMapVar := make(map[string]*models.QuestionAnswerData)
	fmt.Println(len(questionSetVar.QIDList))
	fmt.Println(len(questionSetVar.CorrectAnswerList))
	fmt.Println(len(questionSetVar.MarkList))
	for i, question := range questionSetVar.QIDList {
		questionAnswerData := new(models.QuestionAnswerData)
		questionAnswerData.Correct = questionSetVar.CorrectAnswerList[i]
		questionAnswerData.Selected = make([]int, 0)
		questionAnswerData.QuestionsTotalMark = questionSetVar.MarkList[i]
		questionAnswerData.QuestionsScoredMark = 0
		tempMapVar[question] = questionAnswerData
	}
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
		QuestionID     string `json:"question_id" validate:"required"`
		SelectedAnswer []int  `json:"selected_answer" validate:"required"`
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
	user := c.Locals("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	userID := claims["id"].(string)

	if testSession.TakenByID != userID {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"status":  "error",
			"message": "Unauthorized",
		})
	}
	if testSession.Finished {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Test session is already finished",
		})
	}

	quesAnsData, exists := testSession.QuestionAnswerData[dto.QuestionID]
	if !exists {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Question not found in the test session",
		})
	}

	quesAnsData.Selected = dto.SelectedAnswer
	testSession.QuestionAnswerData[dto.QuestionID] = quesAnsData
	if len(dto.SelectedAnswer) == 0 {
		quesAnsData.QuestionsScoredMark = 0
	} else {
		quesAnsData.QuestionsScoredMark = CalculateScoredMarks(quesAnsData.Correct, quesAnsData.Selected, quesAnsData.QuestionsTotalMark)
	}
	testSession.CurrentQuestionNum++
	testSession.QuestionAnswerData[dto.QuestionID] = quesAnsData

	_, err = utils.Mg.Db.Collection("test_session").UpdateOne(
		c.Context(),
		bson.M{"_id": idObject},
		bson.M{
			"$set": bson.M{
				"question_answer_data": testSession.QuestionAnswerData,
				"currentQuestionNum":   testSession.CurrentQuestionNum,
			},
		},
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
	if testSession.Finished {
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"status":       "test finished",
			"test_session": testSession,
		})
	}
	currentQuestionIndex := testSession.CurrentQuestionNum

	var currentQuestion map[string]interface{}
	if currentQuestionIndex < len(testSession.QuestionIDsOrdered) {
		currentQuestionID := testSession.QuestionIDsOrdered[currentQuestionIndex]
		fmt.Println(currentQuestionID)
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

		currentQuestion = fiber.Map{
			"id":            question.ID,
			"question":      question.Question,
			"options":       question.Options,
			"question_type": question.QuestionType,
			"difficulty":    question.Difficulty,
		}
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":           "success",
		"test_session":     testSession,
		"current_question": currentQuestion,
		//"question_ids_ordered":   testSession.QuestionIDsOrdered,
		//"answered_data":          testSession.QuestionAnswerData,
		//"current_question_index": currentQuestionIndex,
	})
}

func FinishTestSession(c *fiber.Ctx) error {
	// Get the test session ID from the path param
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
	fmt.Println(totalMarks, scoredMarks)

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

//	func ResumeQTest(c *fiber.Ctx) error {
//		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "error", "error": ""})
//	}
//
//	func GetQTestByID(c *fiber.Ctx) error {
//		qTestId := c.Params("id")
//		idObject, err := primitive.ObjectIDFromHex(qTestId)
//		if err != nil {
//			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "error", "error": err.Error()})
//		}
//		qTest := new(models.QTest)
//		err = utils.Mg.Db.Collection("q_test").FindOne(c.Context(), bson.M{"_id": idObject}).Decode(&qTest)
//		if err != nil {
//			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "error", "error": err.Error()})
//		}
//		return c.Status(fiber.StatusOK).JSON(fiber.Map{"status": "success", "qTest": qTest})
//	}
//
//	func TakeTest(c *fiber.Ctx) error {
//		qTestId := c.Params("id")
//		type Temp struct {
//			Question string `json:"question" validate:"required"`
//			Answer   string `json:"answer" validate:"required"`
//		}
//
//		t := new(Temp)
//		if err := c.BodyParser(&t); err != nil {
//			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
//				"status":  "error",
//				"message": "Bad Request",
//				"error":   err.Error(),
//			})
//		}
//		selectedAnswer, err := strconv.Atoi(t.Answer)
//		idObject, err := primitive.ObjectIDFromHex(qTestId)
//		if err != nil {
//			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "error", "error": err.Error()})
//		}
//		qTest := new(models.QTest)
//		qIdObject, err := primitive.ObjectIDFromHex(t.Question)
//
//		if err != nil {
//			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "error", "error": err.Error()})
//		}
//		question := new(models.Question)
//		err = utils.Mg.Db.Collection("q_test").FindOne(c.Context(), bson.M{"_id": idObject}).Decode(&qTest)
//		if err != nil {
//			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "error", "error": err.Error()})
//		}
//		if qTest.Finished {
//			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "error", "error": "Test already finished"})
//		}
//		err = utils.Mg.Db.Collection("questions").FindOne(c.Context(), bson.M{"_id": qIdObject}).Decode(&question)
//		if err != nil {
//			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "error", "error": err.Error()})
//		}
//		var answerSlice []int
//		answerSlice = append(answerSlice, selectedAnswer)
//		answerSlice = append(answerSlice, question.CorrectOptions)
//		user := c.Locals("user").(*jwt.Token)
//		claims := user.Claims.(jwt.MapClaims)
//		if qTest.TakenById != claims["id"].(string) {
//			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"status": "unauthorized"})
//		}
//
//		qTest.AllQuestionsIDs[t.Question] = answerSlice
//
//		update := bson.M{"allQuestionsId": qTest.AllQuestionsIDs}
//		filter := bson.M{"_id": idObject}
//		if selectedAnswer == question.CorrectOptions {
//			update["nCorrectlyAnswered"] = qTest.NCorrectlyAnswered + 1
//		}
//		updateQuery := bson.M{"$set": update}
//		result, err := utils.Mg.Db.Collection("q_test").UpdateOne(c.Context(), filter, updateQuery)
//		if err != nil {
//			return handleError(c, fiber.StatusBadRequest, err.Error())
//		}
//		if qTest.Mode != "exam" {
//			totalScoreSoFar := 0
//			for question := range qTest.AllQuestionsIDs {
//				if qTest.AllQuestionsIDs[question][0] == qTest.AllQuestionsIDs[question][1] && qTest.AllQuestionsIDs[question][0] != 0 {
//					totalScoreSoFar++
//				}
//			}
//			return c.Status(fiber.StatusOK).JSON(fiber.Map{"status": "success", "result": result, "totalScoreSoFar": totalScoreSoFar})
//
//		}
//
//		return c.Status(fiber.StatusOK).JSON(fiber.Map{"status": "success", "result": result})
//
// }
//
// //	func finishQTest(c *fiber.Ctx) error {
// //		qTestId := c.Params("id")
// //		qTest := new(models.QTest)
// //		qTestIdObject,err := primitive.ObjectIDFromHex(qTestId)
// //		if err != nil {
// //			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "error", "error": err.Error()})
// //		}
// //	}
