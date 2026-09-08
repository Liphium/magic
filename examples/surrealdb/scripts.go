package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/Liphium/magic/v4/mconfig"
	"github.com/Liphium/magic/v4/mrunner"
	"github.com/gofiber/fiber/v2"
	"resty.dev/v3"
)

// This method would ideally be created in a shared package between all the scripts.
func GetPath() string {
	return fmt.Sprintf("http://%s", os.Getenv("LISTEN"))
}

// Script for creating a post using the endpoint.
//
// You could go into the database and add it there, but we want to be able to call the endpoint using scripts.
func CreatePost(post Post) error {
	client := resty.New()
	defer client.Close()

	_, err := client.R().
		SetBody(post).
		Post(GetPath() + "/posts")
	return err
}

// Script for printing all the posts using the endpoint.
func PrintPosts() error {
	client := resty.New()
	defer client.Close()

	res, err := client.R().
		Get(GetPath() + "/posts")
	if err != nil || res.StatusCode() != fiber.StatusOK {
		return fmt.Errorf("couldn't get posts: %v", err)
	}

	var posts []Post
	if err := json.Unmarshal(res.Bytes(), &posts); err != nil {
		return err
	}

	better, _ := json.MarshalIndent(posts, "", "   ")
	fmt.Println(string(better))
	return nil
}

// Script to clear all database tables content, but not fully delete them.
//
// Here we just use any to ignore the argument. This can be useful for scripts such as this one.
func ClearDatabases(runner *mrunner.Runner) error {
	log.Println("Clearing database...")

	// Magic can clear all databases for you, don't worry, only data will be deleted meaning your schema is still all good :D
	if err := runner.RunInstruction(mconfig.InstructionClearTables); err != nil {
		return fmt.Errorf("couldn't clear database tables: %w", err)
	}

	log.Println("Database clear completed successfully!")
	return nil
}

// Script to reset the database by dropping all tables.
//
// Here we just use any to ignore the argument. This can be useful for scripts such as this one.
func ResetDatabase(runner *mrunner.Runner) error {
	log.Println("Resetting database...")

	// Magic can drop all databases for you as well, this means that all the tables are actually gone
	if err := runner.RunInstruction(mconfig.InstructionDropTables); err != nil {
		log.Fatalln("Couldn't reset database tables:", err)
	}

	log.Println("Database reset completed successfully!")
	return nil
}

var SamplePosts = []Post{
	{Author: "Alice", Content: "Welcome to our new blog platform! This is the first post."},
	{Author: "Bob", Content: "I love how easy it is to create posts here. Great work!"},
	{Author: "Charlie", Content: "Looking forward to sharing more content with everyone."},
	{Author: "Diana", Content: "The API is so clean and well-designed. Kudos to the developers!"},
}

// Script to seed the database with sample posts
//
// Here we just use any to ignore the argument. This can be useful for scripts such as this one.
func SeedDatabase() error {
	log.Println("Seeding database with sample posts...")

	// Connect to the database (in case this is run as a script without the app running)
	ConnectSurreal()

	// Insert sample posts
	for _, post := range SamplePosts {
		if _, err := createPost(post); err != nil {
			return fmt.Errorf("failed to create sample post: %v", err)
		}
	}

	log.Printf("Successfully seeded database with %d sample posts!", len(SamplePosts))
	return nil
}
