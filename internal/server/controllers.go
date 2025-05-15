package server

import (
	"net/http"
	"test-news/internal/database/models"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (s *Server) HelloWorldHandler(c *gin.Context) {
	resp := make(map[string]string)
	resp["message"] = "Hello World"

	c.JSON(http.StatusOK, resp)
}

func (s *Server) healthHandler(c *gin.Context) {
	c.JSON(http.StatusOK, s.db.Health())
}

// GetPostsHandler returns all posts
func (s *Server) GetPostsHandler(c *gin.Context) {
	posts, err := s.db.GetPosts()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to fetch posts",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":  posts,
		"count": len(posts),
	})
}

// CreatePostHandler creates a new post
func (s *Server) CreatePostHandler(c *gin.Context) {
	var post models.Post
	if err := c.ShouldBindJSON(&post); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid input",
			"details": err.Error(),
		})
		return
	}

	// Set default values
	now := time.Now()
	post.CreatedAt = now
	post.UpdatedAt = now

	err := s.db.CreatePost(&post)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to create post",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, post)
}

// GetPostHandler returns a single post by ID
func (s *Server) GetPostHandler(c *gin.Context) {
	id := c.Param("id")

	// Validate ObjectID format
	if !isValidObjectID(id) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid post ID format"})
		return
	}

	post, err := s.db.GetPost(id)
	if err != nil {
		status := http.StatusInternalServerError
		if err.Error() == "post not found" {
			status = http.StatusNotFound
		}

		c.JSON(status, gin.H{
			"error":   "Failed to retrieve post",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, post)
}

// UpdatePostHandler updates an existing post
func (s *Server) UpdatePostHandler(c *gin.Context) {
	id := c.Param("id")

	// Validate ObjectID format
	if !isValidObjectID(id) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid post ID format"})
		return
	}

	var post models.Post
	if err := c.ShouldBindJSON(&post); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid input",
			"details": err.Error(),
		})
		return
	}

	// Update modification time
	post.UpdatedAt = time.Now()

	err := s.db.UpdatePost(id, &post)
	if err != nil {
		status := http.StatusInternalServerError
		if err.Error() == "post not found" {
			status = http.StatusNotFound
		}

		c.JSON(status, gin.H{
			"error":   "Failed to update post",
			"details": err.Error(),
		})
		return
	}

	// Fetch the updated post to return the complete object
	updatedPost, err := s.db.GetPost(id)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"message": "Post updated successfully but couldn't retrieve updated data"})
		return
	}

	c.JSON(http.StatusOK, updatedPost)
}

// DeletePostHandler removes a post
func (s *Server) DeletePostHandler(c *gin.Context) {
	id := c.Param("id")

	// Validate ObjectID format
	if !isValidObjectID(id) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid post ID format"})
		return
	}

	err := s.db.DeletePost(id)
	if err != nil {
		status := http.StatusInternalServerError
		if err.Error() == "post not found" {
			status = http.StatusNotFound
		}

		c.JSON(status, gin.H{
			"error":   "Failed to delete post",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Post deleted successfully"})
}

// Helper function to validate MongoDB ObjectID
func isValidObjectID(id string) bool {
	_, err := primitive.ObjectIDFromHex(id)
	return err == nil
}
