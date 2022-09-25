package controller

import (
	"context"
	"fmt"
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

var playerCollection *mongo.Collection = configuration.Collection("player")
var validatePlayer = validator.New()

func CreatePlayer() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		var player model.Player
		defer cancel()

		if err := c.BindJSON(&player); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"data":    bson.M{},
				"error":   true,
				"message": err,
			})
			return
		}

		if validationErr := validatePlayer.Struct(&player); validationErr != nil {
			fmt.Println(validatePlayer)
			c.JSON(http.StatusBadRequest, gin.H{
				"data":    bson.M{},
				"error":   true,
				"message": validationErr.Error(),
			})
			return
		}

		data := model.Player{
			ID:        primitive.NewObjectID(),
			Name:      player.Name,
			TeamId:    player.TeamId,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		result, err := playerCollection.InsertOne(ctx, data)
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

func GetPlayer() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		playerId := c.Param("id")
		var player model.Player
		defer cancel()

		objId, _ := primitive.ObjectIDFromHex(playerId)

		err := playerCollection.FindOne(ctx, bson.M{"_id": objId}).Decode(&player)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{
				"data":    bson.M{},
				"error":   true,
				"message": "Player with specified ID not found!",
			})
			return
		}
		player.WithTeam()

		c.JSON(http.StatusOK, gin.H{
			"data":    player,
			"error":   false,
			"message": http.StatusOK,
		})
		return
	}
}

func UpdatePlayer() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		playerId := c.Param("id")
		var player model.Player
		defer cancel()

		objId, _ := primitive.ObjectIDFromHex(playerId)

		if err := c.BindJSON(&player); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"data":    bson.M{},
				"error":   true,
				"message": err,
			})
			return
		}

		if validationErr := validatePlayer.Struct(&player); validationErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"data":    bson.M{},
				"error":   true,
				"message": validationErr.Error(),
			})
			return
		}

		update := bson.M{"name": player.Name, "updated_at": time.Now()}
		_, err := playerCollection.UpdateOne(ctx, bson.M{"_id": objId}, bson.M{"$set": update})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"data":    bson.M{},
				"error":   true,
				"message": err,
			})
			return
		}

		var result model.Player
		err = playerCollection.FindOne(ctx, bson.M{"_id": objId}).Decode(&result)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{
				"data":    bson.M{},
				"error":   true,
				"message": "Player with specified ID not found!",
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

func DeletePlayer() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		playerId := c.Param("id")
		defer cancel()

		objId, _ := primitive.ObjectIDFromHex(playerId)

		result, err := playerCollection.DeleteOne(ctx, bson.M{"_id": objId})

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
				"message": "Player with specified ID not found!",
			})
			return
		}

		c.JSON(http.StatusNoContent, gin.H{
			"data":    bson.M{},
			"error":   true,
			"message": "Player berhasil dihapus!",
		})
	}
}

func GetAllPlayer() gin.HandlerFunc {
	return func(c *gin.Context) {
		var player []model.Player

		findOptions := options.Find()
		findOptions.SetLimit(20)
		cur, err := playerCollection.Find(context.TODO(), bson.D{{}}, findOptions)
		if err != nil {
			log.Fatal(err)
		}
		for cur.Next(context.TODO()) {
			var elem model.Player
			err := cur.Decode(&elem)
			if err != nil {
				log.Fatal(err)
			}
			player = append(player, elem)
		}

		c.JSON(http.StatusOK, gin.H{
			"data":    &player,
			"error":   false,
			"message": http.StatusOK,
		})
		return
	}
}
