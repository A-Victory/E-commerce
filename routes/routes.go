package routes

import (
	"github.com/A-Victory/e-commerce/controllers"
	"github.com/gin-gonic/gin"
)

func UserRoutes(incomingRoutes *gin.Engine) {
	incomingRoutes.POST("/users/signup", controllers.Signup())
	incomingRoutes.POST("/users/login", controllers.Login())
	incomingRoutes.GET("/users/search", controllers.SearchProductByQuery())
	incomingRoutes.GET("/users/productview", controllers.SearchProduct())
	incomingRoutes.POST("/admin/addproduct", controllers.ProductView())
}
