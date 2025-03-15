package main

import (
	"log"
	"os"

	forge_service "github.com/Liphium/magic/spellcast/forge"
	"github.com/Liphium/magic/spellcast/util"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/joho/godotenv"
)

const (
	ForgeServiceTag = "forge"
)

func main() {

	// Get all the stuff from the start arguments
	args := os.Args[1:]

	// Get the stuff needed for execution
	service := args[0]
	token := args[1]

	// Setup environment variables
	if err := godotenv.Load(); err != nil {
		log.Fatalln("Couldn't load env file:", err)
	}

	// Set all necessary stuff
	util.BackendURL = os.Getenv("SC_BACKEND")
	if util.BackendURL == "" {
		log.Fatalln("Please specify the backend URL using the SC_BACKEND environment variable (example: http://localhost:8081)")
	}
	util.InitWorkingDirectory()

	// Set up docker
	util.InitDocker()

	// Setup all the endpoints
	app := setupApp(service, token)

	app.Listen("127.0.0.1:9000")
}

func setupApp(service string, token string) *fiber.App {
	app := fiber.New()

	// Setup middlewares to make life a little easier
	app.Use(logger.New())
	app.Use(cors.New())

	// Setup the correct service
	switch service {
	case ForgeServiceTag:
		app.Route("/", forge_service.SetupRoutes(token))
	}

	return app
}
