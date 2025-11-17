package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"week4-webserver/database"
	"week4-webserver/handlers"
	"week4-webserver/middleware"
)

func main() {
	database.InitMongoDB("mongodb://localhost:27017", "userdb")

	r := gin.Default()

	r.POST("/register", handlers.Register)
	r.POST("/login", handlers.Login)

	auth := r.Group("/")
	auth.Use(middleware.AuthMiddleware())
	{
		auth.GET("/user", handlers.GetUser)
		auth.PUT("/user", handlers.UpdateUser)
		auth.PUT("/password", handlers.ChangePassword)
		auth.POST("/logout", handlers.Logout)
	}

	log.Println("Server running on :8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
