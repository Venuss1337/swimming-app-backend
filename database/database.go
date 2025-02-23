package database

import (
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"time"
)

type Database struct {
	client   *mongo.Client
	database *mongo.Database
}

// Disabbled TLS for debugging
func (d *Database) Connect(uri string) error {
	/*tlsConfig := &tls.Config{
		InsecureSkipVerify: true,
	}*/

	clientOptions := options.Client().ApplyURI(uri).
		/*SetTLSConfig(tlsConfig)*/
		SetConnectTimeout(time.Second * 10)

	client, err := mongo.Connect(clientOptions)
	if err != nil {
		return err
	}

	if err := client.Ping(context.Background(), nil); err != nil {
		return err
	}

	d.client = client
	d.database = client.Database("workout-app")
	return nil
}

func (d *Database) Close() {
	if d.client != nil {
		_ = d.client.Disconnect(context.Background())
	}
}

func (d *Database) ContainsUser(username string) bool {
	collection := d.database.Collection("users")
	filter := bson.D{{"username", username}}
	result := collection.FindOne(context.Background(), filter)
	return result.Err() == nil
}

func (d *Database) RegisterUser(username string, password string) error {
	if d.ContainsUser(username) {
		return errors.New("user already exists")
	}
	_, err := d.database.Collection("users").InsertOne(context.Background(), bson.M{"username": username, "password": password})
	if err != nil {
		return err
	}
	return nil
}
func (d *Database) RetrievePasswordHashAndId(username string) (string, bson.ObjectID, error) {
	dcm := d.database.Collection("users").FindOne(context.Background(), bson.M{"username": username})

	if dcm.Err() != nil {
		return "", bson.NilObjectID, dcm.Err()
	}

	var result struct {
		Id           bson.ObjectID `bson:"_id"`
		PasswordHash string        `bson:"password"`
	}

	if err := dcm.Decode(&result); err != nil {
		return "", bson.NilObjectID, errors.New("error decoding user")
	}
	return result.PasswordHash, result.Id, nil
}
