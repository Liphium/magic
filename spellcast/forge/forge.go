package forge_service

import (
	"context"
	"log"
	"os"
	"os/exec"
	"strings"

	"github.com/Liphium/magic/spellcast/util"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/pkg/archive"
	"github.com/docker/docker/pkg/jsonmessage"
	"github.com/gofiber/fiber/v2"
	"github.com/moby/term"
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

	// Start a new goroutine to clone and build the repository
	go startBuild(body)

	return func(r fiber.Router) {

	}
}

// Start the docker image build
func startBuild(body map[string]interface{}) {

	// Delete everything in the working directory
	if err := util.DeleteAllFiles(util.WorkingDirectory); err != nil {
		log.Fatalln("Couldn't delete files in work dir:", err)
	}

	// Clone the repository
	args := strings.Split(body["clone_command"].(string), " ")[1:]
	if err := exec.Command("git", args...).Run(); err != nil {
		log.Fatalln("Couldn't clone repository:", err)
	}

	// Get the build context for the image build
	buildCtx, _ := archive.TarWithOptions(util.WorkingDirectory, &archive.TarOptions{})

	// Build the image
	resp, err := util.DockerClient.ImageBuild(context.Background(), buildCtx, types.ImageBuildOptions{
		Tags:       []string{"built-image"},
		Dockerfile: "Dockerfile",
	})
	if err != nil {
		log.Fatalln("Couldn't build Docker image:", err)
	}

	// Print messages to os.Stdout
	termFd, _ := term.GetFdInfo(os.Stdout)
	jsonmessage.DisplayJSONMessagesStream(resp.Body, os.Stdout, termFd, false, nil)

	log.Println("repository built")
}
