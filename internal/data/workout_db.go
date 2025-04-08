package database

import (
	"context"
	"go.mongodb.org/mongo-driver/v2/bson"
	"time"
)

func (DB *DB) GetAllWorkouts(userId bson.ObjectID) (map[string]interface{}, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	find, err := DB.Db.Collection("workouts").Find(ctx, bson.M{"user_id": userId})
	if err != nil {
		return nil, err
	}
	result := map[string]interface{}{}
	err = find.All(ctx, &result)
	if err != nil {
		return nil, err
	}
	return result, nil
}
func (DB *DB) UpdateWorkout(userID bson.ObjectID, updatedWorkout map[string]interface{}) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	filter := bson.M{"user_id": userID}
	update := bson.M{"$set": updatedWorkout}

	_, err := DB.Db.Collection("workouts").UpdateOne(ctx, filter, update);
	return err
}
func (DB *DB) DeleteWorkout(userId bson.ObjectID, workoutId string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	filter := bson.M{"user_id": userId, "workout": bson.M{"id": workoutId}}
	_, err := DB.Db.Collection("workouts").DeleteOne(ctx, filter)
	return err
}
func (DB *DB) SaveWorkout(userId bson.ObjectID, workout map[string]interface{}) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	toInsert := bson.M{
		"user_id": userId,
		"workout": workout,
	}
	_, err := DB.Db.Collection("workouts").InsertOne(ctx, toInsert)
	if err != nil {
		return err
	}
	return nil
}

/*
{
	"name":"Trening",
	"timeLong":"00:46:35",
	"distance":30850,
	"workoutDate":"2025-03-28T09:38:20.776Z",
	"poolLength":25,
	"mainType":["Wydolność"],
	"elementsIn":[
		{
			"name":"Nazwa cwiczenia",
			"type":"exercise",
			"id":"459bf7d2-0fc0-437c-9ba5-ce408609499b",
			"distance":30850,
			"time":"34:00",
			"subtype":{
				"label":"Wydolność",
				"value":"wydolnosc"
			},
			"equipment":[
				"monofin",
				"handpaddles"
			]
		},
		{
			"type":"break",
			"id":"badc09cf-c41a-4659-838e-a8ad278f491b",
			"name":"Przerwa",
			"time":"12:35"
		}
	]
*/