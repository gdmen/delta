package api

import (
	"database/sql"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	"github.com/gdmen/delta/src/common"
)

var CREATE_TABLES_SQL = []string{
	CreateMeasurementTypeTableSQL,
	CreateMeasurementTableSQL,
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
	// Allow all origins, methods
	config := cors.Default()
	router.Use(config)

	v1 := router.Group("/api/v1")
	{
		importFiles := v1.Group("/import")
		{
			importFiles.POST("/fitnotes", common.RequestIdMiddleware(), a.importFitnotes)
		}
		measurementType := v1.Group("/measurement_types")
		{
			measurementType.POST("/", common.RequestIdMiddleware(), a.createMeasurementType)
			measurementType.POST("/:id", common.RequestIdMiddleware(), a.updateMeasurementType)
			measurementType.DELETE("/:id", common.RequestIdMiddleware(), a.deleteMeasurementType)
			measurementType.GET("/:id", common.RequestIdMiddleware(), a.getMeasurementType)
			measurementType.GET("/", common.RequestIdMiddleware(), a.listMeasurementType)
		}
		measurement := v1.Group("/measurements")
		{
			measurement.POST("/", common.RequestIdMiddleware(), a.createMeasurement)
			measurement.POST("/:id", common.RequestIdMiddleware(), a.updateMeasurement)
			measurement.DELETE("/:id", common.RequestIdMiddleware(), a.deleteMeasurement)
			measurement.GET("/:id", common.RequestIdMiddleware(), a.getMeasurement)
			measurement.GET("/", common.RequestIdMiddleware(), a.listMeasurement)
		}
	}
	return router
}
