package api

import "github.com/gin-gonic/gin"

// SetUpRouter declare our routes
func SetUpRouter(userHandler *UserHandler) *gin.Engine {
	router := gin.Default()

	userRoutes := router.Group("/users")
	{
		userRoutes.POST("/", userHandler.InsertUser)
	}
	return router
}
