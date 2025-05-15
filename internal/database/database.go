package database

import (
	"context"
	"fmt"
	"log"
	"os"
	"test-news/internal/database/models"
	"time"

	"github.com/joho/godotenv"
	_ "github.com/joho/godotenv/autoload"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
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
	host     = os.Getenv("BLUEPRINT_DB_HOST")
	port     = os.Getenv("BLUEPRINT_DB_PORT")
	username = os.Getenv("BLUEPRINT_DB_USERNAME")
	password = os.Getenv("BLUEPRINT_DB_ROOT_PASSWORD")
	database = "news_feed" // Default database name, can be overridden with env var
)

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found, using default or system environment variables")
	}
	// Check if database name is provided via environment variable
	if dbName := os.Getenv("BLUEPRINT_DB_DATABASE"); dbName != "" {
		database = dbName
	}
}

func New() Service {
	clientOpts := options.Client().
		ApplyURI(fmt.Sprintf("mongodb://%s:%s", host, port))

	// Done for tests, casuse the test container doesn't require authentication
	// and crashes if we try to do this any other way
	if username != "" && password != "" {
		cred := options.Credential{
			Username: username,
			Password: password,
		}
		clientOpts.SetAuth(cred)
	}

	client, err := mongo.Connect(context.Background(), clientOpts)
	if err != nil {
		log.Fatal(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}

	log.Println("Successfully connected to MongoDB")

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

func (s *service) getCollection() *mongo.Collection {
	return s.db.Database(database).Collection("posts")
}

func (s *service) GetPosts() ([]*models.Post, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	collection := s.getCollection()

	//time descending (newest first)
	findOptions := options.Find()
	findOptions.SetSort(bson.D{{Key: "created_at", Value: -1}})

	cursor, err := collection.Find(ctx, bson.M{}, findOptions)
	if err != nil {
		return nil, fmt.Errorf("error fetching posts: %w", err)
	}
	defer cursor.Close(ctx)

	var posts []*models.Post
	if err = cursor.All(ctx, &posts); err != nil {
		return nil, fmt.Errorf("error decoding posts: %w", err)
	}

	return posts, nil
}

func (s *service) CreatePost(post *models.Post) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	collection := s.getCollection()

	now := time.Now()
	post.CreatedAt = now
	post.UpdatedAt = now

	post.ID = primitive.NewObjectID()

	_, err := collection.InsertOne(ctx, post)
	if err != nil {
		return fmt.Errorf("error creating post: %w", err)
	}

	return nil
}

func (s *service) GetPost(id string) (*models.Post, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	collection := s.getCollection()

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, fmt.Errorf("invalid ID format: %w", err)
	}

	var post models.Post
	err = collection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&post)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("post not found")
		}
		return nil, fmt.Errorf("error fetching post: %w", err)
	}

	return &post, nil
}

func (s *service) UpdatePost(id string, post *models.Post) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	collection := s.getCollection()

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return fmt.Errorf("invalid ID format: %w", err)
	}

	// Update the modification time
	post.UpdatedAt = time.Now()

	// don't change the original ID
	post.ID = objectID

	update := bson.M{
		"$set": post,
	}

	_, err = collection.UpdateOne(
		ctx,
		bson.M{"_id": objectID},
		update,
	)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			return fmt.Errorf("post not found")
		}
		return fmt.Errorf("error updating post: %w", err)
	}

	return nil
}

func (s *service) DeletePost(id string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	collection := s.getCollection()

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return fmt.Errorf("invalid ID format: %w", err)
	}

	result, err := collection.DeleteOne(ctx, bson.M{"_id": objectID})
	if err != nil {
		return fmt.Errorf("error deleting post: %w", err)
	}

	if result.DeletedCount == 0 {
		return fmt.Errorf("post not found")
	}

	return nil
}
