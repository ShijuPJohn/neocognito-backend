package controllers

import (
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
	//user := c.Locals("user").(*jwt.Token)
	//claims := user.Claims.(jwt.MapClaims)
	//userID := claims["id"].(string)
	testSession := new(models.TestSession)
	testSession.CurrentQuestionNum = 0
	queAnsData := make(map[string]map[string][]int)

	for i, question := range questionSetVar.QIDList {
		queAnsData[question] = map[string][]int{"correct": questionSetVar.CorrectAnswerList[i], "selected": []int{}}
	}
	testSession.QuesCaSa = queAnsData
	if t.RandomizeQuestions {
		testSession.QuestionIDsOrdered = shuffleQuestionIDs(&questionSetVar.QIDList)
	} else {
		testSession.QuestionIDsOrdered = questionSetVar.QIDList
	}
	currentTime := time.Now()
	testSession.Mode = t.TestMode
	testSession.Started = false
	testSession.Finished = false
	user := c.Locals("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	testSession.TakenByID = claims["id"].(string)
	testSession.StartedTime = &currentTime
	testSession.QuestionSetID = questionSetVar.ID
	_, err = utils.Mg.Db.Collection("test_session").InsertOne(c.Context(), testSession)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "error", "error": err.Error()})
	}
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"status":        "success",
		"test_sesssion": testSession,
	})
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
func shuffleQuestionIDs(input *[]string) []string {
	rand.Seed(time.Now().UnixNano()) // Seed the random number generator

	for i := range *input {
		j := rand.Intn(i + 1)
		(*input)[i], (*input)[j] = (*input)[j], (*input)[i]
	}
	return *input
}
