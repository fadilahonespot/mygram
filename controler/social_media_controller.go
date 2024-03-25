package controler

import (
	"net/http"

	"github.com/dgrijalva/jwt-go"
	"github.com/fadilahonespot/mygram/entity"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type defaultSocialMedia struct {
	db *gorm.DB
}

func NewSocialMediaController(db *gorm.DB) *defaultSocialMedia {
	return &defaultSocialMedia{db}
}

func (s *defaultSocialMedia) CreateSocialMedia(c *gin.Context) {
	var req CreateSocialMediaRequest
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

	mediaData := entity.SocialMedia{
		Name:           req.Name,
		SocialMediaUrl: req.SocialMediaUrl,
		UserID:         reqUser.ID,
	}
	err := s.db.Create(&mediaData).Error
	if err != nil {
		response := generateResponseError(http.StatusInternalServerError, err.Error())
		c.JSON(http.StatusInternalServerError, response)
		return
	}

	respData := CreateSocialMediaResponse{
		ID:             mediaData.ID,
		Name:           mediaData.Name,
		SocialMediaUrl: mediaData.SocialMediaUrl,
		UserID:         mediaData.UserID,
		CreatedAt:      mediaData.CreatedAt,
	}

	response := generateResponseSuccess(http.StatusOK, "Success Create Social Media", respData)
	c.JSON(http.StatusOK, response)
}

func (s *defaultSocialMedia) GetSocialMedia(c *gin.Context) {
	var socialMedia []entity.SocialMedia
	err := s.db.Preload("User").Find(&socialMedia).Error
	if err != nil {
		response := generateResponseError(http.StatusInternalServerError, err.Error())
		c.JSON(http.StatusInternalServerError, response)
		return
	}

	var data []GetSocialMediaResponse
	for _, social := range socialMedia {
		var photo string
		if len(social.User.Photo) != 0 {
			photo = social.User.Photo[0].PhotoUrl
		}
		data = append(data, GetSocialMediaResponse{
			ID:             social.ID,
			Name:           social.Name,
			SocialMediaUrl: social.SocialMediaUrl,
			UserID:         social.UserID,
			CreatedAt:      social.CreatedAt,
			UpdatedAt:      social.UpdatedAt,
			User: UserSocial{
				Id:              social.UserID,
				Username:        social.User.Username,
				ProfileImageUrl: photo,
			},
		})
	}

	response := generateResponseSuccess(http.StatusOK, "Success Get Social Media", data)
	c.JSON(http.StatusOK, response)
}

func (s *defaultSocialMedia) UpdateSocialMedia(c *gin.Context) {
	var req UpdateSocialMediaRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response := generateResponseError(http.StatusBadRequest, err.Error())
		c.JSON(http.StatusBadRequest, response)
		return
	}

	socialMediaId := c.Param("socialMediaId")

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

	var socialMedia entity.SocialMedia
	err := s.db.Take(&socialMedia, "id = ?", socialMediaId).Error
	if err != nil {
		response := generateResponseError(http.StatusNotFound, "Social Media not found")
		c.JSON(http.StatusNotFound, response)
		return
	}

	if reqUser.ID != socialMedia.UserID {
		response := generateResponseError(http.StatusBadRequest, "You are not allowed to delete this social media")
		c.JSON(http.StatusBadRequest, response)
		return
	}
	socialMedia.Name = req.Name
	socialMedia.SocialMediaUrl = req.SocialMediaUrl
	err = s.db.Save(&socialMedia).Error
	if err != nil {
		response := generateResponseError(http.StatusInternalServerError, err.Error())
		c.JSON(http.StatusInternalServerError, response)
		return
	}

	respData := UpdateSocialMediaResponse{
		ID:             socialMedia.ID,
		Name:           socialMedia.Name,
		SocialMediaUrl: socialMedia.SocialMediaUrl,
		UserID:         socialMedia.UserID,
		UpdatedAt:      socialMedia.UpdatedAt,
	}

	response := generateResponseSuccess(http.StatusOK, "Success Update Social Media", respData)
	c.JSON(http.StatusOK, response)
}

func (s *defaultSocialMedia) DeleteSocialMedia(c *gin.Context) {
	socialMediaId := c.Param("socialMediaId")

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

	var socialMedia entity.SocialMedia
	err := s.db.Take(&socialMedia, "id = ?", socialMediaId).Error
	if err != nil {
		response := generateResponseError(http.StatusNotFound, "Social Media not found")
		c.JSON(http.StatusNotFound, response)
		return
	}

	if reqUser.ID != socialMedia.UserID {
		response := generateResponseError(http.StatusBadRequest, "You are not allowed to delete this social media")
		c.JSON(http.StatusBadRequest, response)
		return
	}

	err = s.db.Delete(&socialMedia).Error
	if err != nil {
		response := generateResponseError(http.StatusInternalServerError, "failed delete social media")
		c.JSON(http.StatusInternalServerError, response)
		return
	}

	response := generateResponseSuccess(http.StatusOK, "Your social media has been successfully deleted", nil)
	c.JSON(http.StatusOK, response)
}
