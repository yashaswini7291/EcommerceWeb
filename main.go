package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/yashaswini7291/ecommerceWeb/controllers"
	"github.com/yashaswini7291/ecommerceWeb/database"
	"github.com/yashaswini7291/ecommerceWeb/middleware"
	"github.com/yashaswini7291/ecommerceWeb/routes"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8000" //defaulting to 8000 in case empty port
	}

	app := controllers.NewApplication(database.ProductData(database.Client, "Products"), database.UserData(database.Client, "Users"))

	router := gin.New()
	router.Use(gin.Logger())

	routes.UserRoutes(router)
	router.Use(middleware.Authentication())

	router.GET("/addtocart", app.AddToCart())
	router.GET("/removeitem", app.RemoveCartItem())
	router.GET("/listcart", controllers.GetItemFromCart())
	router.POST("/addaddress", controllers.AddAddress())
	router.PUT("/editaddress", controllers.EditAddress())
	router.GET("/deleteaddress", controllers.DeleteAddress())
	router.GET("/cartcheckout", app.BuyFromCart())
	router.GET("/instantbuy", app.InstantBuy())

	log.Fatal(router.Run(":" + port))

}
