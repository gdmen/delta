package api

import (
	"database/sql"
	"gopkg.in/gin-gonic/gin.v1"
)

type Api struct {
	DB *sql.DB
}

func NewApi(db *sql.DB) *Api {
	return &Api{DB: db}
}

func (a *Api) GetRouter() *gin.Engine {
	router := gin.Default()

	v1 := router.Group("/api/v1")
	{
		user := v1.Group("/u")
		{
			user.POST("/register", ensureNotLoggedIn(), a.registerUser)
		}
	}
	return router
}
