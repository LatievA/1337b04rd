package repository

import (
	"1337b04rd/internal/domain"
	"context"
	"database/sql"
	"testing"
	"time"

	_ "github.com/lib/pq" // PostgreSQL driver
)

func runUserMigrations(db *sql.DB) error {
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS user_sessions (
			id SERIAL PRIMARY KEY,
			session_token TEXT NOT NULL UNIQUE,
			name TEXT NOT NULL,
			avatar_url TEXT NOT NULL,
			expires_at TIMESTAMP NOT NULL
		);
	`)
	return err
}

func cleanupUserDatabase(db *sql.DB) error {
	_, err := db.Exec(`DROP TABLE IF EXISTS user_sessions;`)
	return err
}

func TestUserRepository_Save(t *testing.T) {
	repo := NewUserRepository(testDB)
	ctx := context.Background()

	user := &domain.User{
		SessionToken: "token123",
		Name:         "Test User",
		AvatarURL:    "http://example.com/avatar.jpg",
		ExpiresAt:    time.Now().Add(24 * time.Hour),
	}

	id, err := repo.Save(ctx, user)
	if err != nil {
		t.Fatalf("Save failed: %v", err)
	}
	if id <= 0 {
		t.Error("Expected positive ID")
	}
}

func TestUserRepository_FindByID(t *testing.T) {
	repo := NewUserRepository(testDB)
	ctx := context.Background()

	// Create a user first
	user := &domain.User{
		SessionToken: "token_findbyid",
		Name:         "FindByID Test",
		AvatarURL:    "http://example.com/avatar.jpg",
		ExpiresAt:    time.Now().Add(24 * time.Hour),
	}
	id, _ := repo.Save(ctx, user)

	// Test retrieval
	found, err := repo.FindByID(ctx, id)
	if err != nil {
		t.Fatalf("FindByID failed: %v", err)
	}
	if found.Name != "FindByID Test" {
		t.Errorf("Expected name 'FindByID Test', got '%s'", found.Name)
	}
}

func TestUserRepository_FindBySessionToken(t *testing.T) {
	repo := NewUserRepository(testDB)
	ctx := context.Background()

	// Create a user
	token := "unique_session_token_123"
	user := &domain.User{
		SessionToken: token,
		Name:         "Session Token Test",
		AvatarURL:    "http://example.com/avatar.jpg",
		ExpiresAt:    time.Now().Add(24 * time.Hour),
	}
	repo.Save(ctx, user)

	// Test retrieval
	found, err := repo.FindBySessionToken(ctx, token)
	if err != nil {
		t.Fatalf("FindBySessionToken failed: %v", err)
	}
	if found.SessionToken != token {
		t.Errorf("Expected session token '%s', got '%s'", token, found.SessionToken)
	}
}

func TestUserRepository_GetUserIDBySessionToken(t *testing.T) {
	repo := NewUserRepository(testDB)
	ctx := context.Background()

	// Create a user
	token := "token_for_id_lookup"
	user := &domain.User{
		SessionToken: token,
		Name:         "ID Lookup Test",
		AvatarURL:    "http://example.com/avatar.jpg",
		ExpiresAt:    time.Now().Add(24 * time.Hour),
	}
	expectedID, _ := repo.Save(ctx, user)

	// Test retrieval
	userID, err := repo.GetUserIDBySessionToken(ctx, token)
	if err != nil {
		t.Fatalf("GetUserIDBySessionToken failed: %v", err)
	}
	if userID != expectedID {
		t.Errorf("Expected user ID %d, got %d", expectedID, userID)
	}
}

func TestUserRepository_UpdateName(t *testing.T) {
	repo := NewUserRepository(testDB)
	ctx := context.Background()

	// Create a user
	user := &domain.User{
		SessionToken: "update_token",
		Name:         "Original Name",
		AvatarURL:    "http://example.com/avatar.jpg",
		ExpiresAt:    time.Now().Add(24 * time.Hour),
	}
	id, _ := repo.Save(ctx, user)

	// Update name
	newName := "Updated Name"
	err := repo.UpdateName(ctx, id, newName)
	if err != nil {
		t.Fatalf("UpdateName failed: %v", err)
	}

	// Verify update
	updated, _ := repo.FindByID(ctx, id)
	if updated.Name != newName {
		t.Errorf("Expected name '%s', got '%s'", newName, updated.Name)
	}
}

func TestUserRepository_NotFoundCases(t *testing.T) {
	repo := NewUserRepository(testDB)
	ctx := context.Background()

	// Test FindByID with non-existent ID
	_, err := repo.FindByID(ctx, 9999)
	if err == nil {
		t.Error("Expected error for non-existent user ID")
	}

	// Test FindBySessionToken with non-existent token
	_, err = repo.FindBySessionToken(ctx, "non_existent_token")
	if err == nil {
		t.Error("Expected error for non-existent session token")
	}

	// Test GetUserIDBySessionToken with non-existent token
	_, err = repo.GetUserIDBySessionToken(ctx, "non_existent_token")
	if err == nil {
		t.Error("Expected error for non-existent session token")
	}

	// Test UpdateName with non-existent ID
	err = repo.UpdateName(ctx, 9999, "New Name")
	if err == nil {
		t.Error("Expected error for non-existent user ID")
	}
}

func TestUserRepository_UniqueSessionToken(t *testing.T) {
	repo := NewUserRepository(testDB)
	ctx := context.Background()

	// First user with unique token
	token := "unique_token_123"
	user1 := &domain.User{
		SessionToken: token,
		Name:         "User 1",
		AvatarURL:    "http://example.com/avatar1.jpg",
		ExpiresAt:    time.Now().Add(24 * time.Hour),
	}
	_, err := repo.Save(ctx, user1)
	if err != nil {
		t.Fatalf("First save failed: %v", err)
	}

	// Second user with same token (should fail)
	user2 := &domain.User{
		SessionToken: token,
		Name:         "User 2",
		AvatarURL:    "http://example.com/avatar2.jpg",
		ExpiresAt:    time.Now().Add(24 * time.Hour),
	}
	_, err = repo.Save(ctx, user2)
	if err == nil {
		t.Error("Expected error when saving duplicate session token")
	}
}
