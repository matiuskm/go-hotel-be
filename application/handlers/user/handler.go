package user

import (
	"matiuskm/go-hotel-be/application/middlewares"
	"matiuskm/go-hotel-be/domain/entities"
	"matiuskm/go-hotel-be/infrastructure/etc"
	"matiuskm/go-hotel-be/pkg/payloads"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
)

func RegisterRoutes(r *gin.RouterGroup, db *gorm.DB) {
	// register new user (admin)
	r.POST("/", middlewares.RequireRoles("admin"), func(c *gin.Context) {
		validate := validator.New()
		var req payloads.RegisterUserRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid input"})
			return
		}
		req.Role = strings.ToLower(req.Role)

		if err := validate.Struct(req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
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

	// edit user (admin)
	r.PUT("/:id", middlewares.RequireRoles("admin"), func (c *gin.Context) {
		validate := validator.New()
		var req = payloads.EditUserRequest{}

		if err := c.ShouldBindJSON(&req); err!= nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid input"})
			return
		}

		if err := validate.Struct(req); err!= nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		idParam := c.Param("id")
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

	// edit user profile
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

	// delete user (admin)
	r.DELETE("/:id", middlewares.RequireRoles("admin"), func(c *gin.Context) {
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

	// assign role (admin)
	r.PUT("/:id/role", middlewares.RequireRoles("admin"), func(c *gin.Context) {
		idParam := c.Param("id")
		validate := validator.New()
		var req = payloads.AssignRoleRequest{}
		if err := c.ShouldBindJSON(&req); err!= nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid input"})
			return
		}
		if err := validate.Struct(req); err!= nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		var user entities.User
		if err := db.First(&user, idParam).Error; err!= nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
			return
		}
		user.Role = strings.ToLower(req.Role)
		if err := db.Save(&user).Error; err!= nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "failed to update user"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "user role updated"})
	})

	// get user list (admin)
	r.GET("/", middlewares.RequireRoles("admin"), func(c *gin.Context) {
		var users []entities.User
		if err := db.Find(&users).Error; err!= nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "failed to get users"})
			return
		}
		var resp []payloads.UserResponse
		for _, u := range users {
			resp = append(resp, payloads.UserResponse{
				ID: u.ID,
				Username: u.Username,
				FullName: u.FullName,
				Role: u.Role,
			})
		}

		c.JSON(http.StatusOK, resp)
	})

	// get user by id
	r.GET("/:id", middlewares.RequireRoles("admin"), func(c *gin.Context) {
		idParam := c.Param("id")
		var user entities.User
		if err := db.First(&user, idParam).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
			return
		}

		resp := payloads.UserResponse{
			ID: user.ID,
			Username: user.Username,
			FullName: user.FullName,
			Role: user.Role,
		}

		c.JSON(http.StatusOK, resp)
	})

	// get user profile
	r.GET("/me", func(c *gin.Context) {
		uid, ok := c.Get("user_id")
		if!ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			return
		}
		var user entities.User
		if err := db.First(&user, uid).Error; err!= nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
			return
		}
		resp := payloads.UserResponse{
			ID: user.ID,
			Username: user.Username,
			FullName: user.FullName,
			Role: user.Role,
		}
		c.JSON(http.StatusOK, resp)
	})
}