package routes

import (
	"github.com/gin-gonic/gin"
	"go-mongo-practice/controllers"
	"go-mongo-practice/middleware"
)

func UserRoutes(c *gin.RouterGroup) {
	c.Use(middleware.Authentication())
	c.Use(gin.Logger())

	c.GET("/user/get_all", controllers.GetUsers())
	c.GET("/user/get/:user_id", controllers.GetUser())
	c.DELETE("/user/delete/:user_id", controllers.DeleteUser())
	c.PUT("/user/update/:user_id", controllers.UpdateUser())

}
