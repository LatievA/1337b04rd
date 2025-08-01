package repository

import (
	"1337b04rd/internal/domain"
	"context"
	"database/sql"
	"log"
	"os"
	"testing"
	"time"

	_ "github.com/lib/pq"
)

// need to run before repository testing:

// docker run -d --name test-postgres \
//   -e POSTGRES_USER=postgres \
//   -e POSTGRES_PASSWORD=postgres \
//   -e POSTGRES_DB=testdb \
//   -p 5432:5432 \
//   postgres:15


var testDB *sql.DB

func TestMain(m *testing.M) {
	connStr := "host=localhost port=5432 user=postgres password=postgres dbname=testdb sslmode=disable"
	var err error
	testDB, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}

	if err = testDB.Ping(); err != nil {
		log.Fatal(err)
	}

	setupTestDatabase(testDB)
	code := m.Run()
	cleanupTestDatabase(testDB)
	testDB.Close()
	os.Exit(code)
}

func setupTestDatabase(db *sql.DB) {
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS user_sessions (
			id SERIAL PRIMARY KEY,
			session_token TEXT UNIQUE NOT NULL,
			name TEXT NOT NULL, 
			avatar_url TEXT NOT NULL DEFAULT 'default_avatar.png',
			created_at TIMESTAMP DEFAULT NOW(),
			expires_at TIMESTAMP DEFAULT NOW() + INTERVAL '1 week'
		);
		
		CREATE TABLE IF NOT EXISTS posts (
			id SERIAL PRIMARY KEY,
			session_id INTEGER REFERENCES user_sessions(id) ON DELETE CASCADE,
			username TEXT,
			title TEXT NOT NULL,
			content TEXT NOT NULL,
			image_url TEXT DEFAULT '',
			created_at TIMESTAMP DEFAULT NOW(),
			archived_at TIMESTAMP DEFAULT NOW() + INTERVAL '15 minutes',
			is_archived BOOLEAN DEFAULT FALSE
		);
		
		CREATE TABLE IF NOT EXISTS comments (
			id SERIAL PRIMARY KEY,
			session_id INTEGER REFERENCES user_sessions(id) ON DELETE CASCADE,
			post_id INTEGER REFERENCES posts(id) ON DELETE CASCADE,
			parent_comment_id INTEGER DEFAULT 0,
			content TEXT NOT NULL,
			created_at TIMESTAMP DEFAULT NOW()
		);
	`)
	if err != nil {
		log.Fatal(err)
	}

	_, err = db.Exec(`
		TRUNCATE 
			comments, 
			posts, 
			user_sessions 
		RESTART IDENTITY CASCADE
	`)
	if err != nil {
		log.Fatal(err)
	}
}

func cleanupTestDatabase(db *sql.DB) {
	_, err := db.Exec(`
		DROP TABLE IF EXISTS 
			comments, 
			posts, 
			user_sessions
	`)
	if err != nil {
		log.Fatal(err)
	}
}

func createTestUser(t *testing.T, db *sql.DB, tokenSuffix string) int {
	t.Helper()
	var userID int
	err := db.QueryRow(`
		INSERT INTO user_sessions (session_token, name)
		VALUES ($1, $2)
		RETURNING id
	`, 
		"test_token_"+tokenSuffix, 
		"Test User",
	).Scan(&userID)
	if err != nil {
		t.Fatalf("Failed to create test user: %v", err)
	}
	return userID
}

func createTestPost(t *testing.T, db *sql.DB, userID int) int {
	t.Helper()
	var postID int
	err := db.QueryRow(`
		INSERT INTO posts (session_id, username, title, content)
		VALUES ($1, $2, $3, $4)
		RETURNING id
	`, 
		userID,
		"testuser",
		"Test Post", 
		"Test Content",
	).Scan(&postID)
	if err != nil {
		t.Fatalf("Failed to create test post: %v", err)
	}
	return postID
}

func createTestComment(t *testing.T, db *sql.DB, userID, postID int) int {
	t.Helper()
	var commentID int
	err := db.QueryRow(`
		INSERT INTO comments (session_id, post_id, content, parent_comment_id)
		VALUES ($1, $2, $3, $4)
		RETURNING id
	`, 
		userID,
		postID,
		"Test Comment",
		0, // Explicitly set parent_comment_id to 0
	).Scan(&commentID)
	if err != nil {
		t.Fatalf("Failed to create test comment: %v", err)
	}
	return commentID
}

func TestPostRepository_Save(t *testing.T) {
	repo := NewPostRepository(testDB)
	userID := createTestUser(t, testDB, "save")

	post := &domain.Post{
		UserID:   userID,
		Username: "testuser",
		Title:    "Test Post",
		Content:  "This is a test post",
		ImageURL: "test.jpg",
	}

	id, err := repo.Save(context.Background(), post)
	if err != nil {
		t.Fatalf("Save failed: %v", err)
	}

	if id <= 0 {
		t.Errorf("Expected positive post ID, got %d", id)
	}
}

func TestPostRepository_FindByID(t *testing.T) {
	repo := NewPostRepository(testDB)
	userID := createTestUser(t, testDB, "findbyid")
	postID := createTestPost(t, testDB, userID)
	createTestComment(t, testDB, userID, postID)

	post, err := repo.FindByID(context.Background(), postID)
	if err != nil {
		t.Fatalf("FindByID failed: %v", err)
	}

	if post.ID != postID {
		t.Errorf("Expected post ID %d, got %d", postID, post.ID)
	}
}

func TestPostRepository_FindAll(t *testing.T) {
	repo := NewPostRepository(testDB)
	userID := createTestUser(t, testDB, "findall")
	postID := createTestPost(t, testDB, userID)
	createTestComment(t, testDB, userID, postID)

	// Test active posts
	posts, err := repo.FindAll(context.Background(), false)
	if err != nil {
		t.Fatalf("FindAll failed: %v", err)
	}

	if len(posts) == 0 {
		t.Error("Expected at least 1 post, got 0")
	}

	// Test archived posts
	_, err = testDB.Exec("UPDATE posts SET is_archived = true WHERE id = $1", postID)
	if err != nil {
		t.Fatalf("Failed to archive post: %v", err)
	}

	archivedPosts, err := repo.FindAll(context.Background(), true)
	if err != nil {
		t.Fatalf("FindAll archived failed: %v", err)
	}

	if len(archivedPosts) == 0 {
		t.Error("Expected at least 1 archived post, got 0")
	}
}

func TestPostRepository_Update(t *testing.T) {
	repo := NewPostRepository(testDB)
	userID := createTestUser(t, testDB, "update")
	postID := createTestPost(t, testDB, userID)

	updatedPost := &domain.Post{
		ID:         postID,
		Title:      "Updated Title",
		Content:    "Updated Content",
		ImageURL:   "updated.jpg",
		Archived:   true,
		ArchivedAt: time.Now(),
	}

	err := repo.Update(context.Background(), updatedPost)
	if err != nil {
		t.Fatalf("Update failed: %v", err)
	}

	// Verify the update
	var title string
	err = testDB.QueryRow(
		"SELECT title FROM posts WHERE id = $1",
		postID,
	).Scan(&title)
	if err != nil {
		t.Fatalf("Failed to verify updated post: %v", err)
	}

	if title != updatedPost.Title {
		t.Errorf("Expected title %s, got %s", updatedPost.Title, title)
	}
}

func TestPostRepository_ArchiveExpired(t *testing.T) {
	repo := NewPostRepository(testDB)
	userID := createTestUser(t, testDB, "archive")

	// Create an expired post
	_, err := testDB.Exec(`
		INSERT INTO posts(session_id, username, title, content, archived_at, is_archived)
		VALUES ($1, $2, $3, $4, NOW() - INTERVAL '1 day', false)
	`, userID, "testuser", "Expired Post", "Content")
	if err != nil {
		t.Fatalf("Failed to create test post: %v", err)
	}

	err = repo.ArchiveExpired(context.Background())
	if err != nil {
		t.Fatalf("ArchiveExpired failed: %v", err)
	}

	// Verify the post was archived
	var isArchived bool
	err = testDB.QueryRow(
		"SELECT is_archived FROM posts WHERE is_archived = true",
	).Scan(&isArchived)
	if err != nil {
		t.Fatalf("Failed to verify archived post: %v", err)
	}
}

func TestPostRepository_Add15Min(t *testing.T) {
	repo := NewPostRepository(testDB)
	userID := createTestUser(t, testDB, "add15min")
	postID := createTestPost(t, testDB, userID)

	// Get original archived_at time
	var originalTime time.Time
	err := testDB.QueryRow(
		"SELECT archived_at FROM posts WHERE id = $1", 
		postID,
	).Scan(&originalTime)
	if err != nil {
		t.Fatalf("Failed to get original time: %v", err)
	}

	err = repo.Add15Min(context.Background(), postID)
	if err != nil {
		t.Fatalf("Add15Min failed: %v", err)
	}

	// Verify the time was increased
	var newTime time.Time
	err = testDB.QueryRow(
		"SELECT archived_at FROM posts WHERE id = $1",
		postID,
	).Scan(&newTime)
	if err != nil {
		t.Fatalf("Failed to get updated time: %v", err)
	}

	if !newTime.After(originalTime) {
		t.Errorf("Expected time to be increased, original: %v, new: %v", originalTime, newTime)
	}
}