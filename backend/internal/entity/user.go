package entity

import "time"

type User struct {
	ID        string    `gorm:"primaryKey;column:id;type:varchar(36)"`
	Name      string    `gorm:"column:name;type:varchar(100);not null"`
	Email     string    `gorm:"column:email;type:varchar(100);unique;not null"`
	GoogleID  string    `gorm:"column:google_id;type:varchar(100);unique"`
	Picture   string    `gorm:"column:picture;type:varchar(255)"`
	CreatedAt time.Time `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt time.Time `gorm:"column:updated_at;autoUpdateTime"`
}

func (user *User) TableName() string {
	return "users"
}