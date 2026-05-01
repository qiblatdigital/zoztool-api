package model

import (
	"time"
)

// PostStatus represents the lifecycle status of a post.
type PostStatus string

const (
	PostStatusDraft     PostStatus = "draft"
	PostStatusScheduled PostStatus = "scheduled"
	PostStatusPublished PostStatus = "published"
)

type Post struct {
	BaseModel
	UserID      string     `json:"user_id" gorm:"type:uuid;not null"`
	Title       string     `json:"title" gorm:"type:varchar(255)"`
	Content     string     `json:"content" gorm:"type:text"`
	Status      PostStatus `json:"status" gorm:"type:varchar(20);default:draft"`
	ScheduledAt *time.Time `json:"scheduled_at" gorm:"type:timestamp"`
	PublishedAt *time.Time `json:"published_at" gorm:"type:timestamp"`
}

func (p *Post) TableName() string {
	return "posts"
}
