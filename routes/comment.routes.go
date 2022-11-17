package routes

import (
	"github.com/gin-gonic/gin"
	"go-mongo-practice/controllers"
	"go-mongo-practice/middleware"
)

func CommentRoutes(c *gin.RouterGroup) {
	c.Use(middleware.Authentication())
	c.Use(gin.Logger())

	c.POST("/comment/create/:note_id", controllers.CreateComment())
	c.GET("/comment/get/:comment_id", controllers.GetComment())
	c.GET("/comment/get_all", controllers.GetCommentAll())
	c.PUT("/comment/update/:comment_id", controllers.PutComment())
	c.DELETE("/comment/delete/:comment_id", controllers.DeleteComment())
	c.GET("/comment/get_all_by_notes/:note_id", controllers.GetCommentByNotes())
}
