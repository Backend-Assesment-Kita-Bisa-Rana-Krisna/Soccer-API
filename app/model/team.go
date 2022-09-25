package model

import (
	"context"
	"soccer-api/configuration"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type Team struct {
	ID        primitive.ObjectID `json:"_id" bson:"_id"`
	Name      string             `json:"name" bson:"name"`
	Players   []Player           `json:"players,omitempty"`
	CreatedAt time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt time.Time          `json:"updated_at" bson:"updated_at"`
}

func (u *Team) WithPlayer() *Team {
	var playerCollection *mongo.Collection = configuration.Collection("player")
	cur, _ := playerCollection.Find(context.TODO(), bson.M{"team_id": u.ID.Hex()})
	for cur.Next(context.TODO()) {
		var player Player
		_ = cur.Decode(&player)
		u.Players = append(u.Players, player)
	}

	return u
}
