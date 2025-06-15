package domain

type PostRepository interface {
	Create(post *Post) error
	GetByID(id int) (*Post, error)
	GetAll() ([]*Post, error)
	Update(post *Post) error
	Delete(id int) error
}

type UserRepository interface {
	Create(user *User) error
	GetByID(id int) (*User, error)
}

type CommentRepository interface {
	Create(comment *Comment) error
	GetByID(id int) (*Comment, error)
}