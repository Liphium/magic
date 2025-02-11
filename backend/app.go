package main

import (
	"fmt"
	"log"
	"os"

	"github.com/Liphium/magic/backend/database"
	"github.com/Liphium/magic/backend/routes"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/joho/godotenv"
)

func main() {

	// Load all environment variables
	if err := godotenv.Load(); err != nil {
		panic(err)
	}

	// Load secret environment variables (if there are any)
	godotenv.Load("secret.env")

	// Make sure the environment variables are set correctly
	port := os.Getenv("MAGIC_PORT")
	if port == "" {
		log.Fatal("MAGIC_PORT env variable not set correctly")
	}
	listen := os.Getenv("MAGIC_LISTEN")
	if listen == "" {
		listen = "0.0.0.0"
		log.Println("Listening on 0.0.0.0, you can specify something different by using the MAGIC_LISTEN environment variable..")
	}

	// Connect to the database and initialize
	database.Connect()

	// Start fiber
	app := setupApp()
	if err := app.Listen(fmt.Sprintf("%s:%s", listen, port)); err != nil {
		log.Fatal(err)
	}
}

func setupApp() *fiber.App {
	app := fiber.New()
	app.Use(logger.New())
	app.Use(cors.New())

	// Initialize the router
	app.Route("/", routes.InitializeRoutes)

	return app
}
