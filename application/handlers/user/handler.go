package user

import (
	"matiuskm/go-hotel-be/domain/entities"
	"matiuskm/go-hotel-be/infrastructure/etc"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func RegisterRoutes(r *gin.RouterGroup, db *gorm.DB) {
	r.POST("/", func(c *gin.Context) {
		var req struct {
			Username string `json:"username"`
			Password string `json:"password"`
			FullName string `json:"full_name"`
			Role string `json:"role"`
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid input"})
			return
		}
		
		hash, _ := etc.HashPassword(req.Password)
		user := entities.User {
			Username: req.Username,
			Password: hash,
			FullName: req.FullName,
			Role: req.Role,
		}
		if err := db.Create(&user).Error; err!= nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "failed to create user"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "user created"})
	})

	r.PUT("/:id", func (c *gin.Context) {
		var req struct {
			FullName string `json:"full_name"`
			Role string `json:"role"`
		}
		idParam := c.Param("id")
		if err := c.ShouldBindJSON(&req); err!= nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid input"})
			return
		}
		var user entities.User
		if err := db.First(&user, idParam).Error; err!= nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
			return
		}
		user.FullName = req.FullName
		user.Role = strings.ToLower(req.Role)
		if err := db.Save(&user).Error; err!= nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "failed to update user"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "user updated"})
	})

	r.PUT("/me", func(c *gin.Context) {
		uid, ok := c.Get("user_id")
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			return
		}
		var req struct {
			FullName string `json:"full_name"`
			Password string `json:"password"`
		}
		if err := c.ShouldBindJSON(&req); err!= nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid input"})
			return
		}
		var user entities.User
		if err := db.First(&user, uid).Error; err!= nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
			return
		}
		user.FullName = req.FullName
		if req.Password != "" {
			hash, _ := etc.HashPassword(req.Password)
			user.Password = hash
		}
		if err := db.Save(&user).Error; err!= nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "failed to update user"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "profile updated"})
	})

	r.DELETE("/:id", func(c *gin.Context) {
		role, _ := c.Get("role")
		if role!= "admin" {
			c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
			return
		}
		idParam := c.Param("id")
		var user entities.User
		if err := db.First(&user, idParam).Error; err!= nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
			return
		}
		if err := db.Delete(&user).Error; err!= nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "failed to delete user"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "user deleted"})
	})
}