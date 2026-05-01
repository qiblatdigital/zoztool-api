package model

// PostTargetStatus represents the publishing status of a post target.
type PostTargetStatus string

const (
	PostTargetStatusPending   PostTargetStatus = "pending"
	PostTargetStatusPublished PostTargetStatus = "published"
	PostTargetStatusFailed    PostTargetStatus = "failed"
)

type PostTarget struct {
	BaseModel
	PostID          string           `json:"post_id" gorm:"type:uuid;not null"`
	SocialAccountID string           `json:"social_account_id" gorm:"type:uuid;not null"`
	Status          PostTargetStatus `json:"status" gorm:"type:varchar(20);default:pending"`
	ErrorMessage    *string          `json:"error_message,omitempty" gorm:"type:text"`
}

func (p *PostTarget) TableName() string {
	return "post_targets"
}
