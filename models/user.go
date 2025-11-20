package models

// User represents user data model
type User struct {
	Username string `bson:"username" json:"username" example:"AHuangMeow"`
	Nickname string `bson:"nickname" json:"nickname" example:"阿鍠"`
	Password string `bson:"password" json:"-"`
}

// RegisterRequest represents user registration request
type RegisterRequest struct {
	Username string `json:"username" example:"AHuangMeow" binding:"required"`
	Nickname string `json:"nickname" example:"阿鍠" binding:"required"`
	Password string `json:"password" example:"12345678" binding:"required"`
}

// LoginRequest represents user login request
type LoginRequest struct {
	Username string `json:"username" example:"AHuangMeow" binding:"required"`
	Password string `json:"password" example:"12345678" binding:"required"`
}

// UpdateUserRequest represents user update request
type UpdateUserRequest struct {
	Target  string `json:"target" example:"nickname" binding:"required" enums:"username,nickname"`
	Content string `json:"content" example:"怕酱" binding:"required"`
}

// ChangePasswordRequest represents password change request
type ChangePasswordRequest struct {
	OldPassword string `json:"old_password" example:"12345678" binding:"required"`
	NewPassword string `json:"new_password" example:"87654321" binding:"required"`
}

