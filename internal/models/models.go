package models

import (
	"go.mongodb.org/mongo-driver/v2/bson"
)

type (
	User struct {
		Id           bson.ObjectID `json:"id" bson:"_id"`
		Username     string        `json:"username,omitempty" bson:"username,omitempty"`
		Email        string        `json:"email,omitempty" bson:"email,omitempty"`
		Password     string        `json:"password,omitempty" bson:"password,omitempty"`
		AccessToken  string        `json:"access_token,omitempty" bson:"-"`
		RefreshToken string        `json:"refresh_token,omitempty" bson:"-"`
	}
)

var (
	NilUser = User{}
)