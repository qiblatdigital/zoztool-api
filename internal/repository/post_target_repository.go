package repository

import (
	"errors"
	"time"

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
	err := r.db.Where("post_id = ? AND deleted_at IS NULL", postID).Find(&targets).Error
	return targets, err
}

func (r *PostTargetRepository) FindByID(id string) (*model.PostTarget, error) {
	var target model.PostTarget
	err := r.db.Where("id = ? AND deleted_at IS NULL", id).First(&target).Error
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
		"status":     status,
		"updated_at": time.Now(),
	}
	if errMsg != nil {
		updates["error_message"] = *errMsg
	}
	return r.db.Model(&model.PostTarget{}).Where("id = ?", id).Updates(updates).Error
}
