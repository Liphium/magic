package main

import (
	"encoding/json"
	"testing"

	"github.com/Liphium/magic/v4"
	"github.com/Liphium/magic/v4/mconfig"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"resty.dev/v3"
)

// Let Magic start the app and needed containers right here so it runs before any tests can run.
func TestMain(m *testing.M) {
	magic.PrepareTesting(m, BuildConfig())
}

func TestApp(t *testing.T) {
	ConnectSurreal()

	t.Run("post is added properly", func(t *testing.T) {
		client := resty.New()
		defer client.Close()

		testPost := Post{
			Author:  "Test",
			Content: "Hello world!",
		}

		res, err := client.R().
			SetBody(testPost).
			Post(GetPath() + "/posts")
		assert.Nil(t, err)
		assert.Equal(t, fiber.StatusCreated, res.StatusCode())

		var created Post
		assert.Nil(t, json.Unmarshal(res.Bytes(), &created))

		testPost.ID = created.ID
		assert.EqualValues(t, testPost, created)

		// You can check if it was actually created straight in the database.
		// In this case it might not be so useful, but when you call complex endpoints, direct access to the database can be really handy to be able to fully test if the endpoint did the correct thing.
		post, err := getPost(created.ID)
		assert.Nil(t, err)
		assert.NotNil(t, post)

		assert.EqualValues(t, testPost, *post)
	})

	t.Run("posts can be retrived", func(t *testing.T) {
		client := resty.New()
		defer client.Close()

		// You can clear databases here, but if you don't rely on an empty database for a test, just not doing it is fine, too.
		assert.Nil(t, magic.GetTestRunner().RunInstruction(mconfig.InstructionClearTables))

		// Yes, you can call scripts in here to make your life a little easier.
		if err := SeedDatabase(); err != nil {
			t.Fatal("Couldn't seed database:", err)
		}

		res, err := client.R().
			Get(GetPath() + "/posts")
		assert.Nil(t, err)
		assert.Equal(t, fiber.StatusOK, res.StatusCode())

		var posts []Post
		assert.Nil(t, json.Unmarshal(res.Bytes(), &posts))
		assert.Equal(t, len(SamplePosts), len(posts))

		// SurrealDB doesn't guarantee any ordering for SELECT * FROM posts, so compare without relying on it.
		byAuthor := map[string]string{}
		for _, post := range posts {
			byAuthor[post.Author] = post.Content
		}
		for _, sample := range SamplePosts {
			assert.Equal(t, sample.Content, byAuthor[sample.Author])
		}
	})

	// You may want to add more tests in a real app...
}
