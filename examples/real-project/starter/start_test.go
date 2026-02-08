package starter_test

import (
	"encoding/json"
	"real-project/database"
	"real-project/starter"
	"testing"

	"github.com/Liphium/magic/v2"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"resty.dev/v3"
)

// Let Magic start the app and needed containers right here so it runs before any tests can run.
func TestMain(m *testing.M) {
	magic.PrepareTesting(m, starter.BuildMagicConfig())
}

func TestApp(t *testing.T) {
	database.Connect()

	t.Run("post is added properly", func(t *testing.T) {
		client := resty.New()
		defer client.Close()

		testPost := database.Post{
			Author:  "Test",
			Content: "Hello world!",
		}

		res, err := client.R().
			SetBody(testPost).
			Post(starter.GetPath() + "/posts")
		assert.Nil(t, err)
		assert.Equal(t, fiber.StatusCreated, res.StatusCode())

		var created database.Post
		assert.Nil(t, json.Unmarshal(res.Bytes(), &created))

		testPost.ID = created.ID
		assert.EqualValues(t, testPost, created)

		// You can check if it was actually created straight in the database.
		// In this case it might not be so useful, but when you call complex endpoints, direct access to the database
		// can be really handy to be able to fully test if the endpoint did the correct thing.
		var post database.Post
		err = database.DBConn.Where("id = ?", created.ID).Take(&post).Error
		assert.Nil(t, err)

		assert.EqualValues(t, testPost, created)
	})

	t.Run("posts can be retrived", func(t *testing.T) {
		client := resty.New()
		defer client.Close()

		// You can clear databases here, but if you don't rely on an empty database for a test, just not doing it is fine, too.
		magic.GetTestRunner().ClearDatabases()

		// Yes, you can call scripts in here to make your life a little easier.
		if err := starter.SeedDatabase(); err != nil {
			t.Fatal("Couldn't seed database:", err)
		}

		res, err := client.R().
			Get(starter.GetPath() + "/posts")
		assert.Nil(t, err)
		assert.Equal(t, fiber.StatusOK, res.StatusCode())

		var posts []database.Post
		assert.Nil(t, json.Unmarshal(res.Bytes(), &posts))
		assert.Equal(t, len(starter.SamplePosts), len(posts))

		for i, post := range starter.SamplePosts {
			assert.Equal(t, post.Author, posts[i].Author)
			assert.Equal(t, post.Content, posts[i].Content)
		}
	})

	// You may want to add more tests in a real app...
}
