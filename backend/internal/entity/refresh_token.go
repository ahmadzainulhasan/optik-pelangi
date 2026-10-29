package entity

import "time"

type RefreshToken struct {
	ID        string    `gorm:"primaryKey;column:id;type:uuid"`
	UserID    string    `gorm:"column:user_id;type:uuid;not null;index"`
	Token     string    `gorm:"column:token;type:uuid;unique;not null;index"`
	ExpiresAt time.Time `gorm:"column:expires_at;not null"`
	CreatedAt time.Time `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt time.Time `gorm:"column:updated_at;autoUpdateTime"`
}

func (refreshToken *RefreshToken) TableName() string {
	return "refresh_tokens"
}
