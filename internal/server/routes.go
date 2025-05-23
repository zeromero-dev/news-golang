package server

import (
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	"io/fs"
	"test-news/cmd/web"

	"github.com/a-h/templ"
)

func (s *Server) RegisterRoutes() http.Handler {
	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173"}, // Add your frontend URL
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"},
		AllowHeaders:     []string{"Accept", "Authorization", "Content-Type"},
		AllowCredentials: true, // Enable cookies/auth
	}))

	// this is not a concern of duplication
	// API routes that return JSON

	r.GET("/health", s.healthHandler)

	r.GET("/api/posts", s.GetPostsHandler)
	r.POST("/api/posts", s.CreatePostHandler)
	r.GET("/api/posts/:id", s.GetPostHandler)
	r.PUT("/api/posts/:id", s.UpdatePostHandler)
	r.DELETE("/api/posts/:id", s.DeletePostHandler)

	staticFiles, _ := fs.Sub(web.Files, "assets")
	r.StaticFS("/assets", http.FS(staticFiles))

	r.GET("/web", func(c *gin.Context) {
		templ.Handler(web.HelloForm()).ServeHTTP(c.Writer, c.Request)
	})

	// Web routes that return HTML
	// Serve the HTML page for posts
	r.GET("/web/posts", func(c *gin.Context) {
		web.PostsPageHandler(c.Writer, c.Request)
	})

	r.GET("/api/posts/list", func(c *gin.Context) {
		web.PostsListHandler(c.Writer, c.Request)
	})

	r.GET("/web/upload", func(c *gin.Context) {
		web.UploadPageHandler(c.Writer, c.Request)
	})

	r.POST("/web/upload/submit", func(c *gin.Context) {
		web.UploadSubmitHandler(c.Writer, c.Request)
	})

	r.GET("/web/update", func(c *gin.Context) {
		web.UpdatePageHandler(c.Writer, c.Request)
	})

	//Delete page routes + handlers
	r.GET("/web/delete", func(c *gin.Context) {
		web.DeletePageHandler(c.Writer, c.Request)
	})

	r.POST("/web/delete/confirm", func(c *gin.Context) {
		web.DeleteConfirmHandler(c.Writer, c.Request)
	})

	r.POST("/web/delete/execute/:id", func(c *gin.Context) {
		web.DeleteExecuteHandler(c.Writer, c.Request)
	})

	r.GET("/web/posts/:id", func(c *gin.Context) {
		web.PostDetailPageHandler(c.Writer, c.Request)
	})
	return r
}
