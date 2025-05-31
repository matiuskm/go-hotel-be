package auth

import (
	"matiuskm/go-hotel-be/domain/entities"
	"matiuskm/go-hotel-be/infrastructure/etc"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"gorm.io/gorm"
)

func LoginHandler(db *gorm.DB) gin.HandlerFunc{
	return func(c *gin.Context) {
		var req struct {
			Username string `json:"username"`
			Password string `json:"password"`
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
			return
		}

		var u entities.User
		if err := db.Where("username = ?", req.Username).First(&u).Error; err!= nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
			return
		}
		if !etc.CheckPasswordHash(req.Password, u.Password) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
			return
		}

		token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"user_id": u.ID,
			"role": u.Role,
			"exp": time.Now().Add(time.Hour * 1).Unix(),
		})
		secret := []byte(os.Getenv("JWT_SECRET"))
		tokenString, err := token.SignedString(secret)
		if err!= nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"token": tokenString})
	}
}