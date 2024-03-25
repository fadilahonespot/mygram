package entity

import "time"

type User struct {
	ID          int `gorm:"primaryKey"`
	Username    string
	Email       string
	Password    string
	Age         int
	CreatedAt   time.Time
	UpdatedAt   time.Time
	Photo       []Photo
	Comment     []Comment
	SocialMedia []SocialMedia
}