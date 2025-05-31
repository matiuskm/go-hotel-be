package main

import (
	"log"
	userHandler "matiuskm/go-hotel-be/application/handlers/user"
	authHandler "matiuskm/go-hotel-be/application/handlers/auth"
	"matiuskm/go-hotel-be/config"
	"matiuskm/go-hotel-be/infrastructure/database"
	"net/http"
	"os"
	"strings"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	config.LoadEnv()

	dbConn, err := database.ConnectAndMigrate()
	if err != nil {
		log.Fatal("Failed to connect to the database")
		panic(err)
	}

	r := gin.Default()

	originEnv := os.Getenv("CORS_ORIGINS")
	log.Println("Loaded CORS_ORIGINS:", originEnv)

	if originEnv == "" {
		log.Fatal("CORS_ORIGINS env var is required")
	}

	corsOrigins := strings.Split(originEnv, ",")
	r.Use(cors.New(cors.Config{
		AllowOrigins:  corsOrigins,
		AllowMethods:  []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:  []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders: []string{"Content-Length"},
	}))

	api := r.Group("/api")
	userHandler.RegisterRoutes(api.Group("/users"), dbConn)
	api.POST("/login", authHandler.LoginHandler(dbConn))
	
	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "Hotel Management System in Golang",
		})
	})

	r.Run(":" + os.Getenv("PORT"))
}