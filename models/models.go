package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID              primitive.ObjectID `json:"_id" bson:"_id"`
	First_Name      string             `json:"firstname" validate:"required,min=2,max=20"`
	Last_Name       string             `json:"lastname" validate:"required,min=2,max=20"`
	Password        string             `json:"password" validate:"required,alphanum,min=8"`
	Email           string             `json:"email" validate:"required,email"`
	Phone           string             `json:"phone"`
	User_ID         string             `json:"user_id"`
	UserCart        []ProductUser      `json:"usercart" bson:"usercart"`
	Token           string             `json:"token"`
	Refresh_Token   string             `json:"refresh_token"`
	Address_Details []Address          `json:"address" bson:"address"`
	Order_Status    []Order            `json:"orders" bson:"orders"`
	Created_At      time.Time          `json:"created_at"`
	Updated_At      time.Time          `json:"updated_at"`
}

type Product struct {
	Product_ID   primitive.ObjectID `bson:"_id"`
	Product_Name string             `json:"product_name"`
	Price        uint               `json:"price"`
	Rating       int                `json:"rating"`
}

type ProductUser struct {
	Product_ID   primitive.ObjectID `bson:"_id"`
	Product_Name string             `json:"product_name" bson:"product_name"`
	Price        uint               `json:"price" bson:"price"`
	Rating       int                `json:"rating" bson:"rating"`
}

type Address struct {
	Address_ID primitive.ObjectID `bson:"_id"`
	House      string             `json:"house_name" bson:"house_name"`
	Street     string             `json:"street_name" bson:"street_name"`
	City       string             `json:"city_name" bson:"city_name"`
	Zipcode    string             `json:"zipcode" bson:"zipcode"`
}

type Order struct {
	Order_ID       primitive.ObjectID `bson:"_id"`
	Order_Cart     []ProductUser      `json:"order_list" bson:"order_list"`
	Ordered_At     time.Time          `json:"ordered_at" bson:"ordered_at"`
	Price          int                `json:"total_price" bson:"total_price"`
	Discount       int                `json:"discount" bson:"discount"`
	Payment_Method Payment            `json:"payment_method" bson:"payment_method"`
}

type Payment struct {
	Digital bool
	COD     bool
}
