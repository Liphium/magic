package starter_test

import (
	"real-project/database"
	"real-project/starter"
	"real-project/util"
	"reflect"
	"testing"

	"github.com/Liphium/magic"
	"github.com/Liphium/magic/mrunner"
)

func TestApp(t *testing.T) {

	// By calling this function, Magic will start up your entire app in the background. You can then call its endpoints from in here
	// and test against the output. Additionally, you may access any state your app may store in maps or anything since you're in the same
	// process. Services or something could also be tested.
	magic.TestRunner(t, starter.BuildMagicConfig(), "app", func(t *testing.T, r *mrunner.Runner) {
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

		t.Run("invalid post is rejected", func(t *testing.T) {
			created, err := util.Post[database.Post](starter.GetPath()+"/posts", database.Post{}, util.Headers{})
			if err == nil {
				t.Fatalf("Invalid post could be created: created=%+v", created)
			}
		})

		// You can clear databases here, but if you don't rely on an empty database in tests below, it's fine.
		r.ClearDatabases()

		// Yes, you can call scripts in here to make your life a little easier.
		if err := starter.SeedDatabase(r, ""); err != nil {
			t.Fatal("Couldn't seed database:", err)
		}

		t.Run("posts can be retrived", func(t *testing.T) {

		})

		// You may want to add more tests in a real app...
	})
}
