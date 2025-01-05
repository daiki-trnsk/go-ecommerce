package main

import (
	"log"
	"os"

	"github.com/daiki-trnsk/go-ecommerce/controllers"
	"github.com/daiki-trnsk/go-ecommerce/database"
	customMiddleware "github.com/daiki-trnsk/go-ecommerce/middleware"
	"github.com/daiki-trnsk/go-ecommerce/routes"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8000"
	}
	app := controllers.NewApplication(database.ProductData(database.Client, "Products"), database.UserData(database.Client, "Users"))

	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	routes.UserRoutes(e)

	e.POST("/signup", controllers.SignUp)
	e.POST("/login", controllers.Login)

	auth := e.Group("")
	auth.Use(customMiddleware.Authentication)
	auth.GET("/addtocart", app.AddToCart)
	auth.GET("/removeitem", app.RemoveItem)
	auth.GET("/listcart", controllers.GetItemFromCart)
	auth.POST("/addaddress", controllers.AddAddress)
	auth.PUT("/edithomeaddress", controllers.EditHomeAddress)
	auth.PUT("/editworkaddress", controllers.EditWorkAddress)
	auth.GET("/deleteaddresses", controllers.DeleteAddress)
	auth.GET("/cartcheckout", app.BuyFromCart)
	auth.GET("/instantbuy", app.InstantBuy)

	log.Fatal(e.Start(":" + port))
}
