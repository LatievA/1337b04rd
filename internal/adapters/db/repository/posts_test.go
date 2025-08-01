package repository

import (
	"1337b04rd/internal/domain"
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"testing"
	"time"

	_ "github.com/lib/pq" // PostgreSQL driver
)

var testDB *sql.DB

func TestMain(m *testing.M) {
	// Get test database configuration from environment variables
	dbHost := getEnv("TEST_DB_HOST", "localhost")
	dbPort := getEnv("TEST_DB_PORT", "5432")
	dbUser := getEnv("TEST_DB_USER", "postgres")
	dbPassword := getEnv("TEST_DB_PASSWORD", "postgres")
	dbName := getEnv("TEST_DB_NAME", "testdb")

	// Create test database connection string
	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		dbHost, dbPort, dbUser, dbPassword, dbName)

	// Connect to test database
	var err error
	testDB, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("Failed to connect to test database: %v", err)
	}

	// Verify connection
	if err = testDB.Ping(); err != nil {
		log.Fatalf("Failed to ping test database: %v", err)
	}

	// Run migrations
	if err := runMigrations(testDB); err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}

	// Run tests
	code := m.Run()

	// Cleanup
	if err := cleanupDatabase(testDB); err != nil {
		fmt.Printf("Failed to clean up database: %v\n", err)
	}
	testDB.Close()
	os.Exit(code)
}

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

func runMigrations(db *sql.DB) error {
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS posts (
			id SERIAL PRIMARY KEY,
			session_id INTEGER NOT NULL,
			username TEXT NOT NULL,
			title TEXT NOT NULL,
			content TEXT NOT NULL,
			image_url TEXT,
			created_at TIMESTAMP DEFAULT NOW(),
			archived_at TIMESTAMP NOT NULL,
			is_archived BOOLEAN DEFAULT FALSE
		);
		
		CREATE TABLE IF NOT EXISTS comments (
			id SERIAL PRIMARY KEY,
			user_id INTEGER NOT NULL,
			post_id INTEGER NOT NULL REFERENCES posts(id) ON DELETE CASCADE,
			parent_id INTEGER DEFAULT 0,
			content TEXT NOT NULL,
			created_at TIMESTAMP DEFAULT NOW()
		);
	`)
	return err
}

func cleanupDatabase(db *sql.DB) error {
	_, err := db.Exec(`
		DROP TABLE IF EXISTS comments;
		DROP TABLE IF EXISTS posts;
	`)
	return err
}

func TestPostRepository_Save(t *testing.T) {
	repo := NewPostRepository(testDB)
	ctx := context.Background()

	post := &domain.Post{
		UserID:     1,
		Username:   "testuser",
		Title:      "Test Post",
		Content:    "Test content",
		ImageURL:   "http://example.com/image.jpg",
		ArchivedAt: time.Now().Add(15 * time.Minute),
	}

	id, err := repo.Save(ctx, post)
	if err != nil {
		t.Fatalf("Save failed: %v", err)
	}
	if id <= 0 {
		t.Error("Expected positive ID")
	}
}

func TestPostRepository_FindByID(t *testing.T) {
	repo := NewPostRepository(testDB)
	ctx := context.Background()

	// Create a post first
	post := &domain.Post{
		UserID:     1,
		Username:   "testuser",
		Title:      "Test FindByID",
		Content:    "Content",
		ArchivedAt: time.Now().Add(15 * time.Minute),
	}
	id, _ := repo.Save(ctx, post)

	// Test retrieval
	found, err := repo.FindByID(ctx, id)
	if err != nil {
		t.Fatalf("FindByID failed: %v", err)
	}
	if found.Title != "Test FindByID" {
		t.Errorf("Expected title 'Test FindByID', got '%s'", found.Title)
	}
}

func TestPostRepository_FindAll(t *testing.T) {
	repo := NewPostRepository(testDB)
	ctx := context.Background()

	// Create test posts
	repo.Save(ctx, &domain.Post{
		UserID:     1,
		Username:   "user1",
		Title:      "Active 1",
		Content:    "Content",
		ArchivedAt: time.Now().Add(1 * time.Hour),
	})
	repo.Save(ctx, &domain.Post{
		UserID:     2,
		Username:   "user2",
		Title:      "Archived 1",
		Content:    "Content",
		ArchivedAt: time.Now().Add(-1 * time.Hour),
		Archived:   true,
	})
	repo.Save(ctx, &domain.Post{
		UserID:     3,
		Username:   "user3",
		Title:      "Active 2",
		Content:    "Content",
		ArchivedAt: time.Now().Add(2 * time.Hour),
	})

	// Test active posts
	active, err := repo.FindAll(ctx, false)
	if err != nil {
		t.Fatalf("FindAll(active) failed: %v", err)
	}
	if len(active) != 2 {
		t.Errorf("Expected 2 active posts, got %d", len(active))
	}

	// Test archived posts
	archived, err := repo.FindAll(ctx, true)
	if err != nil {
		t.Fatalf("FindAll(archived) failed: %v", err)
	}
	if len(archived) != 1 {
		t.Errorf("Expected 1 archived post, got %d", len(archived))
	}
}

func TestPostRepository_Update(t *testing.T) {
	repo := NewPostRepository(testDB)
	ctx := context.Background()

	// Create a post
	post := &domain.Post{
		UserID:     1,
		Username:   "user",
		Title:      "Original Title",
		Content:    "Original Content",
		ArchivedAt: time.Now().Add(15 * time.Minute),
	}
	id, _ := repo.Save(ctx, post)

	// Update it
	updatePost := &domain.Post{
		ID:         id,
		Title:      "Updated Title",
		Content:    "Updated Content",
		ArchivedAt: time.Now().Add(30 * time.Minute),
		Archived:   true,
	}
	err := repo.Update(ctx, updatePost)
	if err != nil {
		t.Fatalf("Update failed: %v", err)
	}

	// Verify update
	updated, _ := repo.FindByID(ctx, id)
	if updated.Title != "Updated Title" {
		t.Errorf("Update failed, expected 'Updated Title', got '%s'", updated.Title)
	}
	if !updated.Archived {
		t.Error("Expected post to be archived")
	}
}

func TestPostRepository_ArchiveExpired(t *testing.T) {
	repo := NewPostRepository(testDB)
	ctx := context.Background()

	// Create an expired post
	expiredPost := &domain.Post{
		UserID:     1,
		Username:   "user",
		Title:      "Expired Post",
		Content:    "Content",
		ArchivedAt: time.Now().Add(-1 * time.Hour), // Expired 1 hour ago
	}
	id, _ := repo.Save(ctx, expiredPost)

	// Archive expired posts
	err := repo.ArchiveExpired(ctx)
	if err != nil {
		t.Fatalf("ArchiveExpired failed: %v", err)
	}

	// Verify it was archived
	archivedPost, _ := repo.FindByID(ctx, id)
	if !archivedPost.Archived {
		t.Error("Expired post was not archived")
	}
}

func TestPostRepository_Add15Min(t *testing.T) {
	repo := NewPostRepository(testDB)
	ctx := context.Background()

	// Create a post
	originalTime := time.Now().Add(30 * time.Minute)
	post := &domain.Post{
		UserID:     1,
		Username:   "user",
		Title:      "Test Add15Min",
		Content:    "Content",
		ArchivedAt: originalTime,
	}
	id, _ := repo.Save(ctx, post)

	// Add 15 minutes
	err := repo.Add15Min(ctx, id)
	if err != nil {
		t.Fatalf("Add15Min failed: %v", err)
	}

	// Verify time was extended
	updated, _ := repo.FindByID(ctx, id)
	expected := originalTime.Add(15 * time.Minute)
	if !updated.ArchivedAt.Equal(expected) {
		t.Errorf("Expected archive time %v, got %v", expected, updated.ArchivedAt)
	}
}

func TestPostRepository_NotFoundCases(t *testing.T) {
	repo := NewPostRepository(testDB)
	ctx := context.Background()

	// Test FindByID with non-existent ID
	_, err := repo.FindByID(ctx, 9999)
	if err == nil {
		t.Error("Expected error for non-existent post")
	}

	// Test Update with non-existent ID
	err = repo.Update(ctx, &domain.Post{ID: 9999})
	if err == nil {
		t.Error("Expected error for updating non-existent post")
	}

	// Test Add15Min with non-existent ID
	err = repo.Add15Min(ctx, 9999)
	if err == nil {
		t.Error("Expected error for extending non-existent post")
	}
}
