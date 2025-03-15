package forge_service

import (
	"log"
	"os/exec"
	"strings"

	"github.com/Liphium/magic/spellcast/util"
	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(token string) func(fiber.Router) {

	// Connect to the backend to get all required information
	body, err := util.PostRequestBackendToken("/api/spellcast/forge/connect", token, map[string]interface{}{})
	if err != nil {
		log.Fatalln("couldn't connect to backend")
	}

	// Make sure the request was successful
	if !body["success"].(bool) {
		log.Fatalln("the backend didn't return success")
	}

	// Clone the repository
	args := strings.Split(body["clone_command"].(string), " ")[1:]
	if err := exec.Command("git", args...).Run(); err != nil {
		log.Fatalln("couldn't ")
	}

	log.Println("Repository cloned successfully.")

	return func(r fiber.Router) {

	}
}
