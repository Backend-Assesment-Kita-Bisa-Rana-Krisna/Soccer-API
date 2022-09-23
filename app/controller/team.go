package controller

import (
	"context"
	"net/http"
	"soccer-api/app/model"
	"soccer-api/configuration"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var teamCollection *mongo.Collection = configuration.Collection("team")
var validate = validator.New()

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

		if validationErr := validate.Struct(&team); validationErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"data":    bson.M{},
				"error":   true,
				"message": validationErr,
			})
			return
		}

		data := model.Team{
			ID:   primitive.NewObjectID(),
			Name: team.Name,
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
				"message": err,
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
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
			c.JSON(http.StatusBadRequest, helpers.Get(http.StatusBadRequest, err.Error(), false, nil))
			return
		}

		if validationErr := validate.Struct(&team); validationErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"data":    bson.M{},
				"error":   true,
				"message": validationErr,
			})
			return
		}

		update := bson.M{"name": team.Name}
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
				"message": err,
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
			c.JSON(http.StatusInternalServerError, helpers.Get(http.StatusInternalServerError, err.Error(), false, nil))
			return
		}

		if result.DeletedCount < 1 {
			c.JSON(http.StatusNotFound, helpers.Get(http.StatusNotFound, "Team with specified ID not found!", false, nil))
			return
		}

		c.JSON(http.StatusOK, helpers.Get(http.StatusOK, "Team berhasil dihapus!", true, nil))
	}
}

func GetAllUsers() gin.HandlerFunc {
	type GetAllUsers struct {
		ID        primitive.ObjectID `json:"id" bson:"_id"`
		Name      string             `json:"name"`
		Username  string             `json:"username"`
		Level     string             `json:"level"`
		Provinsi  string             `json:"provinsi"`
		Kota      string             `json:"kota"`
		Kecamatan string             `json:"kecamatan"`
		Kelurahan string             `json:"kelurahan"`
		Wilayah   string             `json:"wilayah" bson:"wilayah"`
		ApiToken  string             `json:"api_token" bson:"api_token"`
		Status    bool               `json:"status"`
		CreatedAt time.Time          `json:"created_at" bson:"created_at"`
		UpdatedAt time.Time          `json:"updated_at" bson:"updated_at"`
	}
	return func(c *gin.Context) {
		provinsi_id, _ := strconv.Atoi(c.Query("provinsi_id"))
		kota_id, _ := strconv.Atoi(c.Query("kota_id"))
		kecamatan_id, _ := strconv.Atoi(c.Query("kecamatan_id"))
		kelurahan_id, _ := strconv.Atoi(c.Query("kelurahan_id"))
		level := c.Query("level")
		status := c.Query("status")

		filter := []bson.M{}
		statusBool, _ := strconv.Atoi(status)
		if status != "" && (statusBool == 0 || statusBool == 1) {
			statusData := false
			if statusBool == 1 {
				statusData = true
			}
			filter = append(filter, bson.M{"status": bson.M{"$eq": statusData}})
		} else {
			filter = append(filter, bson.M{"status": bson.M{"$in": []bool{true, false}}})
		}
		if level != "" {
			filter = append(filter, bson.M{"level": bson.M{"$eq": level}})
		}
		if provinsi_id > 0 {
			filter = append(filter, bson.M{"provinsi_id": bson.M{"$eq": provinsi_id}})
		}
		if kota_id > 0 {
			filter = append(filter, bson.M{"kota_id": bson.M{"$eq": kota_id}})
		}
		if kecamatan_id > 0 {
			filter = append(filter, bson.M{"kecamatan_id": bson.M{"$eq": kecamatan_id}})
		}
		if kelurahan_id > 0 {
			filter = append(filter, bson.M{"kelurahan_id": bson.M{"$eq": kelurahan_id}})
		}
		pipe := []bson.M{
			{
				"$project": bson.M{
					"name":         1,
					"username":     1,
					"level":        1,
					"wilayah":      1,
					"api_token":    1,
					"status":       1,
					"created_at":   1,
					"updated_at":   1,
					"provinsi_id":  bson.M{"$cond": []interface{}{bson.M{"$or": []bson.M{{"$ne": []string{"$wilayah", ""}}, {"$in": []interface{}{"$level", []string{"provinsi", "kota", "kecamatan", "kelurahan"}}}}}, bson.M{"$convert": bson.M{"input": bson.M{"$substr": []interface{}{"$wilayah", 0, 2}}, "to": "long", "onError": 0, "onNull": 0}}, 0}},
					"kota_id":      bson.M{"$cond": []interface{}{bson.M{"$or": []bson.M{{"$ne": []string{"$wilayah", ""}}, {"$in": []interface{}{"$level", []string{"kota", "kecamatan", "kelurahan"}}}}}, bson.M{"$convert": bson.M{"input": bson.M{"$substr": []interface{}{"$wilayah", 0, 4}}, "to": "long", "onError": 0, "onNull": 0}}, 0}},
					"kecamatan_id": bson.M{"$cond": []interface{}{bson.M{"$or": []bson.M{{"$ne": []string{"$wilayah", ""}}, {"$in": []interface{}{"$level", []string{"kecamatan", "kelurahan"}}}}}, bson.M{"$convert": bson.M{"input": bson.M{"$substr": []interface{}{"$wilayah", 0, 7}}, "to": "long", "onError": 0, "onNull": 0}}, 0}},
					"kelurahan_id": bson.M{"$cond": []interface{}{bson.M{"$or": []bson.M{{"$ne": []string{"$wilayah", ""}}, {"$in": []interface{}{"$level", []string{"kelurahan"}}}}}, bson.M{"$convert": bson.M{"input": bson.M{"$substr": []interface{}{"$wilayah", 0, 10}}, "to": "long", "onError": 0, "onNull": 0}}, 0}},
				},
			},
			{
				"$match": bson.M{
					"$and": filter,
				},
			},
			{
				"$lookup": bson.M{
					"from": "m_provinsi",
					"let": bson.M{
						"wilayah_id": bson.M{
							"$toInt": "$provinsi_id",
						},
					},
					"pipeline": []bson.M{
						{
							"$match": bson.M{
								"$expr": bson.M{
									"$eq": []string{
										"$id",
										"$$wilayah_id",
									},
								},
							},
						},
						{
							"$project": bson.M{
								"_id":  0,
								"id":   1,
								"nama": 1,
							},
						},
					},
					"as": "provinsi",
				},
			},
			{
				"$lookup": bson.M{
					"from": "m_kota",
					"let": bson.M{
						"wilayah_id": bson.M{
							"$toInt": "$kota_id",
						},
					},
					"pipeline": []bson.M{
						{
							"$match": bson.M{
								"$expr": bson.M{
									"$eq": []string{
										"$id",
										"$$wilayah_id",
									},
								},
							},
						},
						{
							"$project": bson.M{
								"_id":  0,
								"id":   1,
								"nama": 1,
							},
						},
					},
					"as": "kota",
				},
			},
			{
				"$lookup": bson.M{
					"from": "m_kecamatan",
					"let": bson.M{
						"wilayah_id": bson.M{
							"$toInt": "$kecamatan_id",
						},
					},
					"pipeline": []bson.M{
						{
							"$match": bson.M{
								"$expr": bson.M{
									"$eq": []string{
										"$id",
										"$$wilayah_id",
									},
								},
							},
						},
						{
							"$project": bson.M{
								"_id":  0,
								"id":   1,
								"nama": 1,
							},
						},
					},
					"as": "kecamatan",
				},
			},
			{
				"$lookup": bson.M{
					"from": "m_kelurahan",
					"let": bson.M{
						"wilayah_id": bson.M{
							"$toLong": "$kelurahan_id",
						},
					},
					"pipeline": []bson.M{
						{
							"$match": bson.M{
								"$expr": bson.M{
									"$eq": []string{
										"$id",
										"$$wilayah_id",
									},
								},
							},
						},
						{
							"$project": bson.M{
								"_id":  0,
								"id":   1,
								"nama": 1,
							},
						},
					},
					"as": "kelurahan",
				},
			},
			{"$addFields": bson.M{"provinsi": bson.M{"$ifNull": []interface{}{bson.M{"$arrayElemAt": []interface{}{"$provinsi.nama", 0}}, ""}}}},
			{"$addFields": bson.M{"kota": bson.M{"$ifNull": []interface{}{bson.M{"$arrayElemAt": []interface{}{"$kota.nama", 0}}, ""}}}},
			{"$addFields": bson.M{"kecamatan": bson.M{"$ifNull": []interface{}{bson.M{"$arrayElemAt": []interface{}{"$kecamatan.nama", 0}}, ""}}}},
			{"$addFields": bson.M{"kelurahan": bson.M{"$ifNull": []interface{}{bson.M{"$arrayElemAt": []interface{}{"$kelurahan.nama", 0}}, ""}}}},
		}
		cursor, err := teamCollection.Aggregate(context.TODO(), pipe)
		if err != nil {
			c.JSON(http.StatusInternalServerError, helpers.Get(http.StatusInternalServerError, err.Error(), false, nil))
		}
		var results []GetAllUsers
		if err = cursor.All(context.TODO(), &results); err != nil {
			c.JSON(http.StatusInternalServerError, helpers.Get(http.StatusInternalServerError, err.Error(), false, nil))
		}
		if err := cursor.Close(context.TODO()); err != nil {
			c.JSON(http.StatusInternalServerError, helpers.Get(http.StatusInternalServerError, err.Error(), false, nil))
		}

		c.JSON(http.StatusOK, helpers.Get(http.StatusOK, http.StatusText(http.StatusOK), true, results))
	}
}
