package database

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/A-Victory/e-commerce/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var (
	ErrCantFindProduct    = errors.New("product could not be found")
	ErrInvalidUserID      = errors.New("userID is invalid")
	ErrUserUpdateFailed   = errors.New("failed to update user")
	ErrCantDecodeProduct  = errors.New("product could not be decoded")
	ErrCantRemoveItemCart = errors.New("failed to remove item from cart")
	ErrUnableToGetItem    = errors.New("failed to retrieve product")
	ErrCartItemBuyFailed  = errors.New("unable to buy product from cart")
)

// AddProductToCart adds a product to products database
func AddProductToCart(ctx context.Context, prodCollection *mongo.Collection, userCollection *mongo.Collection, prodID primitive.ObjectID, userID string) error {
	searchfromdb, err := prodCollection.Find(ctx, bson.M{"_id": prodID})
	if err != nil {
		log.Println(err)
		return ErrCantFindProduct
	}

	var productCart []models.ProductUser

	err = searchfromdb.All(ctx, &productCart)
	if err != nil {
		log.Println(err)
		return ErrCantDecodeProduct
	}

	id, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		log.Println(err)
		return ErrInvalidUserID
	}

	filter := bson.D{primitive.E{Key: "_id", Value: id}}
	update := bson.D{{Key: "$push", Value: bson.D{primitive.E{Key: "usercart", Value: bson.D{{Key: "$each", Value: productCart}}}}}}

	_, err = userCollection.UpdateOne(ctx, filter, update)
	if err != nil {
		log.Println(err)
		return ErrUserUpdateFailed
	}

	return nil
}

// RemoveCartItem removes a product from the user's cart
func RemoveCartItem(ctx context.Context, prodCollection *mongo.Collection, userCollection *mongo.Collection, prodID primitive.ObjectID, userID string) error {
	id, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		log.Println(err)
		return ErrInvalidUserID
	}

	filter := bson.D{primitive.E{Key: "_id", Value: id}}
	update := bson.M{"$pull": bson.M{"usercart": bson.M{"_id": prodID}}}
	_, err = userCollection.UpdateMany(ctx, filter, update)
	if err != nil {
		log.Println(err)
		return ErrCantRemoveItemCart
	}

	return nil
}

// BuyItemFromCart is the helper function for database interaction for the BuyFromCart function in controllers.Cart
func BuyItemFromCart(ctx context.Context, userCollection *mongo.Collection, userID string) error {
	id, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		log.Println(err)
		return ErrInvalidUserID
	}

	var getcartitems models.User
	var ordercart models.Order

	ordercart.Ordered_At = time.Now()
	ordercart.Order_ID = primitive.NewObjectID()
	ordercart.Payment_Method.COD = true
	ordercart.Order_Cart = make([]models.ProductUser, 0)

	unwind := bson.D{{Key: "$unwind", Value: primitive.E{Key: "path", Value: "$usercart"}}}
	grouping := bson.D{{Key: "$group", Value: bson.D{primitive.E{Key: "_id", Value: "$_id"}, {Key: "total", Value: bson.D{primitive.E{Key: "$sum", Value: "$usercart.price"}}}}}}

	currentresults, err := userCollection.Aggregate(ctx, mongo.Pipeline{unwind, grouping})
	if err != nil {
		log.Panicln(err)
	}

	var getusercart []bson.M
	err = currentresults.All(ctx, &getusercart)
	if err != nil {
		log.Panicln(err)
	}

	var total_price int

	for _, user_item := range getusercart {
		price := user_item["total"]
		total_price = price.(int)
	}

	ordercart.Price = total_price

	filter := bson.D{primitive.E{Key: "_id", Value: id}}
	update := bson.D{{Key: "$push", Value: bson.D{primitive.E{Key: "orders", Value: ordercart}}}}

	_, err = userCollection.UpdateMany(ctx, filter, update)
	if err != nil {
		log.Println(err)
	}

	if err = userCollection.FindOne(ctx, bson.D{primitive.E{Key: "_id", Value: id}}).Decode(&getcartitems); err != nil {
		log.Println(err)
	}

	filter2 := bson.D{primitive.E{Key: "_id", Value: id}}
	update2 := bson.M{"$push": bson.M{"orders.$[].order_list": bson.M{"$each": getcartitems.UserCart}}}
	_, err = userCollection.UpdateOne(ctx, filter2, update2)
	if err != nil {
		log.Println(err)
	}

	empty_usercart := make([]models.ProductUser, 0)
	filter3 := bson.D{primitive.E{Key: "_id", Value: id}}
	update3 := bson.D{{Key: "$set", Value: bson.D{primitive.E{Key: "usercart", Value: empty_usercart}}}}
	_, err = userCollection.UpdateOne(ctx, filter3, update3)
	if err != nil {
		log.Println(err)
		return ErrCartItemBuyFailed
	}

	return nil
}

// InstantBuy is the helper function for database interaction for the InstantBuy function in controllers.Cart
func InstantBuy(ctx context.Context, prodCollection *mongo.Collection, userCollection *mongo.Collection, prodID primitive.ObjectID, userID string) error {
	id, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		log.Println(err)
		return ErrInvalidUserID
	}

	var product_details models.ProductUser
	var orders_details models.Order

	orders_details.Ordered_At = time.Now()
	orders_details.Order_ID = primitive.NewObjectID()
	orders_details.Payment_Method.COD = true
	orders_details.Order_Cart = make([]models.ProductUser, 0)

	err = prodCollection.FindOne(ctx, bson.D{primitive.E{Key: "_id", Value: prodID}}).Decode(&product_details)
	if err != nil {
		log.Println(err)
		return ErrCantFindProduct
	}

	orders_details.Price = int(product_details.Price)

	filter := bson.D{primitive.E{Key: "_id", Value: id}}
	update := bson.D{{Key: "$push", Value: bson.D{primitive.E{Key: "orders", Value: orders_details}}}}

	_, err = userCollection.UpdateOne(ctx, filter, update)
	if err != nil {
		log.Println(err)
	}

	filter2 := bson.D{primitive.E{Key: "_id", Value: id}}
	update2 := bson.M{"$push": bson.M{"orders.$[].order_list": product_details}}

	_, err = userCollection.UpdateOne(ctx, filter2, update2)
	if err != nil {
		log.Println(err)
		return errors.New("buy failed")
	}

	return nil
}
