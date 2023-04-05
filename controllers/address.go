package controllers

import (
	"context"
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/A-Victory/e-commerce/models"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func AddAddress() gin.HandlerFunc {
	return func(c *gin.Context) {
		user_id := c.Query("id")
		if user_id == "" {
			log.Println("user_id is empty")

			c.AbortWithError(http.StatusNotFound, errors.New("invalid search index"))
			return
		}

		userID, err := primitive.ObjectIDFromHex(user_id)
		if err != nil {
			log.Println("userid is not valid")
			c.AbortWithStatusJSON(http.StatusInternalServerError, errors.New("user id is not valid"))
			return
		}

		var addresses models.Address
		addresses.Address_ID = primitive.NewObjectID()

		err = c.BindJSON(&addresses)
		if err != nil {
			log.Println(err)

			c.AbortWithStatusJSON(http.StatusInternalServerError, errors.New("Internal server error"))
			return
		}

		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()
		match_filter := bson.D{{Key: "$match", Value: bson.D{primitive.E{Key: "_id", Value: userID}}}}
		unwind := bson.D{{Key: "$unwind", Value: bson.D{primitive.E{Key: "path", Value: "$address"}}}}
		group := bson.D{{Key: "$group", Value: bson.D{primitive.E{Key: "path", Value: "$address_id"}, {Key: "count", Value: bson.D{primitive.E{Key: "$sum", Value: 1}}}}}}

		cursor, err := userCollection.Aggregate(ctx, mongo.Pipeline{match_filter, unwind, group})
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		}

		var addressinfo []bson.M
		if err = cursor.All(ctx, &addressinfo); err != nil {
			log.Panicln(err)
		}

		var size int
		for _, address_no := range addressinfo {
			count := address_no["count"]
			size = count.(int)
		}
		if size < 2 {
			filter := bson.D{primitive.E{Key: "_id", Value: userID}}
			update := bson.D{{Key: "$push", Value: primitive.E{Key: "address", Value: addresses}}}
			_, err = userCollection.UpdateOne(ctx, filter, update)
			if err != nil {
				c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
				return
			}
		} else {
			c.AbortWithStatusJSON(http.StatusBadRequest, "Not allowed")
			return
		}

	}
}

func EditHomeAddress() gin.HandlerFunc {
	return func(c *gin.Context) {
		user_id := c.Query("id")
		if user_id == "" {
			log.Println("user_id is empty")

			c.AbortWithError(http.StatusNotFound, errors.New("invalid search index"))
			return
		}

		userID, err := primitive.ObjectIDFromHex(user_id)
		if err != nil {
			log.Println("userid is not valid")
			c.AbortWithStatusJSON(http.StatusInternalServerError, errors.New("Internal server error"))
			return
		}

		var editaddress models.Address
		c.BindJSON(&editaddress)
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()
		filter := bson.D{primitive.E{Key: "_id", Value: userID}}
		update := bson.D{{Key: "$set", Value: bson.D{primitive.E{Key: "address.0.house_name", Value: editaddress.House}, {Key: "address.0.street_name", Value: editaddress.Street}, {Key: "address.0.city", Value: editaddress.City}, {Key: "address.0.zipcode", Value: editaddress.Zipcode}}}}
		_, err = userCollection.UpdateOne(ctx, filter, update)
		if err != nil {
			c.JSON(404, "something went wrong, try again")
			return
		}

		c.JSON(http.StatusOK, "Successfully update home address")
	}
}

func EditWorkAddress() gin.HandlerFunc {
	return func(c *gin.Context) {
		user_id := c.Query("id")
		if user_id == "" {
			log.Println("user_id is empty")

			c.AbortWithError(http.StatusNotFound, errors.New("invalid search index"))
			return
		}

		userID, err := primitive.ObjectIDFromHex(user_id)
		if err != nil {
			log.Println("userid is not valid")
			c.AbortWithStatusJSON(http.StatusInternalServerError, errors.New("Internal server error"))
			return
		}

		var editaddress models.Address
		c.BindJSON(&editaddress)
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()
		filter := bson.D{primitive.E{Key: "_id", Value: userID}}
		update := bson.D{{Key: "$set", Value: bson.D{primitive.E{Key: "address.0.house_name", Value: editaddress.House}, {Key: "address.0.street_name", Value: editaddress.Street}, {Key: "address.0.city", Value: editaddress.City}, {Key: "address.0.zipcode", Value: editaddress.Zipcode}}}}
		_, err = userCollection.UpdateOne(ctx, filter, update)
		if err != nil {
			c.JSON(404, "something went wrong, try again")
			return
		}

		c.JSON(http.StatusOK, "Successfully update home address")
	}
}

func DeleteAddress() gin.HandlerFunc {
	return func(c *gin.Context) {
		user_id := c.Query("id")
		if user_id == "" {
			log.Println("user_id is empty")

			c.AbortWithError(http.StatusNotFound, errors.New("invalid search index"))
			return
		}

		addresses := make([]models.Address, 0)
		userID, err := primitive.ObjectIDFromHex(user_id)
		if err != nil {
			log.Println("userid is not valid")
			c.AbortWithStatusJSON(http.StatusInternalServerError, errors.New("Internal server error"))
			return
		}

		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		filter := bson.D{primitive.E{Key: "_id", Value: userID}}
		update := bson.D{{Key: "$set", Value: primitive.E{Key: "address", Value: addresses}}}
		_, err = userCollection.UpdateOne(ctx, filter, update)
		if err != nil {
			c.JSON(404, "wrong command")
			return
		}
		defer cancel()

		c.JSON(http.StatusOK, "Successfully deleted address")
	}
}
