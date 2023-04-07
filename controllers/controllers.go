package controllers

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/A-Victory/e-commerce/database"
	"github.com/A-Victory/e-commerce/models"
	"github.com/A-Victory/e-commerce/tokens"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var userCollection *mongo.Collection = database.UserData(database.Client, "user")
var prodCollection *mongo.Collection = database.ProductData(database.Client, "product")

func Signup() gin.HandlerFunc {

	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		var user models.User
		if err := c.BindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		validationErr := validate().Struct(user)
		if validationErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": validationErr})
			return
		}

		pass, err := hashPasswd(user.Password)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to encrypt password"})
			return
		}

		user.Password = pass
		user.Created_At, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		user.Updated_At, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		user.ID = primitive.NewObjectID()
		user.User_ID = user.ID.Hex()
		user.Order_Status = make([]models.Order, 0)
		user.Address_Details = make([]models.Address, 0)
		user.UserCart = make([]models.ProductUser, 0)
		token, refreshToken, err := tokens.GenerateToken(user.First_Name, user.Last_Name, user.Email, user.User_ID)
		user.Token = token
		user.Refresh_Token = refreshToken

		_, inserterr := userCollection.InsertOne(ctx, user)
		if inserterr != nil {
			log.Println(err)
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": inserterr})
			return
		}

		c.JSON(http.StatusOK, "Successfully signed in")
	}
}

func Login() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()
		var user models.User

		if err := c.BindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		var founduser models.User
		err := userCollection.FindOne(ctx, bson.M{"email": user.Email}).Decode(&founduser)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Email or password incorrect!"})
			return
		}

		valid, msg := verifyPasswd(founduser.Password, user.Password)
		if !valid {
			c.JSON(http.StatusBadRequest, gin.H{"error": msg})
			return
		}

		token, refreshToken, err := tokens.GenerateToken(founduser.First_Name, founduser.Last_Name, founduser.Email, founduser.User_ID)
		if err != nil {
			log.Println(err)
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err})
			return
		}

		tokens.UpdateToken(token, refreshToken, founduser.User_ID)

		c.JSON(http.StatusOK, founduser)

	}
}

func ProductView() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()
		var products models.Product
		if err := c.BindJSON(&products); err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		products.Product_ID = primitive.NewObjectID()
		_, err := prodCollection.InsertOne(ctx, products)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, "successfully addded product")
	}
}

func SearchProduct() gin.HandlerFunc {
	return func(c *gin.Context) {

		var productlist []models.Product

		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		cursor, err := prodCollection.Find(ctx, bson.D{})
		if err != nil {
			c.JSON(http.StatusInternalServerError, "Something went wrong please try again")
			return
		}
		defer cancel()

		err = cursor.All(ctx, &productlist)
		if err != nil {
			log.Println(err)
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}
		defer cancel()

		if err := cursor.Err(); err != nil {
			log.Println(err)
			c.JSON(300, "invalid")
			return
		}
		c.JSON(http.StatusOK, productlist)
	}
}

func SearchProductByQuery() gin.HandlerFunc {
	return func(c *gin.Context) {
		var searchProduct []models.Product
		queryParam := c.Query("name")

		if queryParam == "" {
			log.Println("query is empty")
			c.JSON(http.StatusNotFound, "query is empty")
			return
		}

		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		searchQuery, err := prodCollection.Find(ctx, bson.M{"product_name": bson.M{"$regex": queryParam}})
		if err != nil {
			c.JSON(http.StatusNotFound, "something went wrong while fetching data")
			return
		}

		err = searchQuery.All(ctx, &searchProduct)
		if err != nil {
			log.Println(err)
			c.JSON(http.StatusNotFound, "invalid")
			return
		}
		defer searchQuery.Close(ctx)

		if err := searchQuery.Err(); err != nil {
			log.Println(err)
			c.JSON(http.StatusNotFound, "invalid request")
			return
		}

		defer cancel()

		c.JSON(http.StatusOK, searchProduct)
	}
}
