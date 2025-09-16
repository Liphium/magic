package starter

import (
	"fmt"
	"log"
	"real-project/database"

	"github.com/Liphium/magic/mrunner"
)

// Script to reset the database by dropping and recreating all tables
//
// Here we just use any to ignore the argument. This can be useful for scripts such as this one.
func resetDatabase(runner *mrunner.Runner, _ any) error {
	log.Println("Resetting database...")

	// Connect to the database
	database.Connect()

	// Drop all tables and recreate them
	if err := database.DBConn.Migrator().DropTable(&database.Post{}); err != nil {
		log.Printf("Warning: Could not drop Post table: %v", err)
	}

	// Recreate tables
	if err := database.DBConn.AutoMigrate(&database.Post{}); err != nil {
		return fmt.Errorf("failed to migrate database: %v", err)
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
func SeedDatabase(runner *mrunner.Runner, _ any) error {
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
