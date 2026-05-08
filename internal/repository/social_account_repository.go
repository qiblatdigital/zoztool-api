package repository

import (
	"errors"

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
	err := r.db.Where("id = ? AND deleted_at IS NULL", id).First(&account).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return &account, nil
}

func (r *SocialAccountRepository) FindByUserID(userID string) ([]model.SocialAccount, error) {
	var accounts []model.SocialAccount
	err := r.db.Where("user_id = ? AND deleted_at IS NULL", userID).Find(&accounts).Error
	return accounts, err
}

func (r *SocialAccountRepository) SoftDelete(id string) error {
	return r.db.Model(&model.SocialAccount{}).
		Where("id = ?", id).
		Update("deleted_at", gorm.Expr("NOW()")).Error
}
