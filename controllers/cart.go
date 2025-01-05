package controllers

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/daiki-trnsk/go-ecommerce/database"
	"github.com/daiki-trnsk/go-ecommerce/models"
	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type Application struct {
	prodCollection *mongo.Collection
	userCollection *mongo.Collection
}

func NewApplication(prodCollection, userCollection *mongo.Collection) *Application {
	return &Application{
		prodCollection: prodCollection,
		userCollection: userCollection,
	}
}

func (app *Application) AddToCart(c echo.Context) error {
	productQueryID := c.QueryParam("id")
	if productQueryID == "" {
		log.Println("product id is empty")
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "product id is empty"})
	}
	userQueryID := c.QueryParam("userID")
	if userQueryID == "" {
		log.Println("user id is empty")
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "user id is empty"})
	}
	productID, err := primitive.ObjectIDFromHex(productQueryID)
	if err != nil {
		log.Println(err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "invalid product id"})
	}
	var ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = database.AddProductToCart(ctx, app.prodCollection, app.userCollection, productID, userQueryID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
	}
	return c.JSON(200, "Successfully Added to the cart")
}

func (app *Application) RemoveItem(c echo.Context) error {
	productQueryID := c.QueryParam("id")
	if productQueryID == "" {
		log.Println("product id is inavalid")
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "product id is empty"})
	}

	userQueryID := c.QueryParam("userID")
	if userQueryID == "" {
		log.Println("user id is empty")
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "UserID is empty"})
	}

	ProductID, err := primitive.ObjectIDFromHex(productQueryID)
	if err != nil {
		log.Println(err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "invalid product id"})
	}

	var ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	err = database.RemoveCartItem(ctx, app.prodCollection, app.userCollection, ProductID, userQueryID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}
	return c.JSON(200, "Successfully removed from cart")
}

func GetItemFromCart(c echo.Context) error {
	user_id := c.QueryParam("id")
	if user_id == "" {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "invalid id"})
	}

	usert_id, _ := primitive.ObjectIDFromHex(user_id)

	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	var filledcart models.User
	err := UserCollection.FindOne(ctx, bson.D{primitive.E{Key: "_id", Value: usert_id}}).Decode(&filledcart)
	if err != nil {
		log.Println(err)
		return c.JSON(500, "not id found")
	}

	filter_match := bson.D{{Key: "$match", Value: bson.D{primitive.E{Key: "_id", Value: usert_id}}}}
	unwind := bson.D{{Key: "$unwind", Value: bson.D{primitive.E{Key: "path", Value: "$usercart"}}}}
	grouping := bson.D{{Key: "$group", Value: bson.D{primitive.E{Key: "_id", Value: "$_id"}, {Key: "total", Value: bson.D{primitive.E{Key: "$sum", Value: "$usercart.price"}}}}}}
	pointcursor, err := UserCollection.Aggregate(ctx, mongo.Pipeline{filter_match, unwind, grouping})
	if err != nil {
		log.Println(err)
	}
	var listing []bson.M
	if err = pointcursor.All(ctx, &listing); err != nil {
		log.Println(err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "internal server error"})
	}
	for _, json := range listing {
		c.JSON(200, json["total"])
		c.JSON(200, filledcart.UserCart)
	}
	return nil
}

func (app *Application) BuyFromCart(c echo.Context) error {
	userQueryID := c.QueryParam("id")
	if userQueryID == "" {
		log.Panicln("user id is empty")
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "UserID is empty"})
	}
	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()
	err := database.BuyItemFromCart(ctx, app.userCollection, userQueryID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}
	return c.JSON(200, "Successfully Placed the order")
}

func (app *Application) InstantBuy(c echo.Context) error {
	UserQueryID := c.QueryParam("userid")
	if UserQueryID == "" {
		log.Println("UserID is empty")
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "UserID is empty"})
	}
	ProductQueryID := c.QueryParam("pid")
	if ProductQueryID == "" {
		log.Println("Product_ID id is empty")
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "product_id is empty"})
	}
	productID, err := primitive.ObjectIDFromHex(ProductQueryID)
	if err != nil {
		log.Println(err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "internal server error"})
	}

	var ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	err = database.InstantBuyer(ctx, app.prodCollection, app.userCollection, productID, UserQueryID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}
	return c.JSON(200, "Successully placed the order")
}
