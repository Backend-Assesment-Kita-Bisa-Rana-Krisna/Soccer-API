package configuration

import (
	"os"
	"context"
	"fmt"
	"log"
	"time"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var Client *mongo.Client

func ConnectDB() *mongo.Client {
	if Client != nil {
		return Client
	}

	err := godotenv.Load(".env")

	if err != nil {
		log.Fatal("Error loading .env file")
	}

	MongoDb := os.Getenv("MONGODB_URL")

	Client, err = mongo.NewClient(options.Client().ApplyURI(MongoDb))
	if err != nil {
		log.Fatal(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)

	defer cancel()
	err = Client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}

	err = Client.Ping(ctx, nil)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Connected to MongoDB.")
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
	var collection *mongo.Collection = client.Database(os.Getenv("MONGODB_DATABASE")).Collection(collectionName)
	return collection
}
