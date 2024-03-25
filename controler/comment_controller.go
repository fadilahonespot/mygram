package controler

import (
	"net/http"

	"github.com/dgrijalva/jwt-go"
	"github.com/fadilahonespot/mygram/entity"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type defaultComment struct {
	db *gorm.DB
}

func NewCommentController(db *gorm.DB) *defaultComment {
	return &defaultComment{db: db}
}

func (s *defaultComment) CreateComment(c *gin.Context) {
	var req AddCommentRequest
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

	var photoData entity.Photo
	err := s.db.Take(&photoData, "id = ?", req.PhotoId).Error
	if err != nil {
		response := generateResponseError(http.StatusNotFound, "Photo not found")
		c.JSON(http.StatusNotFound, response)
		return
	}

	reqPhoto := &entity.Comment{
		Message:   req.Message,
		PhotoID:   req.PhotoId,
		UserID:    reqUser.ID,
	}
	err = s.db.Create(&reqPhoto).Error
	if err != nil {
		response := generateResponseError(http.StatusInternalServerError, err.Error())
		c.JSON(http.StatusInternalServerError, response)
		return
	}

	respData := AddCommentResponse{
		ID:        reqPhoto.ID,
		Message:   reqPhoto.Message,
		PhotoID:   reqPhoto.PhotoID,
		UserID:    reqPhoto.UserID,
		CreatedAt: reqPhoto.CreatedAt,
	}
	response := generateResponseSuccess(http.StatusOK, "Commen added", respData)
	c.JSON(http.StatusOK, response)
}

func (s *defaultComment) GetListComment(c *gin.Context) {
	var comments []entity.Comment
	err := s.db.Preload("User").Preload("Photo").Find(&comments).Error
	if err != nil {
		response := generateResponseError(http.StatusInternalServerError, err.Error())
		c.JSON(http.StatusInternalServerError, response)
		return
	}

	var respData []GetListCommentResponse
	for _, comment := range comments {
		respData = append(respData, GetListCommentResponse{
			ID:        comment.ID,
			Message:   comment.Message,
			PhotoID:   comment.PhotoID,
			UserID:    comment.UserID,
			CreatedAt: comment.CreatedAt,
			UpdatedAt: comment.UpdatedAt,
			User:      UserDataComment{
				Id:       comment.UserID,
				Email:    comment.User.Email,
				Username: comment.User.Username,
			},
			Photo:     PhotoComment{
				Id:       comment.PhotoID,
				Title:    comment.Photo.Title,
				Caption:  comment.Photo.Caption,
				PhotoUrl: comment.Photo.PhotoUrl,
				UserId:   comment.Photo.UserID,
			},
		})
	}

	
	response := generateResponseSuccess(http.StatusOK, "List comment", respData)
	c.JSON(http.StatusOK, response)
}

func (s *defaultComment) UpdateComment(c *gin.Context) {
	var req UpdateCommentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response := generateResponseError(http.StatusBadRequest, err.Error())
		c.JSON(http.StatusBadRequest, response)
		return
	}

	commentId := c.Param("commentId")

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

	var commentData entity.Comment
	err := s.db.Take(&commentData, "id = ?", commentId).Error
	if err != nil {
		response := generateResponseError(http.StatusNotFound, "Comment not found")
		c.JSON(http.StatusNotFound, response)
		return
	}

	if reqUser.ID != commentData.UserID {
		response := generateResponseError(http.StatusForbidden, "You are not allowed to update this comment")
		c.JSON(http.StatusForbidden, response)
		return
	}

	commentData.Message = req.Message
	err = s.db.Save(commentData).Error
	if err != nil {
		response := generateResponseError(http.StatusInternalServerError, err.Error())
		c.JSON(http.StatusInternalServerError, response)
		return
	}

	respData := UpdateCommentResponse{
		Id:        commentData.ID,
		Message:   commentData.Message,
		PhotoID:   commentData.PhotoID,
		UserId:    commentData.UserID,
		UpdatedAt: commentData.UpdatedAt,
	}

	response := generateResponseSuccess(http.StatusOK, "Comment updated", respData)
	c.JSON(http.StatusOK, response)
}

func (s *defaultComment) DeleteComment(c *gin.Context) {
	commentId := c.Param("commentId")

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

	var commentData entity.Comment
	err := s.db.Take(&commentData, "id = ?", commentId).Error
	if err != nil {
		response := generateResponseError(http.StatusNotFound, "Comment not found")
		c.JSON(http.StatusNotFound, response)
		return
	}

	if reqUser.ID != commentData.UserID {
		response := generateResponseError(http.StatusForbidden, "You are not allowed to delete this comment")
		c.JSON(http.StatusForbidden, response)
		return
	}

	err = s.db.Delete(&commentData).Error
	if err != nil {
		response := generateResponseError(http.StatusInternalServerError, err.Error())
		c.JSON(http.StatusInternalServerError, response)
		return
	}
	response := generateResponseSuccess(http.StatusOK, "Your comment has succesfully deleted", nil)
	c.JSON(http.StatusOK, response)
}