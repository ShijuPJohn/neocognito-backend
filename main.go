package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"log"
	"neocognito-backend/routers"
	"neocognito-backend/utils"
	"os"
)

func main() {
	err, deferFunc := utils.MongoDBConnect()
	defer deferFunc()
	if err != nil {
		log.Fatal(err)
	}
	app := fiber.New()
	app.Use(cors.New())
	app.Use(logger.New())
	routers.SetupRoutes(app)
	var port string
	if port = os.Getenv("PORT"); port == "" {
		port = "8080"
	}
	log.Fatal(app.Listen(":" + port))
}
