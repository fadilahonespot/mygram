package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/dgrijalva/jwt-go"
	"github.com/fadilahonespot/mygram/controler"
	"github.com/fadilahonespot/mygram/pkg/database"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()
	db := database.ConnectDB()

	userController := controler.NewUserController(db)
	photoControler := controler.NewPhotoController(db)
	commentController := controler.NewCommentController(db)
	socialMediaControler := controler.NewSocialMediaController(db)

	router := gin.Default()
	router.POST("/users/register", userController.Register)
	router.POST("/users/login", userController.Login)

	tx := router.Use(authMiddleware)
	tx.PUT("/users", userController.UpdateUser)
	tx.DELETE("/users", userController.DeleteUser)
	tx.POST("/photos", photoControler.CreatePhoto)
	tx.GET("/photos", photoControler.GetPhotos)
	tx.PUT("/photos/:photoId", photoControler.UpdatePhoto)
	tx.DELETE("/photos/:photoId", photoControler.DeletePhoto)
	tx.POST("/comments", commentController.CreateComment)
	tx.GET("/comments", commentController.GetListComment)
	tx.PUT("/comments/:commentId", commentController.UpdateComment)
	tx.DELETE("/comments/:commentId", commentController.DeleteComment)
	tx.POST("/socialmedias", socialMediaControler.CreateSocialMedia)
	tx.GET("/socialmedias", socialMediaControler.GetSocialMedia)
	tx.PUT("/socialmedias/:socialMediaId", socialMediaControler.UpdateSocialMedia)
	tx.DELETE("/socialmedias/:socialMediaId", socialMediaControler.DeleteSocialMedia)

	router.Run(fmt.Sprintf(":%v", os.Getenv("PORT")))
}

func authMiddleware(c *gin.Context) {
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		c.Abort()
		return
	}

	tokenString := authHeader[len("Bearer "):]

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(os.Getenv("SECRET_KEY")), nil
	})
	if err != nil || !token.Valid {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		c.Abort()
		return
	}

	c.Set("token", token)
	c.Next()
}
