package main

import (
	"fmt"
	"os"

	"github.com/Liphium/magic/v4"
	"github.com/gofiber/fiber/v2"
)

func main() {
	magic.Start(BuildConfig())
}

func Start() {
	// Connect to the database
	ConnectSurreal()

	// Create the actual web app
	app := fiber.New()

	// This message is just here for explanation.
	fmt.Println()
	fmt.Println("Welcome, wizard! That's how easy it is to get Magic up and running. Well, I hope it actually all went well for you...")
	fmt.Println("Anyway, now that we're up and running, you can open another terminal and run scripts using: 'go run . -r <script>' or list of all of them using 'go run . --scripts'!")
	fmt.Println("Thanks for using Magic!")

	// Basic insertion endpoint to create a new post
	app.Post("/posts", func(c *fiber.Ctx) error {
		var post Post

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
		created, err := createPost(post)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{
				"error": "Failed to create post",
			})
		}

		return c.Status(fiber.StatusCreated).JSON(created)
	})

	// Basic get endpoint to get all posts
	app.Get("/posts", func(c *fiber.Ctx) error {
		posts, err := getPosts()
		if err != nil {
			return c.Status(500).JSON(fiber.Map{
				"error": "Failed to retrieve posts",
			})
		}
		return c.JSON(posts)
	})

	// Basic get endpoint to get a single post
	app.Get("/posts/:id", func(c *fiber.Ctx) error {
		post, err := getPost(c.Params("id"))
		if err != nil {
			if err == ErrInvalidPostID {
				return c.Status(400).JSON(fiber.Map{
					"error": "Invalid post ID format",
				})
			}
			return c.Status(500).JSON(fiber.Map{
				"error": "Failed to retrieve post",
			})
		}

		if post == nil {
			return c.Status(404).JSON(fiber.Map{
				"error": "Post not found",
			})
		}

		return c.JSON(post)
	})

	// Add a startup hook to notify Magic of the app start
	app.Hooks().OnListen(func(listenData fiber.ListenData) error {
		if fiber.IsChild() {
			return nil
		}

		// Tell Magic the app has started: Makes sure tests start to run after this.
		magic.AppStarted()
		return nil
	})

	app.Listen(os.Getenv("LISTEN"))
}
