package entity

import "time"

type Comment struct {
	ID        int `gorm:"primaryKey"`
	Message   string
	PhotoID   int
	Photo     Photo
	UserID    int
	User      User
	CreatedAt time.Time
	UpdatedAt time.Time
}
