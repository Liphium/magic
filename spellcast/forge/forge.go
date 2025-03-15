package forge_service

import (
	"context"
	"log"
	"os/exec"
	"slices"
	"strings"
	"sync"

	"github.com/Liphium/magic/spellcast/util"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/pkg/archive"
	"github.com/docker/docker/pkg/jsonmessage"
	"github.com/gofiber/fiber/v2"
	"github.com/moby/term"
)

type LogStorage struct {
	mutex *sync.Mutex
	logs  []string
	subs  []chan string
}

func (lg *LogStorage) Write(p []byte) (n int, err error) {
	lg.mutex.Lock()
	defer lg.mutex.Unlock()

	// Turn into a string to sent through the channels
	str := string(p)
	for _, ch := range lg.subs {
		ch <- str
	}

	// Add to the logs
	lg.logs = append(lg.logs, str)

	return len(str), nil
}

var currentLogStorage *LogStorage

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

	// Update the current log storage
	currentLogStorage = &LogStorage{
		mutex: &sync.Mutex{},
		logs:  []string{},
		subs:  []chan string{},
	}

	// Start a new goroutine to clone and build the repository
	go startBuild(body)

	return func(r fiber.Router) {

		r.Post("/logs", func(c *fiber.Ctx) error {
			return util.StartEvents(c, func(w *util.EventWriter) {

				// Copy current events
				currentLogStorage.mutex.Lock()
				logCopy := slices.Clone(currentLogStorage.logs)
				currentLogStorage.mutex.Unlock()

				// Add a new subscription
				currentLogStorage.mutex.Lock()
				logChan := make(chan string)
				currentLogStorage.subs = append(currentLogStorage.subs, logChan)
				i := len(currentLogStorage.subs) - 1
				currentLogStorage.mutex.Unlock()

				// Cancel the subscription when the function returns
				defer func() {
					recover()
					currentLogStorage.mutex.Lock()
					currentLogStorage.subs = slices.Delete(currentLogStorage.subs, i, i+1)
					currentLogStorage.mutex.Unlock()
				}()

				// Send all of the current events
				end := false
				for _, currentLog := range logCopy {
					if err := w.SendEvent("", currentLog); err != nil {
						log.Println("Couldn't send event:", err)
						end = true
						break
					}
				}
				if end {
					return
				}

				// Start sending all events
				for {
					if err := w.SendEvent("", <-logChan); err != nil {
						log.Println("Couldn't send event:", err)
						break
					}
				}

			})
		})

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
	termFd, _ := term.GetFdInfo(currentLogStorage)
	jsonmessage.DisplayJSONMessagesStream(resp.Body, currentLogStorage, termFd, false, nil)

	log.Println("repository built")
}
