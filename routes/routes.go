package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/yashaswini7291/ecommerceWeb/controllers"
)

func UserRoutes(inRoute *gin.Engine) {
	inRoute.POST("/users/signup", controllers.SignUp())
	inRoute.POST("/users/login", controllers.Login())
	inRoute.POST("/admin/addproduct", controllers.ProductViewerAdmin())
	inRoute.GET("/users/productView", controllers.SearchProduct())
	inRoute.GET("/users/search", controllers.SearchProductByQuery())
}
