package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"week4-webserver/database"
	"week4-webserver/handlers"
	"week4-webserver/middleware"
)

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found, using system environment variables")
	}
}

func main() {
	mongodbURI := os.Getenv("MONGODB_URI")
	if mongodbURI == "" {
		mongodbURI = "mongodb://localhost:27017"
	}

	dbName := os.Getenv("DATABASE_NAME")
	if dbName == "" {
		dbName = "userdb"
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	database.InitMongoDB(mongodbURI, "userdb")

	r := gin.Default()

	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	auth := r.Group("/")
	auth.Use(middleware.AuthMiddleware())
	{
		auth.GET("/user", handlers.GetUser)
		auth.PUT("/user", handlers.UpdateUser)
		auth.PUT("/password", handlers.ChangePassword)
		auth.POST("/logout", handlers.Logout)
	}

	r.POST("/register", handlers.Register)
	r.POST("/login", handlers.Login)

	log.Println("Server running on :8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
