package services

import (
	"1337b04rd/internal/domain"
	"context"
	"errors"
	"testing"
)

type mockRickAndMortyAPI struct {
	getRandomCharacterErr error
	name                  string
	avatarURL             string
}

func (m *mockRickAndMortyAPI) GetRandomCharacter(ctx context.Context) (string, string, error) {
	if m.getRandomCharacterErr != nil {
		return "", "", m.getRandomCharacterErr
	}
	return m.name, m.avatarURL, nil
}

func TestUserService_GetOrCreateUser(t *testing.T) {
	tests := []struct {
		name             string
		sessionToken     string
		findByTokenErr   error
		apiErr           error
		saveErr          error
		expectedIsNew    bool
		expectedErr      bool
		expectSessionGen bool
	}{
		{
			name:             "existing user",
			sessionToken:     "existing_token",
			findByTokenErr:   nil,
			expectedIsNew:    false,
			expectedErr:      false,
			expectSessionGen: false,
		},
		{
			name:             "new user with existing token",
			sessionToken:     "existing_token",
			findByTokenErr:   errors.New("not found"),
			expectedIsNew:    true,
			expectedErr:      false,
			expectSessionGen: false,
		},
		{
			name:             "new user with empty token",
			sessionToken:     "",
			expectedIsNew:    true,
			expectedErr:      false,
			expectSessionGen: true,
		},
		{
			name:             "api error",
			sessionToken:     "token",
			findByTokenErr:   errors.New("not found"),
			apiErr:           errors.New("api error"),
			expectedIsNew:    true,
			expectedErr:      true,
			expectSessionGen: false,
		},
		{
			name:             "save error",
			sessionToken:     "token",
			findByTokenErr:   errors.New("not found"),
			saveErr:          errors.New("save error"),
			expectedIsNew:    true,
			expectedErr:      true,
			expectSessionGen: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			userRepo := newMockUserRepo()
			userRepo.findByTokenErr = tt.findByTokenErr
			userRepo.saveErr = tt.saveErr

			api := &mockRickAndMortyAPI{
				name:                  "Test User",
				avatarURL:             "http://example.com/avatar.jpg",
				getRandomCharacterErr: tt.apiErr,
			}

			service := NewUserService(userRepo, api)

			// Pre-populate for existing user case
			if tt.sessionToken != "" && tt.findByTokenErr == nil {
				userRepo.Save(context.Background(), &domain.User{
					SessionToken: tt.sessionToken,
					Name:         "Existing User",
					AvatarURL:    "http://example.com/existing.jpg",
				})
			}

			user, isNew, err := service.GetOrCreateUser(context.Background(), tt.sessionToken)

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

			if isNew != tt.expectedIsNew {
				t.Errorf("expected isNew=%v, got %v", tt.expectedIsNew, isNew)
			}

			if tt.expectSessionGen {
				if user.SessionToken == "" {
					t.Error("expected session token to be generated")
				}
				if user.SessionToken == tt.sessionToken {
					t.Error("expected new session token to be different from input")
				}
			} else if tt.sessionToken != "" && user.SessionToken != tt.sessionToken {
				t.Errorf("expected session token %s, got %s", tt.sessionToken, user.SessionToken)
			}

			if user.Name == "" {
				t.Error("expected user name to be set")
			}
			if user.AvatarURL == "" {
				t.Error("expected avatar URL to be set")
			}
		})
	}
}

func TestUserService_GetUserByID(t *testing.T) {
	tests := []struct {
		name        string
		userID      int
		prepopulate bool
		findByIDErr error
		expectedErr bool
	}{
		{
			name:        "existing user",
			userID:      1,
			prepopulate: true,
			findByIDErr: nil,
			expectedErr: false,
		},
		{
			name:        "non-existent user",
			userID:      2,
			prepopulate: false,
			findByIDErr: nil,
			expectedErr: true,
		},
		{
			name:        "repository error",
			userID:      1,
			prepopulate: true,
			findByIDErr: errors.New("db error"),
			expectedErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			userRepo := newMockUserRepo()
			userRepo.findByIDErr = tt.findByIDErr

			if tt.prepopulate {
				userRepo.Save(context.Background(), &domain.User{ID: tt.userID})
			}

			service := NewUserService(userRepo, &mockRickAndMortyAPI{})
			_, err := service.GetUserByID(context.Background(), tt.userID)

			if tt.expectedErr {
				if err == nil {
					t.Error("expected error, got nil")
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
			}
		})
	}
}

func TestUserService_UpdateUserName(t *testing.T) {
	tests := []struct {
		name          string
		userID        int
		newName       string
		prepopulate   bool
		updateNameErr error
		expectedErr   bool
	}{
		{
			name:          "successful update",
			userID:        1,
			newName:       "New Name",
			prepopulate:   true,
			updateNameErr: nil,
			expectedErr:   false,
		},
		{
			name:          "non-existent user",
			userID:        2,
			newName:       "New Name",
			prepopulate:   false,
			updateNameErr: nil,
			expectedErr:   true,
		},
		{
			name:          "repository error",
			userID:        1,
			newName:       "New Name",
			prepopulate:   true,
			updateNameErr: errors.New("update error"),
			expectedErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			userRepo := newMockUserRepo()
			userRepo.updateNameErr = tt.updateNameErr

			if tt.prepopulate {
				userRepo.Save(context.Background(), &domain.User{
					ID:   tt.userID,
					Name: "Old Name",
				})
			}

			service := NewUserService(userRepo, &mockRickAndMortyAPI{})
			err := service.UpdateUserName(context.Background(), tt.userID, tt.newName)

			if tt.expectedErr {
				if err == nil {
					t.Error("expected error, got nil")
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}

				// Verify the name was updated
				user, _ := userRepo.FindByID(context.Background(), tt.userID)
				if user.Name != tt.newName {
					t.Errorf("expected name %s, got %s", tt.newName, user.Name)
				}
			}
		})
	}
}
