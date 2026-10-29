package entity

import "time"

type User struct {
	ID        string    `gorm:"primaryKey;column:id;type:uuid"`
	Name      string    `gorm:"column:name;type:text;not null"`
	Email     string    `gorm:"column:email;type:text;unique;not null"`
	GoogleID  string    `gorm:"column:google_id;type:text;unique"`
	Picture   string    `gorm:"column:picture;type:text"`
	CreatedAt time.Time `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt time.Time `gorm:"column:updated_at;autoUpdateTime"`
}

func (user *User) TableName() string {
	return "users"
}