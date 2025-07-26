package handlers

type CreatePostRequest struct {
	UserSession string `json:"session_token"`
	Title string `json:"title"`
	Content string `json:"content"`
	ImageURL *string `json:"image_url"`
}