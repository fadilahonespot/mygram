package entity

import "time"

type SocialMedia struct {
	ID             int `gorm:"primaryKey"`
	Name           string
	SocialMediaUrl string
	User           User
	UserID         int
	CreatedAt      time.Time
	UpdatedAt      time.Time
}
