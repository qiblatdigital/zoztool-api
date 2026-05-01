package model

import (
	"time"
)

type SocialAccount struct {
	BaseModel
	UserID         string     `json:"user_id" gorm:"type:uuid;not null"`
	Platform       string     `json:"platform" gorm:"type:varchar(50);not null"`
	AccessToken    string     `json:"-" gorm:"type:text"`
	RefreshToken   string     `json:"-" gorm:"type:text"`
	TokenExpiresAt *time.Time `json:"token_expires_at" gorm:"type:timestamp"`
	Username       string     `json:"username" gorm:"type:varchar(255)"`
}

func (s *SocialAccount) TableName() string {
	return "social_accounts"
}
