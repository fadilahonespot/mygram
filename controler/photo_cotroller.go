package controler

import (
	"net/http"

	"github.com/dgrijalva/jwt-go"
	"github.com/fadilahonespot/mygram/entity"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type defaultPhoto struct {
	db *gorm.DB
}

func NewPhotoController(db *gorm.DB) *defaultPhoto {
	return &defaultPhoto{db}
}

func (s *defaultPhoto) CreatePhoto(c *gin.Context) {
	var req CreatePhotoRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response := generateResponseError(http.StatusBadRequest, err.Error())
		c.JSON(http.StatusBadRequest, response)
		return
	}

	token, _ := c.Get("token")
	claims := token.(*jwt.Token).Claims.(jwt.MapClaims)

	email := claims["email"].(string)
	var reqUser entity.User
	s.db.Take(&reqUser, "email = ?", email)
	if reqUser.Email == "" {
		response := generateResponseError(http.StatusNotFound, "Email not found")
		c.JSON(http.StatusBadRequest, response)
		return
	}

	var photo = entity.Photo{
		Title:    req.Title,
		Caption:  req.Caption,
		PhotoUrl: req.PhotoUrl,
		UserID:   reqUser.ID,
	}
	err := s.db.Create(&photo).Error
	if err != nil {
		response := generateResponseError(http.StatusInternalServerError, err.Error())
		c.JSON(http.StatusInternalServerError, response)
		return
	}

	respData := CreatePhotoResponse{
		ID:        photo.ID,
		Title:     photo.Title,
		Caption:   photo.Caption,
		PhotoUrl:  photo.PhotoUrl,
		UserID:    photo.UserID,
		CreatedAt: photo.CreatedAt,
	}

	response := generateResponseSuccess(http.StatusOK, "Photo Created", respData)
	c.JSON(http.StatusOK, response)
}

func (s *defaultPhoto) GetPhotos(c *gin.Context) {
	var photos []entity.Photo
	err := s.db.Preload("User").Find(&photos).Error
	if err != nil {
		response := generateResponseError(http.StatusInternalServerError, err.Error())
		c.JSON(http.StatusInternalServerError, response)
		return
	}

	var respData []GetPhotoResponse
	for _, photo := range photos {
		respData = append(respData, GetPhotoResponse{
			ID:        photo.ID,
			Title:     photo.Title,
			Caption:   photo.Caption,
			PhotoUrl:  photo.PhotoUrl,
			UserID:    photo.UserID,
			CreatedAt: photo.CreatedAt,
			UpdatedAt: photo.UpdatedAt,
			User: UserData{
				Email:    photo.User.Email,
				Username: photo.User.Username,
			},
		})
	}

	response := generateResponseSuccess(http.StatusOK, "Get Photos Success", respData)
	c.JSON(http.StatusOK, response)
}

func (s *defaultPhoto) UpdatePhoto(c *gin.Context) {
	var req CreatePhotoRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response := generateResponseError(http.StatusBadRequest, err.Error())
		c.JSON(http.StatusBadRequest, response)
		return
	}

	token, _ := c.Get("token")
	claims := token.(*jwt.Token).Claims.(jwt.MapClaims)
	email := claims["email"].(string)

	var reqUser entity.User
	s.db.Take(&reqUser, "email = ?", email)
	if reqUser.Email == "" {
		response := generateResponseError(http.StatusNotFound, "Email not found")
		c.JSON(http.StatusBadRequest, response)
		return
	}

	photoID := c.Param("photoId")
	var photo entity.Photo
	err := s.db.Take(&photo, "id = ?", photoID).Error
	if err != nil {
		response := generateResponseError(http.StatusNotFound, "Photo not found")
		c.JSON(http.StatusNotFound, response)
		return
	}

	if reqUser.ID != photo.UserID {
		response := generateResponseError(http.StatusUnauthorized, "photos belong to other users")
		c.JSON(http.StatusUnauthorized, response)
		return
	}

	photo.Title = req.Title
	photo.Caption = req.Caption
	photo.PhotoUrl = req.PhotoUrl

	err = s.db.Save(&photo).Error
	if err != nil {
		response := generateResponseError(http.StatusInternalServerError, err.Error())
		c.JSON(http.StatusInternalServerError, response)
		return
	}

	respData := UpdatePhotoResponse{
		ID:        photo.ID,
		Title:     photo.Title,
		Caption:   photo.Caption,
		PhotoUrl:  photo.PhotoUrl,
		UserID:    photo.UserID,
		UpdatedAt: photo.UpdatedAt,
	}
	response := generateResponseSuccess(http.StatusOK, "Photo Updated", respData)
	c.JSON(http.StatusOK, response)
}

func (s *defaultPhoto) DeletePhoto(c *gin.Context) {
	token, _ := c.Get("token")
	claims := token.(*jwt.Token).Claims.(jwt.MapClaims)
	email := claims["email"].(string)

	var reqUser entity.User
	s.db.Take(&reqUser, "email = ?", email)
	if reqUser.Email == "" {
		response := generateResponseError(http.StatusNotFound, "Email not found")
		c.JSON(http.StatusBadRequest, response)
		return
	}

	photoID := c.Param("photoId")
	var photo entity.Photo
	err := s.db.Take(&photo, "id = ?", photoID).Error
	if err != nil {
		response := generateResponseError(http.StatusNotFound, "Photo not found")
		c.JSON(http.StatusNotFound, response)
		return
	}

	if reqUser.ID != photo.UserID {
		response := generateResponseError(http.StatusUnauthorized, "photos belong to other users")
		c.JSON(http.StatusUnauthorized, response)
	}
	err = s.db.Delete(&photo).Error
	if err != nil {
		response := generateResponseError(http.StatusInternalServerError, err.Error())
		c.JSON(http.StatusInternalServerError, response)
		return
	}

	response := generateResponseSuccess(http.StatusOK, "Your photo has been successfully deleted", nil)
	c.JSON(http.StatusOK, response)
}
