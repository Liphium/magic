package starter_test

import (
	"real-project/database"
	"real-project/starter"
	"real-project/util"
	"reflect"
	"slices"
	"testing"

	"github.com/Liphium/magic"
	"github.com/google/uuid"
)

// Let Magic start the app and needed containers right here so it runs before any tests can run.
func TestMain(m *testing.M) {
	magic.PrepareTesting(m, "main", starter.BuildMagicConfig())
}

func TestApp(t *testing.T) {
	database.Connect()

	t.Run("post is added properly", func(t *testing.T) {
		testPost := database.Post{
			Author:  "Test",
			Content: "Hello world!",
		}

		created, err := util.Post[database.Post](starter.GetPath()+"/posts", testPost, util.Headers{})
		if err != nil {
			t.Fatal("Failed to create post:", err)
		}

		testPost.ID = created.ID
		if !reflect.DeepEqual(testPost, created) {
			t.Fatalf("Post is not equal: expected=%+v got=%+v", testPost, created)
		}

		// You can check if it was actually created straight in the database.
		// In this case it might not be so useful, but when you call complex endpoints, direct access to the database
		// can be really handy to be able to fully test if the endpoint did the correct thing.
		var post database.Post
		if err := database.DBConn.Where("id = ?", created.ID).Take(&post).Error; err != nil {
			t.Fatal("Couldn't retrive created post:", err)
		}

		if !reflect.DeepEqual(testPost, created) {
			t.Fatalf("Post is not equal to one in database: expected=%+v got=%+v", testPost, post)
		}

	})

	t.Run("posts can be retrived", func(t *testing.T) {

		// You can clear databases here, but if you don't rely on an empty database for a test, just not doing it is fine, too.
		magic.GetTestRunner().ClearDatabases()

		// Yes, you can call scripts in here to make your life a little easier.
		if err := starter.SeedDatabase(magic.GetTestRunner(), ""); err != nil {
			t.Fatal("Couldn't seed database:", err)
		}

		posts, err := util.Get[[]database.Post](starter.GetPath()+"/posts", util.Headers{})
		if err != nil {
			t.Fatalf("Couldn't get posts from the backend: %v", err)
		}

		if !slices.EqualFunc(starter.SamplePosts, posts, func(e1, e2 database.Post) bool {
			e1.ID = uuid.Nil
			e2.ID = uuid.Nil
			return reflect.DeepEqual(e1, e2)
		}) {
			t.Fatalf("Gotten posts don't match sample posts: expected=%v got=%v", starter.SamplePosts, posts)
		}
	})

	// You may want to add more tests in a real app...
}
