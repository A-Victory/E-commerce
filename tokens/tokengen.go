package tokens

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/A-Victory/e-commerce/database"
	"github.com/dgrijalva/jwt-go"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	secretkey = os.Getenv("SECRETKEY")
	userData  = database.UserData(database.Client, "user")
)

type SignedDetails struct {
	Email      string
	First_Name string
	Last_Name  string
	Uid        string
	jwt.StandardClaims
}

func GenerateToken(firstname, lastname, email, uid string) (signedtoken string, signedrefreshtoken string, err error) {

	claims := &SignedDetails{
		Email:      email,
		Last_Name:  lastname,
		First_Name: lastname,
		Uid:        uid,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Local().Add(time.Hour * time.Duration(24)).Unix(),
		},
	}

	refreshClaims := &SignedDetails{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Local().Add(time.Hour * time.Duration(72)).Unix(),
		},
	}

	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(secretkey))
	if err != nil {
		return "", "", err
	}

	refreshtoken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims).SignedString([]byte(secretkey))
	if err != nil {
		log.Panic(err)
		return
	}

	return token, refreshtoken, nil

}

func ValidateToken(signedToken string) (claims *SignedDetails, msg string) {
	token, err := jwt.ParseWithClaims(signedToken, &SignedDetails{}, func(t *jwt.Token) (interface{}, error) {
		return []byte(secretkey), nil
	})
	if err != nil {
		msg = err.Error()
		return
	}

	claims, ok := token.Claims.(*SignedDetails)
	if !ok {
		msg = "Invalid token"
		return
	}

	if claims.ExpiresAt < time.Now().Local().Unix() {
		msg = "Token expired"
		return
	}

	return claims, msg

}

func UpdateToken(signedtoken, signedrefreshtoken, userID string) {

	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)

	var updateObj primitive.D

	updateObj = append(updateObj, bson.E{Key: "token", Value: signedtoken})
	updateObj = append(updateObj, bson.E{Key: "refresh_token", Value: signedrefreshtoken})
	updated_at, _ := time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	updateObj = append(updateObj, bson.E{Key: "updated_at", Value: updated_at})

	upsert := true

	opt := options.UpdateOptions{
		Upsert: &upsert,
	}

	filter := bson.M{"user_id": userID}

	userData.UpdateOne(ctx, filter, bson.D{
		{Key: "$set", Value: updateObj},
	}, &opt)

	defer cancel()
}
