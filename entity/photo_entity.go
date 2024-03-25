package entity

import (
	"time"
)

type Photo struct {
	ID        int `gorm:"primaryKey"`
	Title     string
	Caption   string
	PhotoUrl  string
	UserID    int
	User      User
	CreatedAt time.Time
	UpdatedAt time.Time
	Comment   []Comment
}