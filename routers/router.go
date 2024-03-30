package routers

import (
	"github.com/gofiber/fiber/v2"
	"neocognito-backend/controllers"
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
	//auth.Post("/fbLogin", handlers.FBLoginHandler)
	//auth.Post("/googleLogin", handlers.GoogleLoginHandler)
	//auth.Post("/email_verify", middlewares.Protected(), handlers.EmailVerifyHandler)
	//
	//// Deals
	//deals := api.Group("/deals")
	//deals.Post("/", middlewares.Protected(), handlers.PostDeals)
	//deals.Get("/", handlers.GetDeals)
	//deals.Get("/editDeals", middlewares.Protected(), handlers.GetDealsForEdit)
	//deals.Get("/:shortId", handlers.GetDealByShortId)
	//deals.Delete("/:shortId", middlewares.Protected(), handlers.DeleteByShortID)
	//deals.Put("/:shortId", middlewares.Protected(), handlers.EditHandler)
	//
	//// Messages
	//messages := api.Group("/messages")
	////messages.Get("/", handlers.GetMessagesHandler)
	//messages.Post("/", handlers.PostMessageHandler)
	////Demo

}
