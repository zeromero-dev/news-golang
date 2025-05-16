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
	r.GET("/", s.HelloWorldHandler)

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

	r.POST("/hello", func(c *gin.Context) {
		web.HelloWebHandler(c.Writer, c.Request)
	})

	// Web routes that return HTML
	// Serve the HTML page for posts
	r.GET("/web/posts", func(c *gin.Context) {
		web.PostsPageHandler(c.Writer, c.Request)
	})

	r.GET("/api/posts/list", func(c *gin.Context) {
		web.PostsListHandler(c.Writer, c.Request)
	})

	r.GET("/api/posts/detail/:id", func(c *gin.Context) {
		id := c.Param("id")
		c.Request.URL.Path = "/api/posts/" + id
		web.PostDetailHandler(c.Writer, c.Request)
	})

	r.GET("/web/upload", func(c *gin.Context) {
		web.UploadPageHandler(c.Writer, c.Request)
	})

	r.GET("/web/delete", func(c *gin.Context) {
		web.DeletePageHandler(c.Writer, c.Request)
	})

	return r
}
