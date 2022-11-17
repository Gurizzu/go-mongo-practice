package routes

import (
	"github.com/gin-gonic/gin"
	"go-mongo-practice/controllers"
)

func AuthRoutes(c *gin.RouterGroup) {
	c.POST("/user/signup", controllers.SignUp())
	c.POST("/user/login", controllers.Login())
}
