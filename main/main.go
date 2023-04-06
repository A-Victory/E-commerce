package main

import (
	"log"
	"os"

	"github.com/A-Victory/e-commerce/controllers"
	"github.com/A-Victory/e-commerce/database"
	"github.com/A-Victory/e-commerce/middleware"
	"github.com/A-Victory/e-commerce/routes"
	"github.com/gin-gonic/gin"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	app := controllers.NewApplication(database.ProductData(database.Client, "product"), database.UserData(database.Client, "user"))

	router := gin.New()
	router.Use(gin.Logger())

	routes.UserRoutes(router)
	router.Use(middleware.Authentication())

	router.GET("/addtocart", app.AddToCart())
	router.GET("/removeitem", app.RemoveFromCart())
	router.GET("/cartcheckout", app.BuyFromCart())
	router.GET("/instantbuy", app.InstantBuy())

	log.Fatal(router.Run(":" + port))
}
