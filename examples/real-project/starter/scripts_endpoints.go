package starter

import (
	"encoding/json"
	"fmt"
	"os"
	"real-project/database"
	"real-project/util"

	"github.com/Liphium/magic/mrunner"
)

// This method would ideally be created in a shared package between all the scripts.
func GetPath() string {
	return fmt.Sprintf("http://%s", os.Getenv("LISTEN"))
}

// Script for creating a post using the endpoint.
//
// You could go into the database and add it there, but we want to be able to call the endpoint using scripts.
func createPost(_ *mrunner.Runner, post database.Post) error {
	_, err := util.Post[interface{}](GetPath()+"/posts", post, util.Headers{})
	return err
}

// Script for printing all the posts using the endpoint.
func printPosts(_ *mrunner.Runner, _ any) error {
	posts, err := util.Get[[]database.Post](GetPath()+"/posts", util.Headers{})
	if err != nil {
		return err
	}

	better, _ := json.MarshalIndent(posts, "", "   ")
	fmt.Println(string(better))
	return nil
}
