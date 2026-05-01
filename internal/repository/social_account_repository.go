package repository

import (
	"github.com/qiblatdigital/zoztool-api/internal/model"
	"gorm.io/gorm"
)

type SocialAccountRepository struct {
	db *gorm.DB
}

func NewSocialAccountRepository(db *gorm.DB) *SocialAccountRepository {
	return &SocialAccountRepository{db: db}
}

func (r *SocialAccountRepository) Create(account *model.SocialAccount) error {
	return r.db.Create(account).Error
}

func (r *SocialAccountRepository) FindByID(id string) (*model.SocialAccount, error) {
	var account model.SocialAccount
	err := r.db.Where("id = ?", id).First(&account).Error
	if err != nil {
		return nil, err
	}
	return &account, nil
}
