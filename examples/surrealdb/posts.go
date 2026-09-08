package main

import (
	"context"
	"fmt"

	"github.com/surrealdb/surrealdb.go"
	"github.com/surrealdb/surrealdb.go/pkg/models"
)

// Post is the model for a single post on the platform.
type Post struct {
	// The ID of a post. The magic: ignore tag here makes sure this isn't asked for when we use the Post in scripts.
	ID string `json:"id" magic:"ignore"`

	// With the prompt: "" struct tag, we can tell Magic what it should ask for (try running go run . -r create-post c:).
	//
	// Magic additionally supports validation using the validator package. Refer to https://github.com/go-playground/validator for everything possible.
	Author  string `json:"author" prompt:"Author name" validate:"required"`
	Content string `json:"content" prompt:"Content of the post" validate:"required,max=256"`
}

// ErrInvalidPostID is returned when a post ID isn't a valid SurrealDB record ID.
var ErrInvalidPostID = fmt.Errorf("invalid post id")

// dbPost is how a post looks when it comes back from the database, where the ID is a record ID (like posts:7fd0d8...) instead of a plain string.
type dbPost struct {
	ID      *models.RecordID `json:"id"`
	Author  string           `json:"author"`
	Content string           `json:"content"`
}

func (p dbPost) toPost() (Post, error) {
	if p.ID == nil {
		return Post{}, fmt.Errorf("post in database doesn't have an id")
	}
	return Post{
		// Format the record ID as a plain string (posts:<id>) instead of using RecordID.String(),
		// which wraps ID parts containing special characters (like the dashes in UUIDs) in ⟨⟩.
		ID:      fmt.Sprintf("%s:%v", p.ID.Table, p.ID.ID),
		Author:  p.Author,
		Content: p.Content,
	}, nil
}

// createPost saves a new post in the database and returns it with its generated ID.
func createPost(post Post) (Post, error) {
	// SurrealDB returns an array for CREATE on a table (since the ID is generated), so decode into []dbPost.
	created, err := surrealdb.Create[[]dbPost](context.Background(), Surreal, models.Table("posts"), map[string]any{
		"author":  post.Author,
		"content": post.Content,
	})
	if err != nil {
		return Post{}, fmt.Errorf("couldn't create post: %w", err)
	}
	if created == nil || len(*created) != 1 {
		return Post{}, fmt.Errorf("unexpected result creating post: %d records", len(*created))
	}

	return (*created)[0].toPost()
}

// getPosts retrieves all posts from the database.
func getPosts() ([]Post, error) {
	rows, err := surrealdb.Select[[]dbPost](context.Background(), Surreal, models.Table("posts"))
	if err != nil {
		return nil, fmt.Errorf("couldn't retrieve posts: %w", err)
	}

	posts := []Post{}
	for _, row := range *rows {
		post, err := row.toPost()
		if err != nil {
			return nil, err
		}
		posts = append(posts, post)
	}
	return posts, nil
}

// getPost retrieves a single post using its ID (posts:some-id). Returns nil in case the post doesn't exist.
func getPost(id string) (*Post, error) {
	recordID, err := models.ParseRecordID(id)
	if err != nil || recordID.Table != "posts" {
		return nil, ErrInvalidPostID
	}

	row, err := surrealdb.Select[dbPost](context.Background(), Surreal, *recordID)
	if err != nil {
		if surrealdb.IsNotFound(err) {
			return nil, nil
		}
		return nil, fmt.Errorf("couldn't retrieve post: %w", err)
	}
	if row == nil || row.ID == nil {
		return nil, nil
	}

	post, err := row.toPost()
	if err != nil {
		return nil, err
	}
	return &post, nil
}
