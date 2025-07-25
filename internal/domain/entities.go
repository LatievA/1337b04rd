package domain

import "time"

type Post struct {
	ID         int
	Username     string
	Title      string
	Content    string
	ImageURL   *string
	Comments   []*Comment
	CreatedAt  time.Time
	ArchivedAt time.Time
	Archived   bool
}

type Comment struct {
	ID        int
	UserID    int
	PostID    int
	ParentID  *int
	Content   string
	CreatedAt time.Time
}

type User struct {
	ID           int
	SessionToken string
	Name         string
	AvatarURL    *string
	ExpiresAt    time.Time
}
