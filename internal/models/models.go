package models

import (
	"go.mongodb.org/mongo-driver/v2/bson"
)

type (
	User struct {
		Id           bson.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
		Username     string        `json:"username,omitempty" bson:"username,omitempty"`
		Password     string        `json:"password,omitempty" bson:"password,omitempty"`
		Weight       int           `json:"weight,omitempty" bson:"weight,omitempty"`
		IsMale       bool          `json:"isMale,omitempty" bson:"isMale,omitempty"`
		CaloriesGoal int           `json:"caloriesGoal,omitempty" bson:"caloriesGoal,omitempty"`
		AccessToken  string        `json:"access_token,omitempty" bson:"-"`
		RefreshToken string        `json:"refresh_token,omitempty" bson:"-"`
	}
)

var (
	NilUser = User{}
)
