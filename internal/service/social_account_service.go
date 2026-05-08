package service

import (
	"errors"

	"github.com/qiblatdigital/zoztool-api/internal/model"
	"github.com/qiblatdigital/zoztool-api/internal/repository"
)

type SocialAccountService struct {
	repo *repository.SocialAccountRepository
}

func NewSocialAccountService(repo *repository.SocialAccountRepository) *SocialAccountService {
	return &SocialAccountService{repo: repo}
}

type CreateSocialAccountRequest struct {
	Platform      string `json:"platform" binding:"required"`
	Username      string `json:"username" binding:"required"`
	AccessToken   string `json:"access_token" binding:"required"`
	ThreadsUserID string `json:"threads_user_id"`
}

func (s *SocialAccountService) Create(userID string, req CreateSocialAccountRequest) (*model.SocialAccount, error) {
	account := &model.SocialAccount{
		UserID:        userID,
		Platform:      req.Platform,
		Username:      req.Username,
		AccessToken:   req.AccessToken,
		ThreadsUserID: req.ThreadsUserID,
	}
	if err := s.repo.Create(account); err != nil {
		return nil, errors.New("failed to create social account")
	}
	return account, nil
}

func (s *SocialAccountService) List(userID string) ([]model.SocialAccount, error) {
	return s.repo.FindByUserID(userID)
}

func (s *SocialAccountService) GetByID(userID, id string) (*model.SocialAccount, error) {
	account, err := s.repo.FindByID(id)
	if err != nil {
		return nil, errors.New("social account not found")
	}
	if account.UserID != userID {
		return nil, errors.New("social account not found")
	}
	return account, nil
}

func (s *SocialAccountService) Delete(userID, id string) error {
	account, err := s.repo.FindByID(id)
	if err != nil {
		return errors.New("social account not found")
	}
	if account.UserID != userID {
		return errors.New("social account not found")
	}
	return s.repo.SoftDelete(id)
}
