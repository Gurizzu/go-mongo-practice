package controllers

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go-mongo-practice/database"
	"go-mongo-practice/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"net/http"
	"time"
)

var userCollection *mongo.Collection = database.OpenCollection(database.Client, "user")
var validate = validator.New()

// GetUser godoc
// @Summary      Register
// @Description  create account
// @Tags         Users
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "Account ID"
// @Success      200  {object}  models.User
// @Router       /api/user/get/{id} [get]
// @securityDefinitions.apiKey token
// @in header
// @name Authorization
// @Security JWT
func GetUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		userId := c.Param("user_id")
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)

		var userFind models.User
		err := userCollection.FindOne(ctx, bson.M{"user_id": userId}).Decode(&userFind)
		defer cancel()
		if err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}
		c.JSON(200, userFind)
	}
}

// GetUsers godoc
// @Summary      Register
// @Description  create account
// @Tags         Users
// @Produce      json
// @Success      200  {object}  models.User
// @Router       /api/user/get_all [get]
// @securityDefinitions.apiKey token
// @in header
// @name Authorization
// @Security JWT
func GetUsers() gin.HandlerFunc {
	return func(c *gin.Context) {

		//err := helpers.CheckUserType(c, "ADMIN")
		//if err != nil {
		//	c.JSON(404, gin.H{"error": "you dont have permission"})
		//	return
		//}

		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)

		users, err := userCollection.Find(ctx, bson.D{})
		defer cancel()
		if err != nil {
			c.JSON(500, gin.H{"error": "error occured while find users"})
		}
		var allUsers []models.User
		for users.Next(ctx) {
			var userFind models.User
			err := users.Decode(&userFind)
			if err != nil {
				c.JSON(404, gin.H{"error": err.Error()})
			}

			allUsers = append(allUsers, userFind)
		}

		c.JSON(200, allUsers)
	}
}

// DeleteUser godoc
// @Summary      Delete user
// @Description  Delete account
// @Tags         Users
// @Accept       json
// @Param        id   path      string  true  "Account ID"
// @Produce      json
// @Success      200  {string}  string    "ok"
// @Router       /api/user/delete/{id} [delete]
// @securityDefinitions.apiKey token
// @in header
// @name Authorization
// @Security JWT
func DeleteUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		userId := c.Param("user_id")
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		var userFind models.Note
		defer cancel()

		err := userCollection.FindOne(ctx, bson.M{"user_id": userId}).Decode(&userFind)
		if err != nil {
			c.JSON(404, gin.H{"error": "user not found"})
			return
		}

		_, err = dbCollection.DeleteMany(ctx, bson.M{"author.user_id": userId})
		_, err = userCollection.DeleteOne(ctx, bson.M{"user_id": userId})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

	}
}

// UpdateUser godoc
// @Summary      Register
// @Description  create account
// @Tags         Users
// @Accept       json
// @Param        id   path      string  true  "Account ID"
// @Param request body models.UserUpdate true "query params"
// @Produce      json
// @Success      200  {object}  models.UserUpdate
// @Router       /api/user/update/{id} [put]
// @securityDefinitions.apiKey token
// @in header
// @name Authorization
// @Security JWT
func UpdateUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		userId := c.Param("user_id")
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		var user models.UserUpdate
		var userFind models.User
		defer cancel()

		err := userCollection.FindOne(ctx, bson.M{"user_id": userId}).Decode(&userFind)
		if err != nil {
			c.JSON(404, gin.H{"error": "user not found"})
			return
		}

		if err = c.BindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		validationErr := validate.Struct(user)
		if validationErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": validationErr.Error()})
			return
		}

		if user.First_name == nil {
			user.First_name = userFind.First_name
		}
		if user.Last_name == nil {
			user.Last_name = userFind.Last_name
		}
		if user.Email == nil {
			user.Email = userFind.Email
		}
		if user.Phone == nil {
			user.Phone = userFind.Phone
		}
		if user.Password == nil {
			user.Password = userFind.Password
		}

		count, err := userCollection.CountDocuments(ctx, bson.M{"email": user.Email})
		defer cancel()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error occured while checking for the email"})
			return
		}
		count, err = userCollection.CountDocuments(ctx, bson.M{"phone": user.Phone})
		defer cancel()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error occured while checking for the phone number"})
			return
		}

		if count > 0 {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "this email or phone number already exists"})
			return
		}

		update := bson.M{"first_name": user.First_name, "last_name": user.Last_name, "email": user.Email, "phone": user.Phone, "password": user.Password, "updated_at": time.Now().Local().Unix()}

		result, err := userCollection.UpdateOne(ctx, bson.M{"user_id": userId}, bson.M{"$set": update})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		var updatedUser models.UserUpdate
		if result.MatchedCount == 1 {
			err := userCollection.FindOne(ctx, bson.M{"user_id": userId}).Decode(&updatedUser)
			if err != nil {
				c.JSON(404, gin.H{"error": err.Error()})
				return
			}
			c.JSON(200, updatedUser)
		}
	}
}
