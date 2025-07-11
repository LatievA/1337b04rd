package domain

import "time"

type Post struct {
	ID        int
	UserID    int
	Title     string
	Content   string
	ImageURL  *string
	Archived  bool
	Comments  []Comment
	CreatedAt time.Time
	UpdatedAt time.Time
}

type Comment struct {
	ID        int
	UserID    int
	ParentID  int
	Content   string
	CreatedAt time.Time
}

type User struct {
	ID        int
	Name      string
	AvatarURL *string
}
