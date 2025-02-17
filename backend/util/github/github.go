package github_utils

import (
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/Liphium/magic/backend/database"
	"github.com/Liphium/magic/backend/util/constants"
	"github.com/bradleyfalzon/ghinstallation/v2"
	"github.com/gofiber/fiber/v2"
	"github.com/google/go-github/v69/github"
)

// Id of the Magic Github app
var appId int64

func getAppId() int64 {
	if appId == 0 {
		// Initialize the Github app id for Magic
		appStr := os.Getenv("MAGIC_GH_APP")
		log.Println(appStr)
		var err error
		appId, err = strconv.ParseInt(appStr, 10, 64)
		if err != nil {
			log.Fatalln("Couldn't parse GitHub app id to int:", err)
		}
	}

	return appId
}

// Get a GitHub client for the user from context (using database and auth information)
func GetUserFromContext(c *fiber.Ctx) (*github.Client, error) {

	// Get the user access token for Github
	acc := c.Locals(constants.LocalsAccountID)
	var credential database.Credential
	if err := database.DBConn.Where("account = ?", acc).Take(&credential).Error; err != nil {
		return nil, err
	}

	// Get all installations from Github
	return github.NewClient(nil).WithAuthToken(credential.Token), nil
}

// Get a GitHub client for an installation using the private key.
func GetInstallationClient(installationId int64) (*github.Client, error) {

	// Get a new http client that's authenticated for the Github installation
	itr, err := ghinstallation.NewKeyFromFile(http.DefaultTransport, getAppId(), installationId, "github.pem")
	if err != nil {
		return nil, err
	}

	return github.NewClient(&http.Client{Transport: itr}), nil
}
