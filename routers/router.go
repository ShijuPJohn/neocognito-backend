package routers

import (
	"github.com/gofiber/fiber/v2"
	"neocognito-backend/controllers"
	"neocognito-backend/middlewares"
)

func SetupRoutes(app *fiber.App) {

	//Index
	api := app.Group("/api")
	api.Get("/", controllers.Index)
	//api.Get("/home", middlewares.Protected(), handlers.Home)

	//// Auth
	auth := api.Group("/auth")
	auth.Get("/users", controllers.GetAllUsers)
	auth.Post("/signup", controllers.CreateUser)
	auth.Post("/login", controllers.LoginUser)
	auth.Get("/users/:id", middlewares.Protected(), controllers.GetUserDetails)
	//
	//// Questions
	questions := api.Group("/questions")
	questions.Get("/", middlewares.Protected(), controllers.GetQuestions)
	questions.Post("/", middlewares.Protected(), controllers.CreateQuestion)
	questions.Patch("/:id", middlewares.Protected(), controllers.EditQuestion)
	questions.Get("/:id", middlewares.Protected(), controllers.GetQuestionByID)
	questions.Delete("/:id", middlewares.Protected(), controllers.DeleteQuestion)

	questionSet := api.Group("/questionsets")
	questionSet.Post("/", middlewares.Protected(), controllers.CreateQuestionSet)
	questionSet.Get("/", middlewares.Protected(), controllers.GetQuestionSets)
	//questionSet.Patch("/:id", middlewares.Protected(), controllers.EditQuestionSet)
	//questionSet.Get("/:id", controllers.GetQuestionSetByID)
	//questionSet.Delete("/:id", middlewares.Protected(), controllers.DeleteQuestionSet)
	//
	testSession := api.Group("/test_session")
	testSession.Post("/", middlewares.Protected(), controllers.CreateTestSession)
	testSession.Put("/:test_session_id", middlewares.Protected(), controllers.UpdateTestSession)
	testSession.Get("/:test_session_id", middlewares.Protected(), controllers.GetTestSession)
	testSession.Put("/finish/:test_session_id", middlewares.Protected(), controllers.FinishTestSession)
}
