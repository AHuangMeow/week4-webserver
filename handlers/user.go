package handlers

import (
	"context"
	"time"

	"week4-webserver/database"
	"week4-webserver/middleware"
	"week4-webserver/models"
	"week4-webserver/utils"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func Register(c *gin.Context) {
	var req models.RegisterRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "Invalid request data")
		return
	}

	if req.Username == "" || req.Nickname == "" || req.Password == "" {
		utils.BadRequest(c, "Username, nickname and password are required")
		return
	}

	if len(req.Password) < 8 {
		utils.BadRequest(c, "Password must be at least 8 characters")
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var existingUser models.User
	err := database.UserCollection.FindOne(ctx, bson.M{"username": req.Username}).Decode(&existingUser)
	if err == nil {
		utils.BadRequest(c, "Username already exists")
		return
	} else if err != mongo.ErrNoDocuments {
		utils.InternalError(c, "Database error")
		return
	}

	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		utils.InternalError(c, "Failed to hash password")
		return
	}

	newUser := models.User{
		Username: req.Username,
		Nickname: req.Nickname,
		Password: hashedPassword,
	}

	_, err = database.UserCollection.InsertOne(ctx, newUser)
	if err != nil {
		utils.InternalError(c, "Failed to create user")
		return
	}

	utils.Success(c, gin.H{
		"username": newUser.Username,
		"nickname": newUser.Nickname,
		"message":  "User registered successfully",
	})
}

func Login(c *gin.Context) {
	var req models.LoginRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "Invalid request data")
		return
	}

	if req.Username == "" || req.Password == "" {
		utils.BadRequest(c, "Username and password are required")
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var user models.User
	err := database.UserCollection.FindOne(ctx, bson.M{"username": req.Username}).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			utils.Unauthorized(c, "Invalid username or password")
		} else {
			utils.InternalError(c, "Database error")
		}
		return
	}

	if !utils.CheckPasswordHash(req.Password, user.Password) {
		utils.Unauthorized(c, "Invalid username or password")
		return
	}

	token, err := utils.GenerateToken(user.Username)
	if err != nil {
		utils.InternalError(c, "Failed to generate token")
		return
	}

	utils.Success(c, gin.H{
		"token": token,
		"user": gin.H{
			"username": user.Username,
			"nickname": user.Nickname,
		},
		"message": "Login succeed",
	})
}

func GetUser(c *gin.Context) {
	username, exists := middleware.GetUsernameFromContext(c)
	if !exists {
		utils.Unauthorized(c, "User not found in context")
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var user models.User
	err := database.UserCollection.FindOne(ctx, bson.M{"username": username}).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			utils.NotFound(c, "User not found")
		} else {
			utils.InternalError(c, "Database error")
		}
		return
	}

	utils.Success(c, gin.H{
		"username": user.Username,
		"nickname": user.Nickname,
	})
}

func UpdateUser(c *gin.Context) {
	username, exists := middleware.GetUsernameFromContext(c)
	if !exists {
		utils.Unauthorized(c, "User not found in context")
		return
	}

	token, exists := middleware.GetTokenFromContext(c)
	if exists {
		utils.InvalidateToken(token)
	}

	var req models.UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "Invalid request data")
		return
	}

	if req.Target != "username" && req.Target != "nickname" {
		utils.BadRequest(c, "Invalid request data")
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	updateFields := bson.M{}
	if req.Target == "username" {
		var existingUser models.User
		err := database.UserCollection.FindOne(ctx, bson.M{
			"username": req.Content,
		}).Decode(&existingUser)
		if err == nil {
			utils.BadRequest(c, "Username already exists")
			return
		} else if err != mongo.ErrNoDocuments {
			utils.InternalError(c, "Database error")
			return
		}
		updateFields["username"] = req.Content
	} else {
		updateFields["nickname"] = req.Content
	}

	_, err := database.UserCollection.UpdateOne(
		ctx,
		bson.M{"username": username},
		bson.M{"$set": updateFields},
	)
	if err != nil {
		utils.InternalError(c, "Failed to update user")
		return
	}

	utils.Success(c, gin.H{
		"message": "User updated successfully",
	})
}

func ChangePassword(c *gin.Context) {
	username, exists := middleware.GetUsernameFromContext(c)
	if !exists {
		utils.Unauthorized(c, "User not found in context")
		return
	}

	token, exists := middleware.GetTokenFromContext(c)
	if exists {
		utils.InvalidateToken(token)
	}

	var req models.ChangePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "Invalid request data")
		return
	}

	if len(req.NewPassword) < 8 {
		utils.BadRequest(c, "New password must be at least 8 characters")
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var user models.User
	err := database.UserCollection.FindOne(ctx, bson.M{"username": username}).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			utils.NotFound(c, "User not found")
		} else {
			utils.InternalError(c, "Database error")
		}
		return
	}

	if !utils.CheckPasswordHash(req.OldPassword, user.Password) {
		utils.BadRequest(c, "Old password is incorrect")
		return
	}

	hashedPassword, err := utils.HashPassword(req.NewPassword)
	if err != nil {
		utils.InternalError(c, "Failed to hash password")
		return
	}

	_, err = database.UserCollection.UpdateOne(
		ctx,
		bson.M{"username": username},
		bson.M{"$set": bson.M{"password": hashedPassword}},
	)
	if err != nil {
		utils.InternalError(c, "Failed to update password")
		return
	}

	utils.Success(c, gin.H{
		"message": "Password updated successfully",
	})
}

func Logout(c *gin.Context) {
	token, exists := middleware.GetTokenFromContext(c)
	if !exists {
		utils.Success(c, gin.H{
			"message": "Logout succeed",
		})
		return
	}

	err := utils.InvalidateToken(token)
	if err != nil {
		utils.Success(c, gin.H{
			"message": "Logout succeed",
			"warning": "Token may still be valid for a short time",
		})
		return
	}

	utils.Success(c, gin.H{
		"message": "Logout succeed, token has been invalidated",
	})
}
