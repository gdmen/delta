package api

import (
	"database/sql"
	"gopkg.in/gin-gonic/gin.v1"
)

type Api struct {
	DB *sql.DB
}

func NewApi(db *sql.DB) (*Api, error) {
	_, err := db.Exec(CreateUserTableSQL)
	if err != nil {
		return nil, err
	}
	return &Api{DB: db}, nil
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
