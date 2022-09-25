package controller

import (
	"context"
	"log"
	"net/http"
	"soccer-api/app/model"
	"soccer-api/configuration"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var teamCollection *mongo.Collection = configuration.Collection("team")
var validateTeam = validator.New()

func CreateTeam() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		var team model.Team
		defer cancel()

		if err := c.BindJSON(&team); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"data":    bson.M{},
				"error":   true,
				"message": err,
			})
			return
		}

		if validationErr := validateTeam.Struct(&team); validationErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"data":    bson.M{},
				"error":   true,
				"message": validationErr.Error(),
			})
			return
		}

		data := model.Team{
			ID:        primitive.NewObjectID(),
			Name:      team.Name,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		result, err := teamCollection.InsertOne(ctx, data)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"data":    bson.M{},
				"error":   true,
				"message": err,
			})
			return
		}

		c.JSON(http.StatusCreated, gin.H{
			"data":    result,
			"error":   false,
			"message": "Data berhasil ditambahkan!",
		})
		return
	}
}

func GetTeam() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		teamId := c.Param("id")
		var team model.Team
		defer cancel()

		objId, _ := primitive.ObjectIDFromHex(teamId)

		err := teamCollection.FindOne(ctx, bson.M{"_id": objId}).Decode(&team)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{
				"data":    bson.M{},
				"error":   true,
				"message": "Team with specified ID not found!",
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"data":    team,
			"error":   false,
			"message": http.StatusOK,
		})
		return
	}
}

func GetTeamWithPlayer() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		teamId := c.Param("id")
		var team model.Team
		defer cancel()

		objId, _ := primitive.ObjectIDFromHex(teamId)

		err := teamCollection.FindOne(ctx, bson.M{"_id": objId}).Decode(&team)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{
				"data":    bson.M{},
				"error":   true,
				"message": "Team with specified ID not found!",
			})
			return
		}
		team.WithPlayer()

		c.JSON(http.StatusOK, gin.H{
			"data":    team,
			"error":   false,
			"message": http.StatusOK,
		})
		return
	}
}

func UpdateTeam() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		teamId := c.Param("id")
		var team model.Team
		defer cancel()

		objId, _ := primitive.ObjectIDFromHex(teamId)

		if err := c.BindJSON(&team); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"data":    bson.M{},
				"error":   true,
				"message": err,
			})
			return
		}

		if validationErr := validateTeam.Struct(&team); validationErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"data":    bson.M{},
				"error":   true,
				"message": validationErr.Error(),
			})
			return
		}

		update := bson.M{"name": team.Name, "updated_at": time.Now()}
		_, err := teamCollection.UpdateOne(ctx, bson.M{"_id": objId}, bson.M{"$set": update})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"data":    bson.M{},
				"error":   true,
				"message": err,
			})
			return
		}

		var result model.Team
		err = teamCollection.FindOne(ctx, bson.M{"_id": objId}).Decode(&result)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{
				"data":    bson.M{},
				"error":   true,
				"message": "Team with specified ID not found!",
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"data":    result,
			"error":   false,
			"message": "Data berhasil diubah!",
		})
		return
	}
}

func DeleteTeam() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		teamId := c.Param("id")
		defer cancel()

		objId, _ := primitive.ObjectIDFromHex(teamId)

		result, err := teamCollection.DeleteOne(ctx, bson.M{"_id": objId})

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"data":    bson.M{},
				"error":   true,
				"message": err,
			})
			return
		}

		if result.DeletedCount < 1 {
			c.JSON(http.StatusNotFound, gin.H{
				"data":    bson.M{},
				"error":   true,
				"message": "Team with specified ID not found!",
			})
			return
		}

		c.JSON(http.StatusNoContent, gin.H{
			"data":    bson.M{},
			"error":   true,
			"message": "Team berhasil dihapus!",
		})
	}
}

func GetAllTeam() gin.HandlerFunc {
	return func(c *gin.Context) {
		var team []*model.Team

		findOptions := options.Find()
		findOptions.SetLimit(20)
		cur, err := teamCollection.Find(context.TODO(), bson.D{{}}, findOptions)
		if err != nil {
			log.Fatal(err)
		}
		for cur.Next(context.TODO()) {
			var elem *model.Team
			err := cur.Decode(&elem)
			if err != nil {
				log.Fatal(err)
			}
			team = append(team, elem)
		}

		c.JSON(http.StatusOK, gin.H{
			"data":    &team,
			"error":   false,
			"message": http.StatusOK,
		})
		return
	}
}
