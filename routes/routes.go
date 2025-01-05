package routes

import (
	"github.com/daiki-trnsk/go-ecommerce/controllers"
	"github.com/labstack/echo/v4"
)

func UserRoutes(e *echo.Echo) {
	e.POST("/users/signup", controllers.SignUp)
	e.POST("/users/login", controllers.Login)
	e.POST("/admin/addproduct", controllers.ProductViewerAdmin)
	e.GET("/users/productview", controllers.SearchProduct)
	e.GET("/users/search", controllers.SearchProductByQuery)
}