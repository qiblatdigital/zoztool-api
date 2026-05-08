package repository

import (
	"errors"
	"time"

	"github.com/qiblatdigital/zoztool-api/internal/model"
	"gorm.io/gorm"
)

type PostRepository struct {
	db *gorm.DB
}

func NewPostRepository(db *gorm.DB) *PostRepository {
	return &PostRepository{db: db}
}

func (r *PostRepository) Create(post *model.Post) error {
	return r.db.Create(post).Error
}

func (r *PostRepository) FindByID(id string) (*model.Post, error) {
	var post model.Post
	err := r.db.Where("id = ? AND deleted_at IS NULL", id).First(&post).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return &post, nil
}

func (r *PostRepository) FindByUserID(userID string, page, limit int) ([]model.Post, int64, error) {
	var posts []model.Post
	var total int64
	offset := (page - 1) * limit

	q := r.db.Model(&model.Post{}).Where("user_id = ? AND deleted_at IS NULL", userID)
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	err := q.Order("created_at DESC").Offset(offset).Limit(limit).Find(&posts).Error
	return posts, total, err
}

func (r *PostRepository) Update(id string, updates map[string]interface{}) error {
	updates["updated_at"] = time.Now()
	return r.db.Model(&model.Post{}).Where("id = ?", id).Updates(updates).Error
}

func (r *PostRepository) SoftDelete(id string) error {
	return r.db.Model(&model.Post{}).
		Where("id = ?", id).
		Update("deleted_at", gorm.Expr("NOW()")).Error
}

// FindScheduledDue returns posts with status=scheduled and scheduled_at <= now.
func (r *PostRepository) FindScheduledDue() ([]model.Post, error) {
	var posts []model.Post
	err := r.db.Where("status = ? AND scheduled_at <= ? AND deleted_at IS NULL",
		model.PostStatusScheduled, time.Now()).Find(&posts).Error
	return posts, err
}
