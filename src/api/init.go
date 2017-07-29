package api

import (
	"database/sql"
	"gopkg.in/gin-gonic/gin.v1"
)

var CREATE_TABLES_SQL = []string{
	CreateUserTableSQL,
	CreateMeasurementTypeTableSQL,
}

type Api struct {
	DB *sql.DB
}

func NewApi(db *sql.DB) (*Api, error) {
	for _, sql := range CREATE_TABLES_SQL {
		_, err := db.Exec(sql)
		if err != nil {
			return nil, err
		}
	}
	return &Api{DB: db}, nil
}

func (a *Api) GetRouter() *gin.Engine {
	router := gin.Default()

	v1 := router.Group("/api/v1")
	{
		user := v1.Group("/users")
		{
			user.POST("/register", RequestIdMiddleware(), ensureNotLoggedIn(), a.registerUser)
		}
		measurementType := v1.Group("/measurement_types")
		{
			measurementType.POST("/", RequestIdMiddleware(), a.createMeasurementType)
			measurementType.POST("/:id", RequestIdMiddleware(), a.updateMeasurementType)
			measurementType.DELETE("/:id", RequestIdMiddleware(), a.deleteMeasurementType)
			measurementType.GET("/:id", RequestIdMiddleware(), a.getMeasurementType)
			measurementType.GET("/", RequestIdMiddleware(), a.listMeasurementType)
		}
	}
	return router
}
