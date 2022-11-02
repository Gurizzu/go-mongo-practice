package controllers

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go-mongo-practice/database"
	"go-mongo-practice/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"net/http"
	"time"
)

var dbCollection *mongo.Collection = database.OpenCollection(database.Client, "notes")
var validate = validator.New()

// CreateNotes godoc
// @Summary      Create notes
// @Description  create notes
// @Tags         Notes
// @Accept       json
// @Produce      json
// @Param request body models.CreateNote true "query params"
// @Success      200  {object}  models.Note
// @Router       /notes/create [post]
func CreateNotes() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		var note models.Note

		if err := c.BindJSON(&note); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		validationErr := validate.Struct(note)
		if validationErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "validationErr.Error()"})
			return
		}

		note.Created_at = time.Now().Local().Unix()
		note.Updated_at = time.Now().Local().Unix()
		note.ID = primitive.NewObjectID()
		note.Note_id = note.ID.Hex()

		_, err := dbCollection.InsertOne(ctx, note)
		defer cancel()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(200, note)

	}
}

// GetNote godoc
// @Summary      get note
// @Description  get note
// @Tags         Notes
// @Accept       json
// @Param        id   path      string  true  "Account ID"
// @Produce      json
// @Success      200  {object}  models.Note
// @Router       /notes/get_note/{id} [get]
func GetNote() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		notesId := c.Param("notes_id")
		defer cancel()

		var noteFind models.Note
		err := dbCollection.FindOne(ctx, bson.M{"note_id": notesId}).Decode(&noteFind)
		if err != nil {
			c.JSON(404, gin.H{"error": "id not found"})
		} else {
			c.JSON(200, noteFind)
		}

	}
}

// GetNotes godoc
// @Summary      get all notes
// @Description  get notes
// @Tags         Notes
// @Accept       json
// @Produce      json
// @Success      200  {object}  models.Note
// @Router       /notes/get_notes [get]
func GetNotes() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		notes, err := dbCollection.Find(ctx, bson.D{})
		if err != nil {
			c.JSON(404, gin.H{"error": err.Error()})
		}

		var result []models.Note
		for notes.Next(ctx) {
			var noteFind models.Note
			err := notes.Decode(&noteFind)
			if err != nil {
				c.JSON(404, gin.H{"error": err.Error()})
			}

			result = append(result, noteFind)

		}

		c.JSON(200, result)

	}
}

// UpdateNotes godoc
// @Summary      Update Notes
// @Description  update note
// @Tags         Notes
// @Accept       json
// @Param        id   path      string  true  "Account ID"
// @Param request body models.CreateNote true "query params"
// @Produce      json
// @Success      200  {object}  models.Note
// @Router       /notes/update/{id} [put]
func UpdateNotes() gin.HandlerFunc {
	return func(c *gin.Context) {

		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		notesId := c.Param("notes_id")
		var note models.Note
		var noteFind models.Note
		defer cancel()

		err := dbCollection.FindOne(ctx, bson.M{"note_id": notesId}).Decode(&noteFind)
		if err != nil {
			c.JSON(404, gin.H{"error": "id not found"})
			return
		}

		if err := c.BindJSON(&note); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		validationErr := validate.Struct(note)
		if validationErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": validationErr.Error()})
			return
		}

		update := bson.M{"title": note.Title, "description": note.Description, "updated_at": time.Now().Local().Unix()}

		result, err := dbCollection.UpdateOne(ctx, bson.M{"note_id": notesId}, bson.M{"$set": update})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		var updatedNotes models.Note
		if result.MatchedCount == 1 {
			err := dbCollection.FindOne(ctx, bson.M{"note_id": notesId}).Decode(&updatedNotes)
			if err != nil {
				c.JSON(404, gin.H{"error": err.Error()})
				return
			}
			c.JSON(200, updatedNotes)
		}

	}
}

// DeleteNotes godoc
// @Summary      Delete note
// @Description  delete the note
// @Tags         Notes
// @Accept       json
// @Param        id   path      string  true  "Account ID"
// @Produce      json
// @Success      200              {string}  string    "ok"
// @Router       /notes/delete/{id} [delete]
func DeleteNotes() gin.HandlerFunc {
	return func(c *gin.Context) {

		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		notesId := c.Param("notes_id")
		var noteFind models.Note
		defer cancel()

		err := dbCollection.FindOne(ctx, bson.M{"note_id": notesId}).Decode(&noteFind)
		if err != nil {
			c.JSON(404, gin.H{"error": "id not found"})
			return
		}

		_, err = dbCollection.DeleteOne(ctx, bson.M{"note_id": notesId})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

	}
}
