package controllers

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/daiki-trnsk/go-ecommerce/models"
	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func AddAddress(c echo.Context) error {
		user_id := c.QueryParam("id")
		if user_id == "" {
			return c.JSON(http.StatusNotFound, map[string]string{"error": "Invalid code"})
		}
		address, err := primitive.ObjectIDFromHex(user_id)
		if err != nil {
			return c.JSON(500, "Internal Server Error")
		}
		var addresses models.Address
		addresses.Address_id = primitive.NewObjectID()
		if err = c.Bind(&addresses); err != nil {
			return c.JSON(http.StatusNotAcceptable, map[string]string{"error": err.Error()})
		}

		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		match_filter := bson.D{{Key: "$match", Value: bson.D{primitive.E{Key: "_id", Value: address}}}}
		unwind := bson.D{{Key: "$unwind", Value: bson.D{primitive.E{Key: "path", Value: "$address"}}}}
		group := bson.D{{Key: "$group", Value: bson.D{primitive.E{Key: "_id", Value: "$address_id"}, {Key: "count", Value: bson.D{primitive.E{Key: "$sum", Value: 1}}}}}}

		pointcursor, err := UserCollection.Aggregate(ctx, mongo.Pipeline{match_filter, unwind, group})
		if err != nil {
			return c.JSON(500, "Internal Server Error")
		}

		var addressinfo []bson.M
		if err = pointcursor.All(ctx, &addressinfo); err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Internal Server Error"})
		}

		var size int32
		for _, address_no := range addressinfo {
			count := address_no["count"]
			size = count.(int32)
		}
		if size < 2 {
			filter := bson.D{primitive.E{Key: "_id", Value: address}}
			update := bson.D{{Key: "$push", Value: bson.D{primitive.E{Key: "address", Value: addresses}}}}
			_, err := UserCollection.UpdateOne(ctx, filter, update)
			if err != nil {
				fmt.Println(err)
			}
		} else {
			c.JSON(400, "Not Allowed ")
		}
		return c.JSON(http.StatusOK, map[string]string{"message": "Successfully added address"})
	}
func EditHomeAddress(c echo.Context) error {
		user_id := c.QueryParam("id")
		if user_id == "" {
			return c.JSON(http.StatusNotFound, map[string]string{"Error": "Invalid"})
		}
		usert_id, err := primitive.ObjectIDFromHex(user_id)
		if err != nil {
			c.JSON(500, err)
		}
		var editaddress models.Address
		if err := c.Bind(&editaddress); err != nil {
			return c.JSON(http.StatusBadRequest, err.Error())
		}
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()
		filter := bson.D{primitive.E{Key: "_id", Value: usert_id}}
		update := bson.D{{Key: "$set", Value: bson.D{primitive.E{Key: "address.0.house_name", Value: editaddress.House}, {Key: "address.0.street_name", Value: editaddress.Street}, {Key: "address.0.city_name", Value: editaddress.City}, {Key: "address.0.pin_code", Value: editaddress.Pincode}}}}
		_, err = UserCollection.UpdateOne(ctx, filter, update)
		if err != nil {
			return c.JSON(500, "Something Went Wrong")
		}
		defer cancel()
		ctx.Done()
		return c.JSON(200, "Successfully Updated the Home address")
	}

func EditWorkAddress(c echo.Context) error {
		user_id := c.QueryParam("id")
		if user_id == "" {
			return c.JSON(http.StatusNotFound, map[string]string{"Error": "Wrong id not provided"})
		}
		usert_id, err := primitive.ObjectIDFromHex(user_id)
		if err != nil {
			return c.JSON(500, err)
		}
		var editaddress models.Address
		if err := c.Bind(&editaddress); err != nil {
			return c.JSON(http.StatusBadRequest, err.Error())
		}
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()
		filter := bson.D{primitive.E{Key: "_id", Value: usert_id}}
		update := bson.D{{Key: "$set", Value: bson.D{primitive.E{Key: "address.1.house_name", Value: editaddress.House}, {Key: "address.1.street_name", Value: editaddress.Street}, {Key: "address.1.city_name", Value: editaddress.City}, {Key: "address.1.pin_code", Value: editaddress.Pincode}}}}
		_, err = UserCollection.UpdateOne(ctx, filter, update)
		if err != nil {
			return c.JSON(500, "something Went wrong")
		}
		defer cancel()
		ctx.Done()
		return c.JSON(200, "Successfully updated the Work Address")
	}

func DeleteAddress(c echo.Context) error {
		user_id := c.QueryParam("id")
		if user_id == "" {
			return c.JSON(http.StatusNotFound, map[string]string{"Error": "Invalid Search Index"})
		}
		addresses := make([]models.Address, 0)
		usert_id, err := primitive.ObjectIDFromHex(user_id)
		if err != nil {
			return c.JSON(500, "Internal Server Error")
		}
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()
		filter := bson.D{primitive.E{Key: "_id", Value: usert_id}}
		update := bson.D{{Key: "$set", Value: bson.D{primitive.E{Key: "address", Value: addresses}}}}
		_, err = UserCollection.UpdateOne(ctx, filter, update)
		if err != nil {
			return c.JSON(404, "Wromg")
			
		}
		defer cancel()
		ctx.Done()
		return c.JSON(200, "Successfully Deleted!")
	}