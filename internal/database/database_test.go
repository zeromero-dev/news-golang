package database

import (
	"context"
	"log"
	"test-news/internal/database/models"
	"testing"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/mongodb"
)

func mustStartMongoContainer() (func(context.Context, ...testcontainers.TerminateOption) error, error) {
	dbContainer, err := mongodb.Run(context.Background(), "mongo:latest")
	if err != nil {
		return nil, err
	}

	dbHost, err := dbContainer.Host(context.Background())
	if err != nil {
		return dbContainer.Terminate, err
	}

	dbPort, err := dbContainer.MappedPort(context.Background(), "27017/tcp")
	if err != nil {
		return dbContainer.Terminate, err
	}

	host = dbHost
	port = dbPort.Port()

	return dbContainer.Terminate, err
}

func TestMain(m *testing.M) {
	teardown, err := mustStartMongoContainer()
	if err != nil {
		log.Fatalf("could not start mongodb container: %v", err)
	}

	// Test container doesn't require authentication
	username = ""
	password = ""
	database = "testdb"

	m.Run()

	if teardown != nil && teardown(context.Background()) != nil {
		log.Fatalf("could not teardown mongodb container: %v", err)
	}
}

func TestNew(t *testing.T) {
	srv := New()
	if srv == nil {
		t.Fatal("New() returned nil")
	}
}

func TestHealth(t *testing.T) {
	srv := New()

	stats := srv.Health()

	if stats["message"] != "It's healthy" {
		t.Fatalf("expected message to be 'It's healthy', got %s", stats["message"])
	}
}

func TestGetPosts(t *testing.T) {
	srv := New()

	posts, err := srv.GetPosts()
	if err != nil {
		t.Fatalf("GetPosts() returned an error: %v", err)
	}

	if len(posts) != 0 {
		t.Fatalf("expected 0 posts, got %d", len(posts))
	}
}

func TestCreatePost(t *testing.T) {
	srv := New()

	post := &models.Post{
		Title:   "Test Post",
		Content: "This is a test post",
	}

	err := srv.CreatePost(post)
	if err != nil {
		t.Fatalf("CreatePost() returned an error: %v", err)
	}

	if post.ID.IsZero() {
		t.Fatal("expected post ID to be set, got zero value")
	}
}

func TestGetPostByID(t *testing.T) {
	srv := New()

	post := &models.Post{
		Title:   "Test Post",
		Content: "This is a test post",
	}

	err := srv.CreatePost(post)
	if err != nil {
		t.Fatalf("CreatePost() returned an error: %v", err)
	}

	fetchedPost, err := srv.GetPost(post.ID.Hex())

	if err != nil {
		t.Fatalf("GetPostByID() returned an error: %v", err)
	}

	if fetchedPost == nil {
		t.Fatal("expected fetched post to be non-nil")
	}

	if fetchedPost.ID != post.ID {
		t.Fatalf("expected fetched post ID to be %s, got %s", post.ID.Hex(), fetchedPost.ID.Hex())
	}
}

func TestGetPostByInvalidID(t *testing.T) {
	srv := New()

	_, err := srv.GetPost("invalid_id")

	if err == nil {
		t.Fatal("expected error for invalid ID, got nil")
	}
}
func TestUpdatePost(t *testing.T) {
	srv := New()

	post := &models.Post{
		Title:   "Test Post",
		Content: "This is a test post",
	}

	err := srv.CreatePost(post)
	if err != nil {
		t.Fatalf("CreatePost() returned an error: %v", err)
	}

	post.Title = "Updated Title"
	err = srv.UpdatePost(post.ID.Hex(), post)
	if err != nil {
		t.Fatalf("UpdatePost() returned an error: %v", err)
	}

	fetchedPost, err := srv.GetPost(post.ID.Hex())
	if err != nil {
		t.Fatalf("GetPostByID() returned an error: %v", err)
	}

	if fetchedPost.Title != "Updated Title" {
		t.Fatalf("expected fetched post title to be 'Updated Title', got %s", fetchedPost.Title)
	}
}

func TestUpdatePostInvalidID(t *testing.T) {
	srv := New()

	post := &models.Post{
		Title:   "Test Post",
		Content: "This is a test post",
	}

	err := srv.CreatePost(post)
	if err != nil {
		t.Fatalf("CreatePost() returned an error: %v", err)
	}

	post.Title = "Updated Title"
	err = srv.UpdatePost("invalid_id", post)
	if err == nil {
		t.Fatal("expected error for invalid ID, got nil")
	}
}
func TestDeletePost(t *testing.T) {
	srv := New()

	post := &models.Post{
		Title:   "Test Post",
		Content: "This is a test post",
	}

	err := srv.CreatePost(post)
	if err != nil {
		t.Fatalf("CreatePost() returned an error: %v", err)
	}

	err = srv.DeletePost(post.ID.Hex())
	if err != nil {
		t.Fatalf("DeletePost() returned an error: %v", err)
	}

	_, err = srv.GetPost(post.ID.Hex())
	if err == nil {
		t.Fatal("expected error for deleted post, got nil")
	}
}

func TestDeletePostInvalidID(t *testing.T) {
	srv := New()

	post := &models.Post{
		Title:   "Test Post",
		Content: "This is a test post",
	}

	err := srv.CreatePost(post)
	if err != nil {
		t.Fatalf("CreatePost() returned an error: %v", err)
	}

	err = srv.DeletePost("invalid_id")
	if err == nil {
		t.Fatal("expected error for invalid ID, got nil")
	}
}
