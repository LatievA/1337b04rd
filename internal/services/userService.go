package services

import (
	"1337b04rd/internal/domain"
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
)

type UserService struct {
	userRepo        domain.UserRepository
	rickAndMortyAPI domain.RickAndMortyAPI
}

func NewUserService(userRepo domain.UserRepository, api domain.RickAndMortyAPI) domain.UserService {
	return &UserService{userRepo: userRepo, rickAndMortyAPI: api}
}

func (s *UserService) GetOrCreateUser(ctx context.Context, sessionToken string) (*domain.User, bool, error) {
	var isNew bool
	if sessionToken == "" {
		newSessionToken, err := s.generateSession()
		isNew = true
		if err != nil {
			return nil, isNew, fmt.Errorf("failed to create new session token")
		}
		sessionToken = newSessionToken
	} else {
		user, err := s.userRepo.FindBySessionToken(ctx, sessionToken)
		if err == nil {
			return user, isNew, nil
		}
		isNew = true
	}

	name, avatarURL, err := s.rickAndMortyAPI.GetRandomCharacter(ctx)
	if err != nil {
		return nil, isNew, err
	}

	newUser := &domain.User{
		Name:         name,
		AvatarURL:    avatarURL,
		SessionToken: sessionToken,
	}
	id, err := s.userRepo.Save(ctx, newUser)
	if err != nil {
		return nil, isNew, err
	}
	newUser.ID = id
	return newUser, isNew, nil
}

func (s *UserService) GetUserByID(ctx context.Context, userID int) (*domain.User, error) {
	user, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("user not found")
	}
	return user, nil
}

func (s *UserService) UpdateUserName(ctx context.Context, userID int, newName string) error {
	return s.userRepo.UpdateName(ctx, userID, newName)
}

func (s *UserService) generateSession() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}

	return hex.EncodeToString(b), nil
}
