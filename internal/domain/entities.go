package domain

import "time"

type Post struct {
	ID        int
	Title     string
	Content   string
	ImageURL  *string
	UserID    int
	New bool
	CreatedAt time.Time
	UpdatedAt time.Time
}

type Comment struct {
	ID int
	UserID int
	ParentID *int
} 

type User struct {
	ID int
	Name string
	AvatarURL *string
}
