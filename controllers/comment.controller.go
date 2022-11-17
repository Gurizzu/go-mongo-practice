package controllers

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"go-mongo-practice/database"
	"go-mongo-practice/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"net/http"
	"time"
)

var CommentCollection *mongo.Collection = database.OpenCollection(database.Client, "comment")

// CreateComment godoc
// @Summary      answer Notes
// @Description  answer note
// @Tags         Comment
// @Accept       json
// @Param        id   path      string  true  "Note ID"
// @Param request body models.CommentModels true "query params"
// @Produce      json
// @Success      200  {object}  models.Note
// @Router       /api/comment/create/{id} [post]
// @securityDefinitions.apiKey token
// @in header
// @name Authorization
// @Security JWT
func CreateComment() gin.HandlerFunc {
	return func(c *gin.Context) {

		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		var comment models.Comment
		var note models.Note
		noteid := c.Param("note_id")
		defer cancel()

		err := dbCollection.FindOne(ctx, bson.M{"note_id": noteid}).Decode(&note)
		if err != nil {
			c.JSON(404, gin.H{"error": "id not found"})
		}

		if err := c.BindJSON(&comment); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		validationErr := validate.Struct(comment)
		if validationErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": validationErr.Error()})
			return
		}

		comment.Created_at = time.Now().Local().Unix()
		comment.Updated_at = time.Now().Local().Unix()
		comment.ID = primitive.NewObjectID()
		comment.Note_id = note.Note_id
		comment.Author = models.CommentAuthor{c.GetString("uid"), c.GetString("email")}

		_, err = CommentCollection.InsertOne(ctx, comment)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		} else {
			c.JSON(200, comment)
			return
		}

	}
}

// GetComment godoc
// @Summary      get comment
// @Description  get comment
// @Tags         Comment
// @Accept       json
// @Param        id   path      string  true  "Comment ID"
// @Produce      json
// @Success      200  {object}  models.Comment
// @Router       /api/comment/get/{id} [get]
// @securityDefinitions.apiKey token
// @in header
// @name Authorization
// @Security JWT
func GetComment() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		commentId := c.Param("comment_id")
		defer cancel()
		objId, _ := primitive.ObjectIDFromHex(commentId)

		var commentFind models.Comment
		err := CommentCollection.FindOne(ctx, bson.M{"_id": objId}).Decode(&commentFind)
		if err != nil {
			c.JSON(404, gin.H{"error": "id not found"})
		} else {
			c.JSON(200, commentFind)
		}

	}
}

// GetCommentAll godoc
// @Summary      get all comment
// @Description  get Comment
// @Tags         Comment
// @Produce      json
// @Success      200  {object}  models.Comment
// @Router       /api/comment/get_all [get]
// @securityDefinitions.apiKey token
// @in header
// @name Authorization
// @Security JWT
func GetCommentAll() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		comment, err := CommentCollection.Find(ctx, bson.D{})
		if err != nil {
			c.JSON(404, gin.H{"error": err.Error()})
		}
		var result []models.Comment
		for comment.Next(ctx) {
			var commentFind models.Comment
			err := comment.Decode(&commentFind)
			if err != nil {
				c.JSON(404, gin.H{"error": err.Error()})
			}

			result = append(result, commentFind)

		}
		c.JSON(200, result)
	}
}

// PutComment godoc
// @Summary      Update Comment
// @Description  update Comment
// @Tags         Comment
// @Accept       json
// @Param        id   path      string  true  "Comment ID"
// @Param request body models.CommentModels true "query params"
// @Produce      json
// @Success      200  {object}  models.Comment
// @Router       /api/comment/update/{id} [put]
// @securityDefinitions.apiKey token
// @in header
// @name Authorization
// @Security JWT
func PutComment() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		commentId := c.Param("comment_id")
		var comment models.Comment
		var commentFind models.Comment
		defer cancel()
		objId, _ := primitive.ObjectIDFromHex(commentId)

		err := CommentCollection.FindOne(ctx, bson.M{"_id": objId}).Decode(&commentFind)
		if err != nil {
			c.JSON(404, gin.H{"error": "id not found"})
			return
		}

		if err := c.BindJSON(&comment); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		validationErr := validate.Struct(comment)
		if validationErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": validationErr.Error()})
			return
		}
		update := bson.M{"answer": comment.Answer, "updated_at": time.Now().Local().Unix()}
		result, err := CommentCollection.UpdateOne(ctx, bson.M{"_id": objId}, bson.M{"$set": update})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		var updatedComment models.Comment
		if result.MatchedCount == 1 {
			err := CommentCollection.FindOne(ctx, bson.M{"_id": objId}).Decode(&updatedComment)
			if err != nil {
				c.JSON(404, gin.H{"error": err.Error()})
				return
			}
			c.JSON(200, updatedComment)
			return
		}

	}
}

// DeleteComment godoc
// @Summary      Delete comment
// @Description  delete the comment
// @Tags         Comment
// @Accept       json
// @Param        id   path      string  true  "Comment ID"
// @Produce      json
// @Success      200              {string}  string    "ok"
// @Router       /api/comment/delete/{id} [delete]
// @securityDefinitions.apiKey token
// @in header
// @name Authorization
// @Security JWT
func DeleteComment() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		commentId := c.Param("comment_id")
		var commentFind models.Comment
		defer cancel()
		objId, _ := primitive.ObjectIDFromHex(commentId)

		err := CommentCollection.FindOne(ctx, bson.M{"_id": objId}).Decode(&commentFind)
		if err != nil {
			c.JSON(404, gin.H{"error": "id not found"})
			return
		}

		_, err = CommentCollection.DeleteOne(ctx, bson.M{"_id": objId})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}
}

// GetCommentByNotes godoc
// @Summary      get all comment
// @Description  get Comment
// @Tags         Comment
// @Produce      json
// @Param		 id   path      string  true  "Comment by note ID"
// @Success      200  {object}  models.Comment
// @Router       /api/comment/get_all_by_notes/{id} [get]
// @securityDefinitions.apiKey token
// @in header
// @name Authorization
// @Security JWT
func GetCommentByNotes() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()
		idNotes := c.Param("note_id")
		fmt.Println(idNotes)

		comment, err := CommentCollection.Find(ctx, bson.M{"note_id": idNotes})
		if err != nil {
			c.JSON(404, gin.H{"error": err.Error()})
		}
		fmt.Println(comment)
		var result []models.Comment
		for comment.Next(ctx) {
			var commentFind models.Comment
			fmt.Println(commentFind)
			err := comment.Decode(&commentFind)
			if err != nil {
				c.JSON(404, gin.H{"error": err.Error()})
			}

			result = append(result, commentFind)

		}
		c.JSON(200, result)
	}
}
