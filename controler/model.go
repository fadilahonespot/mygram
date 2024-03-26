package controler

import "time"

type Response struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
}

type RegisterRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
	Age      int    `json:"age" binding:"required"`
	Email    string `json:"email" binding:"required"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type LoginResponse struct {
	Token string `json:"token"`
}

type UpdateUserRequest struct {
	Username string `json:"username" binding:"required"`
	Email    string `json:"email" binding:"required"`
}

type UpdateUserResponse struct {
	ID        int       `json:"id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	Age       int       `json:"age"`
	UpdatedAt time.Time `json:"updated_at"`
}

type CreatePhotoRequest struct {
	Title    string `json:"title" binding:"required"`
	Caption  string `json:"caption"`
	PhotoUrl string `json:"photo_url" binding:"required"`
}

type CreatePhotoResponse struct {
	ID        int       `json:"id"`
	Title     string    `json:"title"`
	Caption   string    `json:"caption"`
	PhotoUrl  string    `json:"photo_url"`
	UserID    int       `json:"user_id"`
	CreatedAt time.Time `json:"created_at"`
}

type GetPhotoResponse struct {
	ID        int       `json:"id"`
	Title     string    `json:"title"`
	Caption   string    `json:"caption"`
	PhotoUrl  string    `json:"photo_url"`
	UserID    int       `json:"user_id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	User      UserData  `json:"user"`
}

type UserData struct {
	Email    string `json:"email"`
	Username string `json:"username"`
}

type UpdatePhotoResponse struct {
	ID        int       `json:"id"`
	Title     string    `json:"title"`
	Caption   string    `json:"caption"`
	PhotoUrl  string    `json:"photo_url"`
	UserID    int       `json:"user_id"`
	UpdatedAt time.Time `json:"updated_at"`
}

type AddCommentRequest struct {
	Message string `json:"message" binding:"required"`
	PhotoId int    `json:"photo_id"`
}

type AddCommentResponse struct {
	ID        int       `json:"id"`
	Message   string    `json:"message"`
	PhotoID   int       `json:"photo_id"`
	UserID    int       `json:"user_id"`
	CreatedAt time.Time `json:"created_at"`
}

type GetListCommentResponse struct {
	ID        int             `json:"id"`
	Message   string          `json:"message"`
	PhotoID   int             `json:"photo_id"`
	UserID    int             `json:"user_id"`
	CreatedAt time.Time       `json:"created_at"`
	UpdatedAt time.Time       `json:"updated_at"`
	User      UserDataComment `json:"user"`
	Photo     PhotoComment    `json:"photo"`
}

type UserDataComment struct {
	Id       int    `json:"id"`
	Email    string `json:"email"`
	Username string `json:"username"`
}

type PhotoComment struct {
	Id       int    `json:"id"`
	Title    string `json:"title"`
	Caption  string `json:"caption"`
	PhotoUrl string `json:"photo_url"`
	UserId   int    `json:"user_id"`
}

type UpdateCommentRequest struct {
	Message string `json:"message" binding:"required"`
}

type UpdateCommentResponse struct {
	Id        int       `json:"id"`
	Message   string    `json:"message"`
	PhotoID   int       `json:"photoId"`
	UserId    int       `json:"userId"`
	UpdatedAt time.Time `json:"updated_at"`
}

type CreateSocialMediaRequest struct {
	Name           string `json:"name" binding:"required"`
	SocialMediaUrl string `json:"social_media_url" binding:"required"`
}

type CreateSocialMediaResponse struct {
	ID             int       `json:"id"`
	Name           string    `json:"name"`
	SocialMediaUrl string    `json:"social_media_url"`
	UserID         int       `json:"user_id"`
	CreatedAt      time.Time `json:"created_at"`
}

type GetSocialMediaResponse struct {
	ID             int        `json:"id"`
	Name           string     `json:"name"`
	SocialMediaUrl string     `json:"social_media_url"`
	UserID         int        `json:"user_id"`
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at"`
	User           UserSocial `json:"user"`
}

type UserSocial struct {
	Id              int    `json:"id"`
	Username        string `json:"username"`
	ProfileImageUrl string `json:"profile_image_url"`
}

type UpdateSocialMediaRequest struct {
	Name           string `json:"name" binding:"required"`
	SocialMediaUrl string `json:"social_media_url" binding:"required"`
}

type UpdateSocialMediaResponse struct {
	ID             int       `json:"id"`
	Name           string    `json:"name"`
	SocialMediaUrl string    `json:"social_media_url"`
	UserID         int       `json:"user_id"`
	UpdatedAt      time.Time `json:"updated_at"`
}
