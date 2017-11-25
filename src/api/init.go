package api

import (
	"database/sql"
	"strings"

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
			if strings.Contains(err.Error(), "already exists") {
				continue
			}
			return nil, err
		}
	}
	return &Api{DB: db}, nil
}

func (a *Api) GetRouter() *gin.Engine {
	router := gin.Default()
	// Allow all origins, methods
	router.Use(cors.Default())

	v1 := router.Group("/api/v1")
	{
		measurementType := v1.Group("/measurement_types")
		{
			measurementType.POST("", common.RequestIdMiddleware(), a.createMeasurementType)
			measurementType.POST("/", common.RequestIdMiddleware(), a.createMeasurementType)
			measurementType.POST("/:id", common.RequestIdMiddleware(), a.updateMeasurementType)
			measurementType.DELETE("/:id", common.RequestIdMiddleware(), a.deleteMeasurementType)
			measurementType.GET("/:id", common.RequestIdMiddleware(), a.getMeasurementType)
			measurementType.GET("", common.RequestIdMiddleware(), a.listMeasurementType)
			measurementType.GET("/", common.RequestIdMiddleware(), a.listMeasurementType)
		}
		measurement := v1.Group("/measurements")
		{
			measurement.POST("", common.RequestIdMiddleware(), a.createMeasurement)
			measurement.POST("/", common.RequestIdMiddleware(), a.createMeasurement)
			measurement.POST("/:id", common.RequestIdMiddleware(), a.updateMeasurement)
			measurement.DELETE("/:id", common.RequestIdMiddleware(), a.deleteMeasurement)
			measurement.GET("/:id", common.RequestIdMiddleware(), a.getMeasurement)
			measurement.GET("", common.RequestIdMiddleware(), a.listMeasurement)
			measurement.GET("/", common.RequestIdMiddleware(), a.listMeasurement)
		}
		importFiles := v1.Group("/import")
		{
			importFiles.POST("/fitnotes", common.RequestIdMiddleware(), a.importFitnotes)
			importFiles.POST("/fitocracy", common.RequestIdMiddleware(), a.importFitocracy)
			importFiles.POST("/strava", common.RequestIdMiddleware(), a.importStrava)
		}
		data := v1.Group("/data")
		{
			data.GET("/drilldown", common.RequestIdMiddleware(), a.getDrilldown)
			data.GET("/maxes", common.RequestIdMiddleware(), a.getMaxes)
		}
	}
	return router
}
