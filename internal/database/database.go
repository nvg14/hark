package database

import (
	"context"
	"log"

	"github.com/spf13/viper"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func NewDatabase() (*mongo.Client, error) {
	clientOptions := options.Client().ApplyURI(viper.GetString("mongodb.URI"))

	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	err = client.Ping(context.TODO(), nil)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	log.Println("Connected to MongoDB!")
	return client, nil
}
