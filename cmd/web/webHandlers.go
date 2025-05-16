package web

import (
	"bytes"
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

func PostDetailPageHandler(w http.ResponseWriter, r *http.Request) {
	// Extract the post ID from the URL
	// The URL format is /web/posts/:id
	pathParts := strings.Split(r.URL.Path, "/")
	id := pathParts[len(pathParts)-1]

	// Fetch the post from the API endpoint
	resp, err := http.Get("http://localhost:8080/api/posts/" + id)
	if err != nil {
		http.Error(w, "Failed to fetch post: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		http.Error(w, "Post not found", http.StatusNotFound)
		return
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		http.Error(w, "Failed to read response body: "+err.Error(), http.StatusInternalServerError)
		return
	}

	var post models.Post
	if err := json.Unmarshal(body, &post); err != nil {
		http.Error(w, "Failed to parse post data: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Render the post detail page
	PostDetailPage(post).Render(r.Context(), w)
}
func UploadPageHandler(w http.ResponseWriter, r *http.Request) {
	UploadPage("", "").Render(r.Context(), w)
}

func UploadSubmitHandler(w http.ResponseWriter, r *http.Request) {
	// Parse the form data
	if err := r.ParseForm(); err != nil {
		UploadPage("", "Failed to parse form data: "+err.Error()).Render(r.Context(), w)
		return
	}

	// Get the form values
	title := r.FormValue("title")
	author := r.FormValue("author")
	content := r.FormValue("content")

	// Validate the form values
	if title == "" || author == "" || content == "" {
		UploadPage("", "All fields are required").Render(r.Context(), w)
		return
	}

	// Create the payload
	payload := map[string]string{
		"title":   title,
		"author":  author,
		"content": content,
	}

	// Convert the payload to JSON
	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		UploadPage("", "Failed to create JSON payload: "+err.Error()).Render(r.Context(), w)
		return
	}

	// Create a new request to create the post
	req, err := http.NewRequest("POST", "http://localhost:8080/api/posts", bytes.NewBuffer(jsonPayload))
	if err != nil {
		UploadPage("", "Failed to create request: "+err.Error()).Render(r.Context(), w)
		return
	}
	req.Header.Set("Content-Type", "application/json")

	// Send the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		UploadPage("", "Failed to create post: "+err.Error()).Render(r.Context(), w)
		return
	}
	defer resp.Body.Close()

	// Check if the creation was successful
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		// Read the error response
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			UploadPage("", "Failed to read error response: "+err.Error()).Render(r.Context(), w)
			return
		}

		UploadPage("", "Failed to create post: "+string(body)).Render(r.Context(), w)
		return
	}

	// Render the upload page with a success message
	UploadPage("Post created successfully!", "").Render(r.Context(), w)
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
