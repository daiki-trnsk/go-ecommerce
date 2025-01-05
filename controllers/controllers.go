package controllers

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/daiki-trnsk/go-ecommerce/database"
	"github.com/daiki-trnsk/go-ecommerce/models"
	generate "github.com/daiki-trnsk/go-ecommerce/tokens"
	"github.com/labstack/echo/v4"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

var UserCollection *mongo.Collection = database.UserData(database.Client, "Users")
var ProductCollection *mongo.Collection = database.ProductData(database.Client, "Products")
var Validate = validator.New()

func HashPassword(password string) string {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		log.Panic(err)
	}
	return string(bytes)
}

func VerifyPassword(userpassword string, givenpassword string) (bool, string) {
	err := bcrypt.CompareHashAndPassword([]byte(givenpassword), []byte(userpassword))
	valid := true
	msg := ""
	if err != nil {
		msg = "Login Or Passowrd is Incorerct"
		valid = false
	}
	return valid, msg
}

func SignUp(c echo.Context) error{
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()
		var user models.User
		if err := c.Bind(&user); err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
		}
		validationErr := Validate.Struct(user)
		if validationErr != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": validationErr.Error()})
		}

		count, err := UserCollection.CountDocuments(ctx, bson.M{"email": user.Email})
		if err != nil {
			log.Panic(err)
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
		}
		if count > 0 {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "User already exists"})
		}
		count, err = UserCollection.CountDocuments(ctx, bson.M{"phone": user.Phone})
		defer cancel()
		if err != nil {
			log.Panic(err)
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
		}
		if count > 0 {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "this phone number is already in use"})
		}
		password := HashPassword(*user.Password)
		user.Password = &password

		user.Created_At, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		user.Updated_At, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		user.ID = primitive.NewObjectID()
		user.User_ID = user.ID.Hex()
		token, refreshtoken, _ := generate.TokenGenerator(*user.Email, *user.First_Name, *user.Last_Name, user.User_ID)
		user.Token = &token
		user.Refresh_Token = &refreshtoken
		user.UserCart = make([]models.ProductUser, 0)
		user.Address_Details = make([]models.Address, 0)
		user.Order_Status = make([]models.Order, 0)
		_, inserterr := UserCollection.InsertOne(ctx, user)
		if inserterr != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "the user did not created"})
		}
		defer cancel()
		return c.JSON(http.StatusCreated, "Succesfully signed in!")
}

func Login(c echo.Context) error{
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()
		var user models.User
		var founduser models.User
		if err := c.Bind(&user); err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
		}
		err := UserCollection.FindOne(ctx, bson.M{"email": user.Email}).Decode(&founduser)
		defer cancel()
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "User not found"})
		}
		PasswordIsValid, msg := VerifyPassword(*user.Password, *founduser.Password)
		defer cancel()
		if !PasswordIsValid {
			fmt.Println(msg)
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": msg})
		}
		token, refreshToken, _ := generate.TokenGenerator(*founduser.Email, *founduser.First_Name, *founduser.Last_Name, founduser.User_ID)
		defer cancel()
		generate.UpdateAllTokens(token, refreshToken, founduser.User_ID)
		return c.JSON(http.StatusFound, founduser)
	}

func ProductViewerAdmin(c echo.Context) error {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		var products models.Product
		defer cancel()
		if err := c.Bind(&products); err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
		}
		products.Product_ID = primitive.NewObjectID()
		_, anyerr := ProductCollection.InsertOne(ctx, products)
		if anyerr != nil {
			return c.JSON(http.StatusInternalServerError, "Not Created")
		}
		defer cancel()
		return c.JSON(http.StatusOK, "Successfully added our Product Admin!!")
	}

func SearchProduct(c echo.Context) error{
		var productlist []models.Product
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()
		cursor, err := ProductCollection.Find(ctx, bson.D{{}})
		if err != nil {
			return c.JSON(http.StatusInternalServerError, "Someting Went Wrong Please Try After Some Time")
		}
		err = cursor.All(ctx, &productlist)
		if err != nil {
			log.Println(err)
			return c.JSON(http.StatusInternalServerError, "Error occurred while fetching products")
		}
		defer cursor.Close(ctx)
		if err := cursor.Err(); err != nil {
			// Don't forget to log errors. I log them really simple here just
			// to get the point across.
			log.Println(err)
			return c.JSON(400, "invalid")
		}
		defer cancel()
		return c.JSON(200, productlist)
	}

func SearchProductByQuery(c echo.Context) error{
		var searchproducts []models.Product
		queryParam := c.QueryParam("name")
		if queryParam == "" {
			log.Println("query is empty")
			return c.JSON(http.StatusNotFound, map[string]string{"error": "Invalid Search Index"})
		}
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()
		searchquerydb, err := ProductCollection.Find(ctx, bson.M{"product_name": bson.M{"$regex": queryParam, "$options": "i"}})
		if err != nil {
			return c.JSON(404, "something went wrong in fetching the dbquery")
		}
		err = searchquerydb.All(ctx, &searchproducts)
		if err != nil {
			log.Println(err)
			return c.JSON(400, "invalid")
		}
		defer searchquerydb.Close(ctx)
		if err := searchquerydb.Err(); err != nil {
			log.Println(err)
			return c.JSON(400, "invalid request")
		}
		defer cancel()
		return c.JSON(200, searchproducts)
	}
