package controllers

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/pizdetz/go-ecommerce/models"
)

/**
* #Hash Password Function
*
*
*
**/
func HashPassword(password string) string{

}

/**
* #Verify Password Function
*
*
*
**/
func VerifyPassword(userPassword string, givenPassword string) (bool, string) {

}

/**
* #Sign Up Function
*
* @returns gin.HandlerFunc context object for parsing request.
*
* 1) Receives context from Gin HandlerFunc; if error, timeout.
* 2) Attempt to bind JSON data to `user` object, return Bad Request if error returns.
* 3) If `user` cannot be validated, return another Bad Request.
* 4) Receive # of `user`'s in UserCollection mongodb collection as `count`, if `count > 0` return error as user already exists.
* 5) Checks if phone number is already in UserCollection, returns error if true.
* 6) Hash `user.Password` and set hashed password back to the `user.Password` variable.
* 7) Get current time and set that value to both `user.Create_At` and `user.Updated_At`
* 8) Set `user.ID` and `user.User_ID` to new object ID and a hex cast of the ID respectively.
* 9) Set `token` and `refreshToken` from tokengen.go in /tokens
*
**/
func SignUp() gin.HandlerFunc {
	return func(c *gin.Context){
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		var user models.User
		if err := c.BindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		validationErr := Validate.Struct(user)
		if validationErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.validationErr()})
			return
		}

		count, err := UserCollection.CountDocuments(ctx, bson.M{})
		if err != nil {
			log.Panic(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err})
			return
		}

		if count > 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "user already exists"})
		}

		count, err = UserCollection.CountDocuments(ctx, bson.M{"phone": user.Phone})

		defer cancel()
		if err != nil {
			log.Panic(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err})
			return
		}

		if count>0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "this phone number already exists"})
			return
		}

		password := HashPassword(*user.Password)
		user.Password = &password

		user.Created_At, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		user.Updated_At, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		user.ID = primitive.NewObjectID()
		user.User_ID = user.ID.Hex()

		token, refreshToken, _ = generate.TokenGenerator(*user.Email, *user.First_Name, *user.Last_Name, user.User_ID)
		user.Token = &token
		user.Refresh_Token = &refreshToken
		user.UserCart = make([]models.ProductUser, 0)
		user.Address_Details = make([]models.Address, 0)
		user.Order_Status = make([]models.Order, 0)
		_, insertErr := UserCollection.InsertOne(ctx, user)
		if insertErr != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "the user did not get created"})
			return 
		}
		defer cancel()

		c.JSON(http.StatusCreated, "Successfully signed in!")
	}
}


/**
* #Login Function
*
*
*
**/
func Login() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		var user models.User
		if err := c.BindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err})
			return
		}

		err := UserCollection.FindOne(ctx, bson.M{"email": user.Email}).Decode(&foundUser)
		defer cancel()

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "login or password incorrect"})
			return
		}

		PasswordIsValid, msg := VerifyPassword(*user.Password, *foundUser.Password)

		defer cancel()

		if !PasswordIsValid {
			c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
			fmt.Println(msg)
			return
		}

		token, refreshToken, _ := generate.TokenGenerator(*foundUser.Email, *foundUser.First_Name, *foundUser.Last_Name)
		defer cancel()

		generate.UpdateAllTokens(token, refreshToken, foundUser.User_ID)

		c.JSON(http.StatusFound, foundUser)
	}
}

/**
* #ProductViewerAdmin Function
*
*
*
**/
func ProductViewerAdmin() gin.HandlerFunc {

}

/**
* #Search Product Function
*
*
*
**/
func SearchProduct() gin.HandlerFunc {

}

/**
* #Search Product By Query Function
*
*
*
**/
func SearchProductByQuery() gin.HandlerFunc {

}