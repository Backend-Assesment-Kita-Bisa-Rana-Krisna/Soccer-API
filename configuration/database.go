package configuration

import (
	"context"
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var Client *mongo.Client

func ConnectDB() *mongo.Client {
	if Client != nil {
		return Client
	}

	var err error
	clientOptions := options.Client().ApplyURI("mongodb://127.0.0.1:27017")
	Client, err = mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	err = Client.Ping(context.Background(), nil)
	if err != nil {
		log.Fatal(err)
	}

	return Client
}

func CloseConnectDB() {
	if Client == nil {
		return
	}

	err := Client.Disconnect(context.TODO())
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Connection to MongoDB closed.")
}

func Collection(collectionName string) *mongo.Collection {
	var client *mongo.Client = ConnectDB()
	var collection *mongo.Collection = client.Database("soccer-api").Collection(collectionName)
	return collection
}
