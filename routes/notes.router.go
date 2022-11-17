package routes

import (
	"github.com/gin-gonic/gin"
	"go-mongo-practice/controllers"
	"go-mongo-practice/middleware"
)

func NotesRouter(c *gin.RouterGroup) {
	c.Use(middleware.Authentication())
	c.Use(gin.Logger())

	c.POST("/notes/create", controllers.CreateNotes())
	c.PUT("/notes/update/:notes_id", controllers.UpdateNotes())
	c.DELETE("/notes/delete/:notes_id", controllers.DeleteNotes())
	c.GET("/notes/get_note/:notes_id", controllers.GetNote())
	c.GET("/notes/get_notes", controllers.GetNotes())

}
