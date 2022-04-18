package controllers

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/pizdetz/go-ecommerce/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

/**
* #Add Address Function
*
*
**/
func AddAddress() gin.HandlerFunc {
	return func(c *gin.Context) {
		user_id := c.Query("id")
		if user_id == "" {
			c.Header("Content-Type", "application/json")
			c.JSON(http.StatusNotFound, gin.H{"error":"Invalid code"})
			c.Abort()
			return
		}

		address, err := ObjectIDFromHex(user_id)
		if err != nil {
			c.IndentedJSON(500, "Internal Server Error")
		}

		var addresses models.Address

		addresses.Address_ID = primitive.NewObjectID()

		if err = c.BindJSON(&address); err != nil {
			c.IndentedJSON(http.StatusNotAcceptable, err.Error())
		}

		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)

		match_filter := bson.D{
			Key: "$match", 
			Value: bson.D{
				primitive.E{
					Key: "_id", 
					Value: address
				}
			}
		}
		unwind := bson.D{
			Key: "$unwind", 
			Value: bson.D{
				primitive.E{
					Key: "path", 
					Value: "$address",
				},
			},
		}
		group := bson.D{
			Key: "$group", 
			Value: bson.D{
				primitive.E{
					Key: "_id", 
					Value: "$address_id",
					}, 
				Key: "count", 
				Value: bson.D{
					primitive.E{
						Key: "$sum", 
						Value: 1,
					},
				},
			},
		}
		pointcursor, err := UserCollection.Aggregate(ctx, mongo.Pipeline{match_filter, unwind, group})
		if err != nil {
			c.IndentedJSON(500, "Internal Server Error")
		}

		var addressInfo []bson.M
		if err = pointcursor.All(ctx, &addressInfo); err != nil {
			panic(err)
		} 

		var size int32
		for _, address_no := range addressInfo {
			count := address_no["count"]
			size := count.(int32)
		}
		if size < 2 {
			filter := bson.D{primitive.E{Key: "_id", Value: address}}
			update := bson.D{Key: "$push", Value: bson.D{primitive.E{Key: "address", Value: addresses}}}
			_, err := UserCollection.UpdateOne(ctx, filter, update)
			if err != nil {
				fmt.Println(err)
			}
		} else {
			c.IndentedJSON(400, "Not Allowed")
		}
		defer cancel()
		ctx.Done()
	}
}

/**
* #Edit Home Address Function
*
*
**/
func EditHomeAddress() gin.HandlerFunc {
	
}

/**
* #Edit Work Address Function
*
*
**/
func EditWorkAddress() gin.HandlerFunc {
	
}

/**
* #Get Item From Cart Function
*
*
**/
func DeleteAddress() gin.HandlerFunc {
	return func(c *gin.Context) {
		user_id := c.Query("id")

		if user_id == "" {
			c.Header("Content-Type", "application/json")
			c.IndentedJSON(http.StatusNotFound, gin.H{"Error": "Invalid Search Index"})
			c.Abort()
			return
		}

		addresses := make([]models.Address, 0)
		usern_id, err := primitive.ObjectIDFromHex(user_id)
		if err != nil {
			c.IndentedJSON(500, "Internal Server Error")
		}

		var ctx, cancel = ctx.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()
		filter := bson.D{primitive.E{Key: "_id", Value: usern_id}}
		update := bson.D{{Key: "$set", Value: bson.D{primitive.E{Key: "address", Value: addresses}}}}
		_, err = UserCollection.UpdateOne(ctx, filter, update)
		if err != nil {
			c.IndentedJSON(404, "Wrong command")
			return
		}
		defer cancel()
		ctx.Done()
		c.IndentedJSON(200, "Successfully deleted")
	}
}