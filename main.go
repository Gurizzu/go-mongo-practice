package main

import (
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	docs "go-mongo-practice/docs"
	"go-mongo-practice/routes"
)

// @title           Swagger Example API
// @version         1.0
// @description     This is a sample server celler server.

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      localhost:2005

func main() {

	router := gin.New()
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	docs.SwaggerInfo.BasePath = ""
	routes.NotesRouter(router)
	router.Run(":2005")

}
