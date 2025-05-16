package web

import (
	"encoding/json"
	"net/http"
	"test-news/internal/database/models"
)

func PostsPageHandler(w http.ResponseWriter, r *http.Request) {
	PostsPage().Render(r.Context(), w)
}

func PostsListHandler(w http.ResponseWriter, r *http.Request) {
	// Fetch posts from the API endpoint
	resp, err := http.Get("http://localhost:8080/api/posts")
	if err != nil {
		http.Error(w, "Failed to fetch posts", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	var posts []models.Post
	if err := json.NewDecoder(resp.Body).Decode(&posts); err != nil {
		http.Error(w, "Failed to parse posts data", http.StatusInternalServerError)
		return
	}

	// Render the posts list component
	PostsList(posts).Render(r.Context(), w)
}

func PostDetailHandler(w http.ResponseWriter, r *http.Request) {
	// Get the post ID from the URL
	id := r.URL.Path[len("/api/posts/"):]

	// Fetch the post from the API endpoint
	resp, err := http.Get("http://localhost:8080/api/posts/" + id)
	if err != nil {
		http.Error(w, "Failed to fetch post", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	var post models.Post
	if err := json.NewDecoder(resp.Body).Decode(&post); err != nil {
		http.Error(w, "Failed to parse post data", http.StatusInternalServerError)
		return
	}

	PostDetail(post).Render(r.Context(), w)
}
