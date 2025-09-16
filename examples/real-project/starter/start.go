package starter

import (
	"fmt"
	"os"
	"real-project/database"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func Start() {

	// Connect to the database
	database.Connect()

	// Create the actual web app
	app := fiber.New()

	// This message is just here for explanation.
	fmt.Println()
	fmt.Println("Welcome, wizard! That's how easy it is to get Magic up and running. Well, I hope it actually all went well for you...")
	fmt.Println("Anyway, now that we're up and running, you can open another terminal and run scripts using: 'go run . -r <script>' or list of all of them using 'go run . --scripts'!")
	fmt.Println("Thanks for using Magic!")

	// Basic insertion endpoint to create a new post
	app.Post("/posts", func(c *fiber.Ctx) error {
		var post database.Post

		// Parse the JSON body
		if err := c.BodyParser(&post); err != nil {
			return c.Status(400).JSON(fiber.Map{
				"error": "Invalid JSON format",
			})
		}

		// Validate required fields
		if post.Author == "" || post.Content == "" {
			return c.Status(400).JSON(fiber.Map{
				"error": "Author and content are required",
			})
		}

		// Save to database
		if err := database.DBConn.Create(&post).Error; err != nil {
			return c.Status(500).JSON(fiber.Map{
				"error": "Failed to create post",
			})
		}

		return c.Status(fiber.StatusCreated).JSON(post)
	})

	// Basic get endpoint to get all posts
	app.Get("/posts", func(c *fiber.Ctx) error {
		var posts []database.Post
		if err := database.DBConn.Find(&posts).Error; err != nil {
			return c.Status(500).JSON(fiber.Map{
				"error": "Failed to retrieve posts",
			})
		}
		return c.JSON(posts)
	})

	// Basic get endpoint to get a single post
	app.Get("/posts/:id", func(c *fiber.Ctx) error {
		id := c.Params("id")

		// Validate UUID format
		if _, err := uuid.Parse(id); err != nil {
			return c.Status(400).JSON(fiber.Map{
				"error": "Invalid UUID format",
			})
		}

		var post database.Post
		if err := database.DBConn.First(&post, "id = ?", id).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				return c.Status(404).JSON(fiber.Map{
					"error": "Post not found",
				})
			}
			return c.Status(500).JSON(fiber.Map{
				"error": "Failed to retrieve post",
			})
		}

		return c.JSON(post)
	})

	app.Listen(os.Getenv("LISTEN"))
}
