package database

import (
	"context"
	"crypto/tls"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"time"
)

type Database struct {
	client   *mongo.Client
	database *mongo.Database
}

func (d *Database) Connect(uri string) error {
	tlsConfig := &tls.Config{
		InsecureSkipVerify: true,
	}

	clientOptions := options.Client().ApplyURI(uri).
		SetTLSConfig(tlsConfig).
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

func (d *Database) RegisterUser(username string, password string) error {
	_, err := d.database.Collection("users").InsertOne(context.Background(), bson.M{"username": username, "password": password})
	if err != nil {
		return err
	}
	return nil
}
func (d *Database) RetrievePasswordHash(username string) (string, error) {
	dcm, err := d.database.Collection("users").Find(context.Background(), bson.M{"username": username})
	if err != nil {
		return "", err
	}
	var result struct {
		PasswordHash string `bson:"password"`
	}

	defer dcm.Close(context.Background())
	if err := dcm.Decode(&result); err != nil {
		return "", err
	}
	return result.PasswordHash, nil
}
