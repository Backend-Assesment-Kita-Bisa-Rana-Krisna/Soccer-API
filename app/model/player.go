package model

import (
	"context"
	"soccer-api/configuration"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type Player struct {
	ID        primitive.ObjectID `json:"_id" bson:"_id"`
	Name      string             `json:"name" bson:"name" validate:"required"`
	TeamId    string             `json:"team_id" bson:"team_id" validate:"required"`
	Team      *Team              `json:"team,omitempty"`
	CreatedAt time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt time.Time          `json:"updated_at" bson:"updated_at"`
}

func (u *Player) WithTeam() *Player {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	objId, _ := primitive.ObjectIDFromHex(u.TeamId)
	var teamCollection *mongo.Collection = configuration.Collection("team")
	_ = teamCollection.FindOne(ctx, bson.M{"_id": objId}).Decode(&u.Team)
	defer cancel()
	return u
}
