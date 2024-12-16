package router

import (
	"fmt"
	"net/http"
	"os"
	"working-day-api/config"
	"working-day-api/internal/container"
	"working-day-api/internal/middlewares"

	"github.com/gin-gonic/gin"

	"github.com/gin-contrib/cors"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func LoadRoutes(config *config.AppVars, container *container.Container) {
	router := route(config.GinMode)

	router.GET("/alive", container.HealthController.AliveHandler)

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	login := router.Group("/v1/login")
	{
		login.POST("", container.AuthController.LoginHandler)
	}

	router.POST("/v1/user", container.UserController.CreateUserHandler)

	user := router.Group("/v1/user", middlewares.Auth())
	{
		user.GET("/:id", container.UserController.GetUserHandler)
	}

	task := router.Group("/v1/task", middlewares.Auth())
	{
		task.GET("/:id", container.TaskController.GetTaskHandler)
		task.GET("", container.TaskController.GetTasksHandler)
		task.POST("", container.TaskController.CreateTaskHandler)
		task.PUT("/:id", container.TaskController.UpdateTaskHandler)
		task.DELETE("/:id", container.TaskController.DeleteTaskHandler)
	}

	router.NoRoute(endpointNotFound)
	router.Run(fmt.Sprint(":" + os.Getenv("PORT")))

}

func endpointNotFound(c *gin.Context) {
	c.JSON(http.StatusNotFound, gin.H{"message": "there's no endpoint for that!"})
}

func route(ginMode string) *gin.Engine {
	router := gin.Default()

	gin.SetMode(ginMode)

	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"*"}
	config.AllowHeaders = []string{"*"}
	config.AllowCredentials = true
	config.AllowBrowserExtensions = true

	config.AddAllowMethods("GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS")

	router.Use(cors.New(config))

	return router
}
