package web

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"test-news/internal/database/models"
)

type PostsResponse struct {
	Count int           `json:"count"`
	Data  []models.Post `json:"data"`
}

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

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		http.Error(w, "Failed to read response body", http.StatusInternalServerError)
		return
	}

	var postsResp PostsResponse
	if err := json.Unmarshal(body, &postsResp); err != nil {
		http.Error(w, "Failed to parse posts data", http.StatusInternalServerError)
		return
	}

	// Render the posts list component
	PostsList(postsResp.Data).Render(r.Context(), w)
}

func PostDetailHandler(w http.ResponseWriter, r *http.Request) {
	// Extract the post ID from the URL
	id := r.URL.Path[len("/api/posts/detail/"):]

	// Fetch the post from the API endpoint
	resp, err := http.Get("http://localhost:8080/api/posts/" + id)
	if err != nil {
		http.Error(w, "Failed to fetch post", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		http.Error(w, "Failed to read response body", http.StatusInternalServerError)
		return
	}

	var post models.Post
	if err := json.Unmarshal(body, &post); err != nil {
		http.Error(w, "Failed to parse post data", http.StatusInternalServerError)
		return
	}

	// Render the post detail component
	PostDetail(post).Render(r.Context(), w)
}

func UploadPageHandler(w http.ResponseWriter, r *http.Request) {
	UploadPage().Render(r.Context(), w)
}
func DeletePageHandler(w http.ResponseWriter, r *http.Request) {
	DeletePage().Render(r.Context(), w)
}