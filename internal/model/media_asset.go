package model

type MediaAsset struct {
	BaseModel
	UserID       string  `json:"user_id" gorm:"type:uuid;not null"`
	PostID       *string `json:"post_id" gorm:"type:uuid"`
	Filename     string  `json:"filename" gorm:"type:varchar(255);not null"`
	OriginalName string  `json:"original_name" gorm:"type:varchar(255);not null"`
	MimeType     string  `json:"mime_type" gorm:"type:varchar(100);not null"`
	Size         int64   `json:"size" gorm:"not null"`
	URL          string  `json:"url" gorm:"type:text;not null"`
}

func (m *MediaAsset) TableName() string {
	return "media_assets"
}
