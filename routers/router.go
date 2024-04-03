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
	questions.Post("/", middlewares.Protected(), controllers.CreateQuestion)
	questions.Get("/", controllers.GetQuestions)

}
