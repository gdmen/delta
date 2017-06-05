package api

import (
	"gopkg.in/gin-gonic/gin.v1"
)

func GetRouterV1() *gin.Engine {
	router := gin.Default()

	v1 := router.Group("/api/v1")
	{
		user := v1.Group("/u")
		{
			// Handle POST requests at /u/register
			// Ensure that the user is not logged in by using the middleware
			user.POST("/register", ensureNotLoggedIn(), registerUser)
		}
	}
	return router
}
