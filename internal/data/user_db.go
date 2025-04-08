package database

import (
	"context"
	"errors"
	"github.com/labstack/echo/v4"
"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"net/http"
"testProject/internal/models"
	"time"
)

type DB struct {
	Db *mongo.Database
}

func (DB *DB) GetUser(id bson.ObjectID) (models.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	result := DB.Db.Collection("users").FindOne(ctx, bson.D{{"_id", id}})
	if err := result.Err(); err != nil { return models.NilUser, err }

	var user models.User
	if err := result.Decode(&user); err != nil { return models.NilUser, err }

	return user, nil
}
func (DB *DB) GetUserByName(username string) (models.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	result := DB.Db.Collection("users").FindOne(ctx, bson.D{{"username", username}})
	if err := result.Err(); err != nil { return models.NilUser, err }
	var user models.User
	if err := result.Decode(&user); err != nil { return models.NilUser, err }
	return user, nil
}
func (DB *DB) Exists(username string) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    filter := bson.M{
		"username": username,
    }

	var result models.User
	err := DB.Db.Collection("users").FindOne(ctx, filter).Decode(&result)
	if errors.Is(err, mongo.ErrNoDocuments) {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return true, nil
}
func (DB *DB) ExistsID(id bson.ObjectID) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	filter := bson.D{{"_id", id}}
	var result models.User

	err := DB.Db.Collection("users").FindOne(ctx, filter).Decode(&result)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}
func (DB *DB) UpdateAccountInfo(id bson.ObjectID, weight int, isMale bool, caloriesGoal int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	filter := bson.M{"_id":id};
	updated := bson.M{"$set": bson.M{ "weight":weight, "isMale": isMale, "caloriesGoal":caloriesGoal }};

	result := DB.Db.Collection("users").FindOneAndUpdate(ctx, filter, updated)
	if errors.Is(result.Err(), mongo.ErrNoDocuments) {
		return &echo.HTTPError{Code: http.StatusNotFound, Message: "user not found"}
	}
	return result.Err()
}
func (DB *DB) NewUser(id bson.ObjectID, username string, password string, weight int, isMale bool, caloriesGoal int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	user := models.User{Id: id, Username: username, Password: password, Weight: weight, IsMale: isMale, CaloriesGoal: caloriesGoal};

	_, err := DB.Db.Collection("users").InsertOne(ctx, user)
	if err != nil { return err }

	return nil
}