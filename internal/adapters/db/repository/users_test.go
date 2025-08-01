package repository

import (
	"1337b04rd/internal/domain"
	"context"
	"strings"
	"testing"
	"time"

	_ "github.com/lib/pq"
)

// docker run -d --name test-postgres \
//   -e POSTGRES_USER=postgres \
//   -e POSTGRES_PASSWORD=postgres \
//   -e POSTGRES_DB=testdb \
//   -p 5432:5432 \
//   postgres:15


func TestUserRepository_Save(t *testing.T) {
	repo := NewUserRepository(testDB)

	t.Run("success", func(t *testing.T) {
		testUser := &domain.User{
			SessionToken: "unique_token_" + time.Now().Format("20060102150405"),
			Name:         "New User",
			AvatarURL:    "new_avatar.png",
		}

		id, err := repo.Save(context.Background(), testUser)
		if err != nil {
			t.Fatalf("Save failed: %v", err)
		}

		if id <= 0 {
			t.Errorf("Expected positive user ID, got %d", id)
		}

		// Verify the user was actually saved
		var dbUser domain.User
		err = testDB.QueryRow(
			"SELECT id, session_token, name, avatar_url FROM user_sessions WHERE id = $1",
			id,
		).Scan(&dbUser.ID, &dbUser.SessionToken, &dbUser.Name, &dbUser.AvatarURL)
		if err != nil {
			t.Fatalf("Failed to verify saved user: %v", err)
		}

		if dbUser.Name != testUser.Name {
			t.Errorf("Expected saved user name %s, got %s", testUser.Name, dbUser.Name)
		}
	})

	t.Run("duplicate token", func(t *testing.T) {
		// Create a unique token for this test
		testToken := "duplicate_test_token_" + time.Now().Format("20060102150405")
		testUser := &domain.User{
			SessionToken: testToken,
			Name:         "Duplicate Test User",
			AvatarURL:    "duplicate_avatar.png",
		}

		// First save should succeed
		_, err := repo.Save(context.Background(), testUser)
		if err != nil {
			t.Fatalf("First save failed: %v", err)
		}

		// Second save with same token should fail
		_, err = repo.Save(context.Background(), testUser)
		if err == nil {
			t.Error("Expected error for duplicate session token, got nil")
		} else {
			// Verify it's the expected constraint violation
			if !isDuplicateKeyError(err) {
				t.Errorf("Expected duplicate key error, got: %v", err)
			}
		}
	})
}

// Helper function to check for duplicate key errors
func isDuplicateKeyError(err error) bool {
	if err == nil {
		return false
	}
	// This is a simple check - you might need to adjust based on your database driver
	return strings.Contains(err.Error(), "duplicate key") || 
	       strings.Contains(err.Error(), "violates unique constraint")
}


func TestUserRepository_FindByID(t *testing.T) {
	repo := NewUserRepository(testDB)

	// Setup test data
	testUser := domain.User{
		SessionToken: "find_by_id_token",
		Name:         "Test User",
		AvatarURL:    "test_avatar.png",
		ExpiresAt:    time.Now().Add(24 * time.Hour),
	}

	// Insert test user
	var id int
	err := testDB.QueryRow(
		"INSERT INTO user_sessions(session_token, name, avatar_url, expires_at) VALUES ($1, $2, $3, $4) RETURNING id",
		testUser.SessionToken, testUser.Name, testUser.AvatarURL, testUser.ExpiresAt,
	).Scan(&id)
	if err != nil {
		t.Fatalf("Failed to insert test user: %v", err)
	}

	t.Run("success", func(t *testing.T) {
		user, err := repo.FindByID(context.Background(), id)
		if err != nil {
			t.Fatalf("FindByID failed: %v", err)
		}

		if user.ID != id || user.Name != testUser.Name {
			t.Errorf("Expected user ID %d with name %s, got ID %d with name %s",
				id, testUser.Name, user.ID, user.Name)
		}
	})

	t.Run("not found", func(t *testing.T) {
		_, err := repo.FindByID(context.Background(), 9999)
		if err == nil {
			t.Error("Expected error for non-existent user, got nil")
		}
	})
}

func TestUserRepository_FindBySessionToken(t *testing.T) {
	repo := NewUserRepository(testDB)

	// Setup test data
	testToken := "test_session_token_123"
	testUser := domain.User{
		SessionToken: testToken,
		Name:         "Session Token User",
		AvatarURL:    "session_avatar.png",
		ExpiresAt:    time.Now().Add(24 * time.Hour),
	}

	// Insert test user
	_, err := testDB.Exec(
		"INSERT INTO user_sessions(session_token, name, avatar_url, expires_at) VALUES ($1, $2, $3, $4)",
		testUser.SessionToken, testUser.Name, testUser.AvatarURL, testUser.ExpiresAt,
	)
	if err != nil {
		t.Fatalf("Failed to insert test user: %v", err)
	}

	t.Run("success", func(t *testing.T) {
		user, err := repo.FindBySessionToken(context.Background(), testToken)
		if err != nil {
			t.Fatalf("FindBySessionToken failed: %v", err)
		}

		if user.SessionToken != testToken || user.Name != testUser.Name {
			t.Errorf("Expected user with token %s and name %s, got token %s and name %s",
				testToken, testUser.Name, user.SessionToken, user.Name)
		}
	})

	t.Run("not found", func(t *testing.T) {
		_, err := repo.FindBySessionToken(context.Background(), "non_existent_token")
		if err == nil {
			t.Error("Expected error for non-existent token, got nil")
		}
	})
}

func TestUserRepository_GetUserIDBySessionToken(t *testing.T) {
	repo := NewUserRepository(testDB)

	// Setup test data
	testToken := "get_user_id_token"
	_, err := testDB.Exec(
		"INSERT INTO user_sessions(session_token, name) VALUES ($1, $2)",
		testToken, "Get User ID Test",
	)
	if err != nil {
		t.Fatalf("Failed to insert test user: %v", err)
	}

	t.Run("success", func(t *testing.T) {
		userID, err := repo.GetUserIDBySessionToken(context.Background(), testToken)
		if err != nil {
			t.Fatalf("GetUserIDBySessionToken failed: %v", err)
		}

		if userID <= 0 {
			t.Errorf("Expected positive user ID, got %d", userID)
		}
	})

	t.Run("not found", func(t *testing.T) {
		_, err := repo.GetUserIDBySessionToken(context.Background(), "non_existent_token")
		if err == nil {
			t.Error("Expected error for non-existent token, got nil")
		}
	})
}


func TestUserRepository_UpdateName(t *testing.T) {
	repo := NewUserRepository(testDB)

	// Setup test data
	testUser := domain.User{
		SessionToken: "update_name_token",
		Name:         "Original Name",
	}
	var userID int
	err := testDB.QueryRow(
		"INSERT INTO user_sessions(session_token, name) VALUES ($1, $2) RETURNING id",
		testUser.SessionToken, testUser.Name,
	).Scan(&userID)
	if err != nil {
		t.Fatalf("Failed to insert test user: %v", err)
	}

	t.Run("success", func(t *testing.T) {
		newName := "Updated Name"
		err := repo.UpdateName(context.Background(), userID, newName)
		if err != nil {
			t.Fatalf("UpdateName failed: %v", err)
		}

		// Verify the name was updated
		var currentName string
		err = testDB.QueryRow(
			"SELECT name FROM user_sessions WHERE id = $1",
			userID,
		).Scan(&currentName)
		if err != nil {
			t.Fatalf("Failed to verify updated name: %v", err)
		}

		if currentName != newName {
			t.Errorf("Expected name %s, got %s", newName, currentName)
		}
	})

	t.Run("user not found", func(t *testing.T) {
		err := repo.UpdateName(context.Background(), 9999, "New Name")
		if err == nil {
			t.Error("Expected error for non-existent user, got nil")
		}
	})
}