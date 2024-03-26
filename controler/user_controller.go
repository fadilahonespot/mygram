package controler

import (
	"net/http"

	"github.com/fadilahonespot/mygram/entity"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"github.com/dgrijalva/jwt-go"
)

type defautUser struct {
	db *gorm.DB
}

func NewUserController(db *gorm.DB) *defautUser {
	return &defautUser{db}
}

func (s *defautUser) Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response := generateResponseError(http.StatusBadRequest, err.Error())
		c.JSON(400, response)
		return
	}

	pass, _ := hashPassword(req.Password)
	reqUser := entity.User{
		Username: req.Username,
		Email:    req.Email,
		Password: pass,
		Age:      req.Age,
	}

	var userData entity.User
	s.db.Take(&userData, "email = ?", reqUser.Email)
	if userData.Email != "" {
		response := generateResponseError(http.StatusBadRequest, "Email already exists")
		c.JSON(http.StatusBadRequest, response)
		return
	}

	err := s.db.Create(&reqUser).Error
	if err != nil {
		response := generateResponseError(http.StatusInternalServerError, err.Error())
		c.JSON(http.StatusInternalServerError, response)
		return
	}

	req.Password = userData.Password
	response := generateResponseSuccess(http.StatusCreated, "Register Success", req)
	c.JSON(http.StatusCreated, response)
}

func (s *defautUser) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response := generateResponseError(http.StatusBadRequest, err.Error())
		c.JSON(400, response)
		return
	}

	var reqUser entity.User
	s.db.Take(&reqUser, "email = ?", req.Email)
	if reqUser.Email == "" {
		response := generateResponseError(http.StatusNotFound, "Email not found")
		c.JSON(http.StatusNotFound, response)
		return
	}

	err := comparePassword(reqUser.Password, req.Password)
	if err != nil {
		response := generateResponseError(http.StatusInternalServerError, err.Error())
		c.JSON(http.StatusBadRequest, response)
		return
	}

	var dataResp LoginResponse
	token, err := createToken(reqUser.Username, reqUser.Email)
	if err != nil {
		response := generateResponseError(http.StatusInternalServerError, err.Error())
		c.JSON(http.StatusInternalServerError, response)
		return
	}
	dataResp.Token = token

	response := generateResponseSuccess(http.StatusOK, "Login Success", dataResp)
	c.JSON(http.StatusOK, response)
}

func (s *defautUser) UpdateUser(c *gin.Context) {
	token, _ := c.Get("token")
	claims := token.(*jwt.Token).Claims.(jwt.MapClaims)
	username := claims["username"].(string)

	var req UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response := generateResponseError(http.StatusBadRequest, err.Error())
		c.JSON(400, response)
		return
	}

	var reqUser entity.User
	s.db.Take(&reqUser, "username = ?", username)
	if reqUser.ID == 0 {
		response := generateResponseError(http.StatusNotFound, "User not found")
		c.JSON(http.StatusBadRequest, response)
		return
	}

	reqUser.Username = req.Username
	reqUser.Email = req.Email
	err := s.db.Model(&reqUser).Updates(&reqUser).Error
	if err != nil {
		response := generateResponseError(http.StatusInternalServerError, err.Error())
		c.JSON(http.StatusInternalServerError, response)
		return
	}

	var respData = UpdateUserResponse{
		ID:        reqUser.ID,
		Username:  reqUser.Username,
		Email:     reqUser.Email,
		Age:       reqUser.Age,
		UpdatedAt: reqUser.UpdatedAt,
	}
	response := generateResponseSuccess(http.StatusOK, "Update Success", respData)
	c.JSON(http.StatusOK, response)
}

func (s *defautUser) DeleteUser(c *gin.Context) {
	token, _ := c.Get("token")
	claims := token.(*jwt.Token).Claims.(jwt.MapClaims)

	email := claims["email"].(string)
	var reqUser entity.User
	s.db.Take(&reqUser, "email = ?", email)
	if reqUser.ID == 0 {
		response := generateResponseError(http.StatusNotFound, "User not found")
		c.JSON(http.StatusBadRequest, response)
		return
	}

	err := s.db.Delete(&reqUser).Error
	if err != nil {
		response := generateResponseError(http.StatusInternalServerError, err.Error())
		c.JSON(http.StatusInternalServerError, response)
		return
	}

	response := generateResponseSuccess(http.StatusOK, "You account has been succesfully deleted", nil)
	c.JSON(http.StatusOK, response)
}
