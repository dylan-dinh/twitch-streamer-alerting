package api

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// SetUpRouter declare our routes
func SetUpRouter(userHandler *UserHandler) *gin.Engine {
	router := gin.Default()

	// Enable CORS
	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"http://localhost:3000"}
	config.AllowMethods = []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"}
	config.AllowHeaders = []string{"Origin", "Content-Type", "Accept", "Authorization"}
	router.Use(cors.New(config))

	userRoutes := router.Group("/api").Group("/users")
	{
		userRoutes.POST("/register", userHandler.InsertUser)
		userRoutes.POST("/login", userHandler.Login)
	}
	return router
}
