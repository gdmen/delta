package api

import (
	"net/http"

	"gopkg.in/gin-gonic/gin.v1"
)

func GetRouterV1() *gin.Engine {
	router := gin.Default()

	v1 := router.Group("/api/v1")
	{
		v1.GET("/welcome", func(c *gin.Context) {
			firstname := c.DefaultQuery("firstname", "Guest")
			lastname := c.Query("lastname") // shortcut for c.Request.URL.Query().Get("lastname")

			c.String(http.StatusOK, "Hello %s %s", firstname, lastname)
		})
	}
	return router
}
