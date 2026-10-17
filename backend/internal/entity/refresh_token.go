package entity

import "time"

type RefreshToken struct {
	ID        string    `gorm:"primaryKey;column:id;type:varchar(36)"`
	UserID    string    `gorm:"column:user_id;type:varchar(36);not null;index"`
	Token     string    `gorm:"column:token;type:varchar(255);unique;not null;index"`
	ExpiresAt time.Time `gorm:"column:expires_at;not null"`
	CreatedAt time.Time `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt time.Time `gorm:"column:updated_at;autoUpdateTime"`
}

func (refreshToken *RefreshToken) TableName() string {
	return "refresh_tokens"
}
