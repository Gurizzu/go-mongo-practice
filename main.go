package main

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"go-mongo-practice/docs"
	"go-mongo-practice/routes"
)

// @title           Tanyaa Project
// @version         1.0
// @description     This is a sample server celler server.

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      localhost:2005

// @securityDefinitions.apiKey JWT
// @in header
// @name token
func main() {

	router := gin.New()
	//router.Use(cors.Default())
	config := cors.DefaultConfig()
	config.AllowAllOrigins = true
	config.AllowCredentials = true
	config.AddAllowMethods("OPTIONS")
	router.Use(cors.New(config))
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	docs.SwaggerInfo.BasePath = ""

	public := router.Group("/api/")
	routes.AuthRoutes(public)

	private := router.Group("/api/")
	routes.NotesRouter(private)
	routes.UserRoutes(private)
	routes.CommentRoutes(private)

	err := router.Run(":2005")
	if err != nil {
		return
	}

}
