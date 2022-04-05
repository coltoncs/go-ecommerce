package controllers

import (
	"context"
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/pizdetz/go-ecommerce/database"
	"github.com/pizdetz/go-ecommerce/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type Application struct {
	prodCollection *mongo.Collection
	userCollection *mongo.Collection
}

func NewApplication(prodCollection, userCollection *mongo.Collection) *Application {
	return &Application{
		prodCollection: prodCollection,
		userCollection: userCollection,
	}
}

/**
* #Add To Cart Function
*
* @returns gin.Handler func for passing context to controller
*
* 1) Retrieve queried product ID from context
* 2) If product ID is empty, return bad request with error and return from func
* 3) Retrieve user ID from context
* 4) If user ID is empty, return bad request with error and return from func
* 5) Create a new ObjectID from the user ID hex string
* 6) If error returns from ObjectID conversion, return from func
* 7) Prepare for database call by passing 5s Timeout to Context
* 8) Call AddProductToCart with the current product and user collections, product to add ID, 
*    and the user ID of the cart to add to
*
**/
func (app *Application) AddToCart() gin.Handler {
	return func(c *gin.Context) {
		productQueryID := c.Query("id") //1
		if productQueryID == "" {
			log.Println("Product ID is empty")
			_ = c.AbortWithError(http.StatusBadRequest, errors.New("Product ID is empty.")) //2
			return
		}

		userQueryID := c.Query("userID") //3
		if userQueryID == "" {
			log.Println("User ID is empty")
			_ = c.AbortWithError(http.StatusBadRequest, errors.New("User ID is empty.")) //4
			return
		}

		productID, err := primitive.ObjectIDFromHex(productQueryID) //5

		if err != nil {
			log.Println(err)
			c.AbortWithError(http.StatusInternalServerError) //6
			return
		}

		var ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second) //7

		defer cancel()

		err = database.AddProductToCart(ctx, app.prodCollection, app.userCollection, productID, userQueryID) //8
		if err != nil {
			c.IndentedJSON(http.StatusInternalServerError, err)
		}
		c.IndentedJSON(200, "Successfully added to the cart")
	}
}

/**
* #Remove Item Function
*
* @returns gin.Handler func for passing context to controller
*
* 1) Retrieve queried product ID from context
* 2) If product ID is empty, return bad request with error and return from func
* 3) Retrieve user ID from context
* 4) If user ID is empty, return bad request with error and return from func
* 5) Create a new ObjectID from the user ID hex string
* 6) If error returns from ObjectID conversion, return from func
* 7) Prepare for database call by passing 5s Timeout to Context
*
**/
func (app *Application) RemoveItem() gin.HandlerFunc {
	return func(c *gin.Context){
		productQueryID := c.Query("id") //1
		if productQueryID == "" {
			log.Println("Product ID is empty")
			_ = c.AbortWithError(http.StatusBadRequest, errors.New("Product ID is empty.")) //2
			return
		}

		userQueryID := c.Query("userID") //3
		if userQueryID == "" {
			log.Println("User ID is empty")
			_ = c.AbortWithError(http.StatusBadRequest, errors.New("User ID is empty.")) //4
			return
		}

		productID, err := primitive.ObjectIDFromHex(productQueryID) //5

		if err != nil {
			log.Println(err)
			c.AbortWithError(http.StatusInternalServerError) //6
			return
		}

		var ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second) //7

		defer cancel()

		err = database.RemoveCartItem(ctx, app.prodCollection, app.userCollection, productID, userQueryID)
		if err != nil {
			c.IndentedJSON(http.StatusInternalServerError, err)
			return
		}
		c.IndentedJSON(200, "Successfully removed item from cart.")
	}
}

/**
* #Get Item From Cart Function
*
*
**/
func GetItemFromCart() gin.HandlerFunc {
	return func(c *gin.Context){
		user_id := c.Query("id")

		if user_id == ""{
			c.Header("Content-Type", "application/json")
			c.JSON(http.StatusNotFound, gin.H{"error":"invalid id"})
			c.Abort()
			return
		}

		usern_id, _ := primitive.ObjectIDFromHex(user_id)

		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		var filledCart models.User
		err := UserCollection.FindOne(ctx, bson.D{primitive.E{Key:"_id", Value: usern_id}}).Decode(&filledCart)

		if err != nil {
			log.Println(err)
			c.IndentedJSON(500, "not found")
			return
		}

		filter_match := bson.D{{Key: "$match", Value: bson.D{primitive.E{Key: "_id", Value: usern_id}}}}
		unwind := bson.D{{Key: "$unwind", Value: bson.D{primitive.E{Key: "path", Value: "$usercart"}}}}
		grouping := bson.D{{Key: "$group", Value: bson.D{primitive.E{Key: "_id", Value: "$_id"}, {Key: "total", Value: bson.D{primitive.E{Key: "$sum", Value: "$usercart.price"}}}}}}
		pointCursor, err := UserCollection.Aggregate(ctx, mongo.Pipeline{filter_match, unwind, grouping})
		if err != nil {
			log.Println(err)
		}
		var listing []bson.M
		if err = pointCursor.All(ctx, &listing); err != nil {
			log.Println(err)
			c.AbortWithStatus(http.StatusInternalServerError)
		}
		for _, json := range listing {
			c.IndentedJSON(200, json["total"])
			c.IndentedJSON(200, filledCart.UserCart)
		}
		ctx.Done()
	}
}

/**
* #Buy From Cart Function
*
*
**/
func (app *Application) BuyFromCart() gin.HandlerFunc {
	return func(c *gin.Context){
		userQueryID := c.Query("id")
		if userQueryID == ""{
			log.Panicln("user id is empty!")
			_ = c.AbortWithError(http.StatusBadRequest, errors.New("UserID is empty"))
		}

		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)

		defer cancel()

		err := database.BuyItemFromCart(ctx, app.userCollection, userQueryID)
		if err != nil {
			c.IndentedJSON(http.StatusInternalServerError, err)
		}

		c.IndentedJSON(200, "successfully placed the order")
	}
}

/**
* #Instant Buy Function
*
* @returns gin.Handler func for passing context to controller
*
* 1) Retrieve queried product ID from context
* 2) If product ID is empty, return bad request with error and return from func
* 3) Retrieve user ID from context
* 4) If user ID is empty, return bad request with error and return from func
* 5) Create a new ObjectID from the user ID hex string
* 6) If error returns from ObjectID conversion, return from func
* 7) Prepare for database call by passing 5s Timeout to Context
*
**/
func (app *Application) InstantBuy() gin.HandlerFunc {
	return func(c *gin.Context){
		productQueryID := c.Query("id") //1
		if productQueryID == "" {
			log.Println("Product ID is empty")
			_ = c.AbortWithError(http.StatusBadRequest, errors.New("Product ID is empty.")) //2
			return
		}

		userQueryID := c.Query("userID") //3
		if userQueryID == "" {
			log.Println("User ID is empty")
			_ = c.AbortWithError(http.StatusBadRequest, errors.New("User ID is empty.")) //4
			return
		}

		productID, err := primitive.ObjectIDFromHex(productQueryID) //5

		if err != nil {
			log.Println(err)
			c.AbortWithError(http.StatusInternalServerError) //6
			return
		}

		var ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second) //7

		defer cancel()

		err = database.InstantBuy(ctx, app.prodCollection, app.userCollection, productID, userQueryID)
		if err != nil {
			c.IndentedJSON(http.StatusInternalServerError, err)
		}

		c.IndentedJSON(200, "Successfully bought item")
	}
}