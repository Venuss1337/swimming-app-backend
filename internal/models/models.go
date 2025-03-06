package models

import "go.mongodb.org/mongo-driver/v2/bson"

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Email    string `json:"email"`
}

type AuthUser struct {
	ID           bson.ObjectID `bson:"_id"`
	Username     string        `bson:"username"`
	PasswordHash string        `bson:"passwordHash"`
}

type LoginReponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	UserID       string `json:"user_id"`
}
