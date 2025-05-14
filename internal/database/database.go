package database

import (
	"context"
	"fmt"
	"log"
	"os"
	"test-news/internal/database/models"
	"time"

	_ "github.com/joho/godotenv/autoload"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Service interface {
	Health() map[string]string
	GetPosts() ([]*models.Post, error)
	CreatePost(post *models.Post) error
	GetPost(id string) (*models.Post, error)
	UpdatePost(id string, post *models.Post) error
	DeletePost(id string) error
}

type service struct {
	db *mongo.Client
}

var (
	host = os.Getenv("BLUEPRINT_DB_HOST")
	port = os.Getenv("BLUEPRINT_DB_PORT")
	//database = os.Getenv("BLUEPRINT_DB_DATABASE")
)

func New() Service {
	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(fmt.Sprintf("mongodb://%s:%s", host, port)))

	if err != nil {
		log.Fatal(err)

	}
	return &service{
		db: client,
	}
}

func (s *service) Health() map[string]string {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	err := s.db.Ping(ctx, nil)
	if err != nil {
		log.Fatalf("db down: %v", err)
	}

	return map[string]string{
		"message": "It's healthy",
	}
}

func (s *service) GetPosts() ([]*models.Post, error) {
	collection := s.db.Database("test-news").Collection("posts")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	cursor, err := collection.Find(ctx, map[string]interface{}{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var posts []*models.Post
	for cursor.Next(ctx) {
		var post models.Post
		if err := cursor.Decode(&post); err != nil {
			return nil, err
		}
		posts = append(posts, &post)
	}
	if err := cursor.Err(); err != nil {
		return nil, err
	}
	return posts, nil
}

func (s *service) CreatePost(post *models.Post) error {
	collection := s.db.Database("test-news").Collection("posts")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	post.CreatedAt = time.Now()
	post.UpdatedAt = time.Now()

	_, err := collection.InsertOne(ctx, post)
	if err != nil {
		return err
	}
	return nil
}

func (s *service) GetPost(id string) (*models.Post, error) {
	collection := s.db.Database("test-news").Collection("posts")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var post models.Post
	err := collection.FindOne(ctx, map[string]interface{}{"_id": id}).Decode(&post)
	if err != nil {
		return nil, err
	}
	return &post, nil
}

func (s *service) UpdatePost(id string, post *models.Post) error {
	collection := s.db.Database("test-news").Collection("posts")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	post.UpdatedAt = time.Now()

	_, err := collection.UpdateOne(ctx, map[string]interface{}{"_id": id}, map[string]interface{}{"$set": post})
	if err != nil {
		return err
	}
	return nil
}

func (s *service) DeletePost(id string) error {
	return fmt.Errorf("not implemented")
}
