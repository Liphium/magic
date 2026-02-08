package starter

import (
	"fmt"
	"log"
	"real-project/database"

	"github.com/Liphium/magic/v2/mrunner"
)

// Script to clear all database tables content, but not fully delete them.
//
// Here we just use any to ignore the argument. This can be useful for scripts such as this one.
func ClearDatabases(runner *mrunner.Runner) error {
	log.Println("Resetting database...")

	// Magic can clear all databases for you, don't worry, only data will be deleted meaning your schema is still all good :D
	if err := runner.ClearTables(); err != nil {
		log.Fatalln("Couldn't clear database tables:", err)
	}

	log.Println("Database reset completed successfully!")
	return nil
}

// Script to reset the database by dropping all tables.
//
// Here we just use any to ignore the argument. This can be useful for scripts such as this one.
func ResetDatabase(runner *mrunner.Runner) error {
	log.Println("Resetting database...")

	// Magic can drop all databases for you as well, this means that all the tables are actually gone
	if err := runner.DropTables(); err != nil {
		log.Fatalln("Couldn't clear database tables:", err)
	}

	log.Println("Database reset completed successfully!")
	return nil
}

var SamplePosts = []database.Post{
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

	// Connect to the database
	database.Connect()

	// Insert sample posts
	for _, post := range SamplePosts {
		if err := database.DBConn.Create(&post).Error; err != nil {
			return fmt.Errorf("failed to create sample post: %v", err)
		}
	}

	log.Printf("Successfully seeded database with %d sample posts!", len(SamplePosts))
	return nil
}
