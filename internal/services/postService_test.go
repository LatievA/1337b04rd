package services

import (
	"1337b04rd/internal/domain"
	"context"
	"errors"
	"sync"
	"testing"
	"time"
)

type mockPostRepository struct {
	posts       map[int]*domain.Post
	lastID      int
	mu          sync.Mutex // For thread safety
	saveErr     error
	findByIDErr error
	findAllErr  error
	updateErr   error
	archiveErr  error
	add15MinErr error
}

func newMockPostRepo() *mockPostRepository {
	return &mockPostRepository{
		posts:  make(map[int]*domain.Post),
		lastID: 0,
	}
}

func (m *mockPostRepository) Save(ctx context.Context, post *domain.Post) (int, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.saveErr != nil {
		return 0, m.saveErr
	}

	m.lastID++
	post.ID = m.lastID
	// Set default timestamps if not set
	if post.CreatedAt.IsZero() {
		post.CreatedAt = time.Now()
	}
	if post.ArchivedAt.IsZero() {
		post.ArchivedAt = time.Now().Add(15 * time.Minute)
	}
	m.posts[post.ID] = post
	return post.ID, nil
}

func (m *mockPostRepository) FindByID(ctx context.Context, id int) (*domain.Post, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.findByIDErr != nil {
		return nil, m.findByIDErr
	}

	post, exists := m.posts[id]
	if !exists {
		return nil, errors.New("post not found")
	}
	return post, nil
}

func (m *mockPostRepository) FindAll(ctx context.Context, archived bool) ([]*domain.Post, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.findAllErr != nil {
		return nil, m.findAllErr
	}

	var result []*domain.Post
	for _, post := range m.posts {
		if post.Archived == archived {
			result = append(result, post)
		}
	}
	return result, nil
}

func (m *mockPostRepository) Update(ctx context.Context, post *domain.Post) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.updateErr != nil {
		return m.updateErr
	}

	if _, exists := m.posts[post.ID]; !exists {
		return errors.New("post not found")
	}

	m.posts[post.ID] = post
	return nil
}

func (m *mockPostRepository) ArchiveExpired(ctx context.Context) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.archiveErr != nil {
		return m.archiveErr
	}

	now := time.Now()
	for _, post := range m.posts {
		if !post.Archived && now.After(post.ArchivedAt) {
			post.Archived = true
		}
	}
	return nil
}

func (m *mockPostRepository) Add15Min(ctx context.Context, postID int) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	post, exists := m.posts[postID]
	if !exists {
		return errors.New("post not found")
	}

	// Add 15 minutes to the ArchivedAt time
	post.ArchivedAt = post.ArchivedAt.Add(15 * time.Minute)
	return nil
}

type mockCommentRepository struct {
	comments        map[int]*domain.Comment
	postComments    map[int][]*domain.Comment
	lastID          int
	saveErr         error
	findByIDErr     error
	findByPostIDErr error
}

type mockUserRepository struct {
	users          map[int]*domain.User
	usersByToken   map[string]*domain.User
	lastID         int
	mu             sync.Mutex
	findByIDErr    error
	findByTokenErr error
	saveErr        error
	updateNameErr  error
	getUserIDErr   error
}

func newMockUserRepo() *mockUserRepository {
	return &mockUserRepository{
		users:        make(map[int]*domain.User),
		usersByToken: make(map[string]*domain.User),
		lastID:       0,
	}
}

func (m *mockUserRepository) FindByID(ctx context.Context, userID int) (*domain.User, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.findByIDErr != nil {
		return nil, m.findByIDErr
	}

	user, exists := m.users[userID]
	if !exists {
		return nil, errors.New("user not found")
	}
	return user, nil
}

func (m *mockUserRepository) FindBySessionToken(ctx context.Context, sessionToken string) (*domain.User, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.findByTokenErr != nil {
		return nil, m.findByTokenErr
	}

	user, exists := m.usersByToken[sessionToken]
	if !exists {
		return nil, errors.New("user not found")
	}
	return user, nil
}

func (m *mockUserRepository) Save(ctx context.Context, user *domain.User) (int, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.saveErr != nil {
		return 0, m.saveErr
	}

	m.lastID++
	user.ID = m.lastID
	m.users[user.ID] = user
	if user.SessionToken != "" {
		m.usersByToken[user.SessionToken] = user
	}
	return user.ID, nil
}

func (m *mockUserRepository) UpdateName(ctx context.Context, userID int, newName string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.updateNameErr != nil {
		return m.updateNameErr
	}

	user, exists := m.users[userID]
	if !exists {
		return errors.New("user not found")
	}

	user.Name = newName
	return nil
}

func (m *mockUserRepository) GetUserIDBySessionToken(ctx context.Context, sessionToken string) (int, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.getUserIDErr != nil {
		return 0, m.getUserIDErr
	}

	user, exists := m.usersByToken[sessionToken]
	if !exists {
		return 0, errors.New("user not found")
	}
	return user.ID, nil
}

func (m *mockCommentRepository) Save(ctx context.Context, comment *domain.Comment) (int, error) {
	if m.saveErr != nil {
		return 0, m.saveErr
	}
	m.lastID++
	comment.ID = m.lastID
	m.comments[comment.ID] = comment
	if m.postComments[comment.PostID] == nil {
		m.postComments[comment.PostID] = []*domain.Comment{}
	}
	m.postComments[comment.PostID] = append(m.postComments[comment.PostID], comment)
	return comment.ID, nil
}

func (m *mockCommentRepository) FindByID(ctx context.Context, id int) (*domain.Comment, error) {
	if m.findByIDErr != nil {
		return nil, m.findByIDErr
	}
	comment, exists := m.comments[id]
	if !exists {
		return nil, errors.New("not found")
	}
	return comment, nil
}

func (m *mockCommentRepository) FindByPostID(ctx context.Context, postID int) ([]*domain.Comment, error) {
	if m.findByPostIDErr != nil {
		return nil, m.findByPostIDErr
	}
	return m.postComments[postID], nil
}

func newMockCommentRepo() *mockCommentRepository {
	return &mockCommentRepository{
		comments:     make(map[int]*domain.Comment),
		postComments: make(map[int][]*domain.Comment),
		lastID:       0,
	}
}

func TestPostService_CreatePost(t *testing.T) {
	tests := []struct {
		name        string
		userID      int
		username    string
		title       string
		content     string
		imageURL    string
		saveErr     error
		expectedErr bool
	}{
		{
			name:        "successful creation",
			userID:      1,
			username:    "testuser",
			title:       "Test Post",
			content:     "Test content",
			imageURL:    "http://example.com/image.jpg",
			saveErr:     nil,
			expectedErr: false,
		},
		{
			name:        "repository error",
			userID:      1,
			username:    "testuser",
			title:       "Test Post",
			content:     "Test content",
			imageURL:    "http://example.com/image.jpg",
			saveErr:     errors.New("database error"),
			expectedErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := newMockPostRepo()
			repo.saveErr = tt.saveErr

			service := NewPostService(repo, &mockCommentRepository{}, &mockUserRepository{})
			post, err := service.CreatePost(context.Background(), tt.userID, tt.username, tt.title, tt.content, tt.imageURL)

			if tt.expectedErr {
				if err == nil {
					t.Error("expected error, got nil")
				}
				return
			}

			if err != nil {
				t.Errorf("unexpected error: %v", err)
				return
			}

			if post.ID == 0 {
				t.Error("expected post ID to be set")
			}
			if post.UserID != tt.userID {
				t.Errorf("expected UserID %d, got %d", tt.userID, post.UserID)
			}
			if post.Username != tt.username {
				t.Errorf("expected Username %s, got %s", tt.username, post.Username)
			}
			if post.Title != tt.title {
				t.Errorf("expected Title %s, got %s", tt.title, post.Title)
			}
			if post.Content != tt.content {
				t.Errorf("expected Content %s, got %s", tt.content, post.Content)
			}
			if post.ImageURL != tt.imageURL {
				t.Errorf("expected ImageURL %s, got %s", tt.imageURL, post.ImageURL)
			}
			if post.CreatedAt.IsZero() {
				t.Error("expected CreatedAt to be set")
			}
			if post.ArchivedAt.IsZero() {
				t.Error("expected ArchivedAt to be set")
			}
		})
	}
}

func TestPostService_GetPostByID(t *testing.T) {
	tests := []struct {
		name        string
		postID      int
		prepopulate bool
		findByIDErr error
		expectedErr bool
	}{
		{
			name:        "existing post",
			postID:      1,
			prepopulate: true,
			findByIDErr: nil,
			expectedErr: false,
		},
		{
			name:        "non-existent post",
			postID:      2,
			prepopulate: false,
			findByIDErr: nil,
			expectedErr: true,
		},
		{
			name:        "repository error",
			postID:      1,
			prepopulate: true,
			findByIDErr: errors.New("database error"),
			expectedErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := newMockPostRepo()
			repo.findByIDErr = tt.findByIDErr

			if tt.prepopulate {
				repo.Save(context.Background(), &domain.Post{
					UserID:     1,
					Username:   "testuser",
					Title:      "Test Post",
					Content:    "Test content",
					ImageURL:   "http://example.com/image.jpg",
					CreatedAt:  time.Now(),
					ArchivedAt: time.Now().Add(15 * time.Minute),
				})
			}

			service := NewPostService(repo, &mockCommentRepository{}, &mockUserRepository{})
			post, err := service.GetPostByID(context.Background(), tt.postID)

			if tt.expectedErr {
				if err == nil {
					t.Error("expected error, got nil")
				}
				return
			}

			if err != nil {
				t.Errorf("unexpected error: %v", err)
				return
			}

			if post.ID != tt.postID {
				t.Errorf("expected post ID %d, got %d", tt.postID, post.ID)
			}
		})
	}
}

func TestPostService_ListPosts(t *testing.T) {
	tests := []struct {
		name        string
		archived    bool
		prepopulate []*domain.Post
		findAllErr  error
		expectedLen int
		expectedErr bool
	}{
		{
			name:     "list active posts",
			archived: false,
			prepopulate: []*domain.Post{
				{ID: 1, Archived: false},
				{ID: 2, Archived: true},
				{ID: 3, Archived: false},
			},
			findAllErr:  nil,
			expectedLen: 2,
			expectedErr: false,
		},
		{
			name:     "list archived posts",
			archived: true,
			prepopulate: []*domain.Post{
				{ID: 1, Archived: false},
				{ID: 2, Archived: true},
				{ID: 3, Archived: false},
			},
			findAllErr:  nil,
			expectedLen: 1,
			expectedErr: false,
		},
		{
			name:     "repository error",
			archived: false,
			prepopulate: []*domain.Post{
				{ID: 1, Archived: false},
			},
			findAllErr:  errors.New("database error"),
			expectedLen: 0,
			expectedErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := newMockPostRepo()
			repo.findAllErr = tt.findAllErr

			for _, post := range tt.prepopulate {
				repo.Save(context.Background(), post)
			}

			service := NewPostService(repo, &mockCommentRepository{}, &mockUserRepository{})
			posts, err := service.ListPosts(context.Background(), tt.archived)

			if tt.expectedErr {
				if err == nil {
					t.Error("expected error, got nil")
				}
				return
			}

			if err != nil {
				t.Errorf("unexpected error: %v", err)
				return
			}

			if len(posts) != tt.expectedLen {
				t.Errorf("expected %d posts, got %d", tt.expectedLen, len(posts))
			}

			for _, post := range posts {
				if post.Archived != tt.archived {
					t.Errorf("expected archived=%v, got %v", tt.archived, post.Archived)
				}
			}
		})
	}
}

func TestPostService_ArchiveOldPosts(t *testing.T) {
	tests := []struct {
		name        string
		prepopulate []*domain.Post
		archiveErr  error
		expectedErr bool
	}{
		{
			name: "archive expired posts",
			prepopulate: []*domain.Post{
				{ID: 1, Archived: false, ArchivedAt: time.Now().Add(-1 * time.Hour)},
				{ID: 2, Archived: false, ArchivedAt: time.Now().Add(1 * time.Hour)},
			},
			archiveErr:  nil,
			expectedErr: false,
		},
		{
			name: "repository error",
			prepopulate: []*domain.Post{
				{ID: 1, Archived: false, ArchivedAt: time.Now().Add(-1 * time.Hour)},
			},
			archiveErr:  errors.New("database error"),
			expectedErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := newMockPostRepo()
			repo.archiveErr = tt.archiveErr

			for _, post := range tt.prepopulate {
				repo.Save(context.Background(), post)
			}

			service := NewPostService(repo, &mockCommentRepository{}, &mockUserRepository{})
			err := service.ArchiveOldPosts(context.Background())

			if tt.expectedErr {
				if err == nil {
					t.Error("expected error, got nil")
				}
				return
			}

			if err != nil {
				t.Errorf("unexpected error: %v", err)
				return
			}

			for _, post := range tt.prepopulate {
				storedPost, _ := repo.FindByID(context.Background(), post.ID)
				if time.Now().After(post.ArchivedAt) && !storedPost.Archived {
					t.Errorf("expected post %d to be archived", post.ID)
				}
			}
		})
	}
}

func TestPostService_AddTimeToPostLifetime(t *testing.T) {
	tests := []struct {
		name        string
		postID      int
		prepopulate bool
		add15MinErr error
		expectedErr bool
	}{
		{
			name:        "successful time extension",
			postID:      1,
			prepopulate: true,
			expectedErr: false,
		},
		{
			name:        "post not found",
			postID:      2,
			prepopulate: false,
			expectedErr: true,
		},
		{
			name:        "repository error",
			postID:      1,
			prepopulate: true,
			add15MinErr: errors.New("db error"),
			expectedErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := newMockPostRepo()
			repo.add15MinErr = tt.add15MinErr

			var originalTime time.Time
			if tt.prepopulate {
				post := &domain.Post{
					ID:         tt.postID,
					ArchivedAt: time.Now().Add(10 * time.Minute),
				}
				repo.Save(context.Background(), post)
				originalTime = post.ArchivedAt
			}

			service := NewPostService(repo, &mockCommentRepository{}, &mockUserRepository{})
			err := service.AddTimeToPostLifetime(context.Background(), tt.postID)

			if tt.expectedErr {
				if err == nil {
					t.Error("expected error but got nil")
				}
				return
			}

			if err != nil {
				t.Errorf("unexpected error: %v", err)
				return
			}

			updatedPost, _ := repo.FindByID(context.Background(), tt.postID)
			if !updatedPost.ArchivedAt.Equal(originalTime.Add(15 * time.Minute)) {
				t.Errorf("expected archive time to be increased by 15 minutes")
			}
		})
	}
}

func TestCommentService_AddComment(t *testing.T) {
	tests := []struct {
		name        string
		userID      int
		postID      int
		parentID    int
		content     string
		postExists  bool
		saveErr     error
		expectedErr bool
	}{
		{
			name:        "successful comment",
			userID:      1,
			postID:      1,
			content:     "Test comment",
			postExists:  true,
			expectedErr: false,
		},
		{
			name:        "post not found",
			postID:      2,
			postExists:  false,
			expectedErr: true,
		},
		{
			name:        "repository error",
			postID:      1,
			postExists:  true,
			saveErr:     errors.New("save error"),
			expectedErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			postRepo := newMockPostRepo()
			commentRepo := newMockCommentRepo()
			commentRepo.saveErr = tt.saveErr

			if tt.postExists {
				postRepo.Save(context.Background(), &domain.Post{ID: tt.postID})
			}

			service := NewCommentService(commentRepo, postRepo)
			comment, err := service.AddComment(context.Background(), tt.userID, tt.postID, tt.parentID, tt.content)

			if tt.expectedErr {
				if err == nil {
					t.Error("expected error, got nil")
				}
				return
			}

			if err != nil {
				t.Errorf("unexpected error: %v", err)
				return
			}

			if comment.ID == 0 {
				t.Error("expected comment ID to be set")
			}
			if comment.PostID != tt.postID {
				t.Errorf("expected PostID %d, got %d", tt.postID, comment.PostID)
			}
			if comment.Content != tt.content {
				t.Errorf("expected content %s, got %s", tt.content, comment.Content)
			}
		})
	}
}

func TestCommentService_GetCommentsByPostID(t *testing.T) {
	tests := []struct {
		name            string
		postID          int
		prepopulate     bool
		findByPostIDErr error
		expectedLen     int
		expectedErr     bool
	}{
		{
			name:        "comments found",
			postID:      1,
			prepopulate: true,
			expectedLen: 2,
		},
		{
			name:        "no comments",
			postID:      2,
			prepopulate: false,
			expectedLen: 0,
		},
		{
			name:            "repository error",
			postID:          1,
			prepopulate:     true,
			findByPostIDErr: errors.New("find error"),
			expectedErr:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			commentRepo := newMockCommentRepo()
			commentRepo.findByPostIDErr = tt.findByPostIDErr

			if tt.prepopulate {
				commentRepo.Save(context.Background(), &domain.Comment{PostID: tt.postID})
				commentRepo.Save(context.Background(), &domain.Comment{PostID: tt.postID})
			}

			service := NewCommentService(commentRepo, newMockPostRepo())
			comments, err := service.GetCommentsByPostID(context.Background(), tt.postID)

			if tt.expectedErr {
				if err == nil {
					t.Error("expected error, got nil")
				}
				return
			}

			if err != nil {
				t.Errorf("unexpected error: %v", err)
				return
			}

			if len(comments) != tt.expectedLen {
				t.Errorf("expected %d comments, got %d", tt.expectedLen, len(comments))
			}
		})
	}
}

func TestPostRepository_ConcurrentAccess(t *testing.T) {
	repo := newMockPostRepo()
	post := &domain.Post{Title: "Test"}
	id, _ := repo.Save(context.Background(), post)

	var wg sync.WaitGroup
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			_, _ = repo.FindByID(context.Background(), id)
			_ = repo.Update(context.Background(), post)
		}()
	}
	wg.Wait()

	// Verify no data corruption occurred
	updatedPost, err := repo.FindByID(context.Background(), id)
	if err != nil {
		t.Fatalf("failed to find post: %v", err)
	}
	if updatedPost.Title != "Test" {
		t.Errorf("expected title 'Test', got '%s'", updatedPost.Title)
	}
}
