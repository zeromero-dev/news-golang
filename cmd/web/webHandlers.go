package web

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
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
	DeletePage("", "").Render(r.Context(), w)
}

func DeleteConfirmHandler(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		DeletePage("", "Failed to parse form: "+err.Error()).Render(r.Context(), w)
		return
	}

	postId := r.FormValue("postId")
	if postId == "" {
		DeletePage("", "Post ID is required").Render(r.Context(), w)
		return
	}

	// Check if the post exists
	resp, err := http.Get("http://localhost:8080/api/posts/" + postId)
	if err != nil {
		DeletePage("", "Failed to check post: "+err.Error()).Render(r.Context(), w)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		DeletePage("", "Post not found. Please check the ID and try again.").Render(r.Context(), w)
		return
	}

	// Render the confirmation page
	DeleteConfirmPage(postId).Render(r.Context(), w)
}

func DeleteExecuteHandler(w http.ResponseWriter, r *http.Request) {
	// Get the post ID from the URL
	// The URL format is /web/delete/execute/:id
	pathParts := strings.Split(r.URL.Path, "/")
	postId := pathParts[len(pathParts)-1]

	// Create a new request to delete the post
	req, err := http.NewRequest("DELETE", "http://localhost:8080/api/posts/"+postId, nil)
	if err != nil {
		DeletePage("", "Failed to create request: "+err.Error()).Render(r.Context(), w)
		return
	}

	// Send the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		DeletePage("", "Failed to delete post: "+err.Error()).Render(r.Context(), w)
		return
	}
	defer resp.Body.Close()

	// Check if the delete was successful
	if resp.StatusCode != http.StatusOK {
		DeletePage("", "Failed to delete post. Please try again.").Render(r.Context(), w)
		return
	}

	// Render the delete page with a success message
	DeletePage("Post deleted successfully!", "").Render(r.Context(), w)
}
func UpdatePageHandler(w http.ResponseWriter, r *http.Request) {
	UpdatePage().Render(r.Context(), w)
}
