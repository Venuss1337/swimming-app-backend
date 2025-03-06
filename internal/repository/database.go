package repository

import (
	"context"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"testProject/internal/models"
	"time"
)

type Database struct {
	client   *mongo.Client
	database *mongo.Database
}

var LocalService *Database = &Database{}

// Disabbled TLS for testing
func Connect(uri string) error {
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

	LocalService.client = client
	LocalService.database = client.Database("workout-app")

	return nil
}

func (d *Database) GetUserAndPassword(username string) (models.AuthUser, error) {
	result := d.database.Collection("users").FindOne(context.Background(), bson.M{"username": username})
	if result.Err() != nil {
		return models.AuthUser{}, result.Err()
	}
	var user models.AuthUser
	err := result.Decode(&user)
	if err != nil {
		return models.AuthUser{}, err
	}

	return user, nil
}

func (d *Database) Close() {
	if d.client != nil {
		_ = d.client.Disconnect(context.Background())
	}
}
