package repository

import (
	"errors"

	"github.com/qiblatdigital/zoztool-api/internal/model"
	"gorm.io/gorm"
)

type PostTargetRepository struct {
	db *gorm.DB
}

func NewPostTargetRepository(db *gorm.DB) *PostTargetRepository {
	return &PostTargetRepository{db: db}
}

func (r *PostTargetRepository) CreateBatch(targets []model.PostTarget) error {
	return r.db.Create(&targets).Error
}

func (r *PostTargetRepository) FindByPostID(postID string) ([]model.PostTarget, error) {
	var targets []model.PostTarget
	err := r.db.Where("post_id = ?", postID).Limit(100).Find(&targets).Error
	return targets, err
}

func (r *PostTargetRepository) FindByID(id string) (*model.PostTarget, error) {
	var target model.PostTarget
	err := r.db.Where("id = ?", id).First(&target).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return &target, nil
}

func (r *PostTargetRepository) UpdateStatus(id string, status model.PostTargetStatus, errMsg *string) error {
	updates := map[string]interface{}{
		"status": status,
	}
	if errMsg != nil {
		updates["error_message"] = *errMsg
	}
	result := r.db.Model(&model.PostTarget{}).Where("id = ?", id).Updates(updates)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return ErrNotFound
	}
	return nil
}
