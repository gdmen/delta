package api

import (
	"database/sql"
	"github.com/golang/glog"
	"net/http"
	"strconv"

	"gopkg.in/gin-gonic/gin.v1"
)

func (a *Api) createMeasurement(c *gin.Context) {
	logPrefix := GetLogPrefix(c)
	glog.Infof("%s fcn start", logPrefix)

	// Parse input
	model := &Measurement{}
	err := c.Bind(model)
	if err != nil {
		msg := "Couldn't parse input form"
		glog.Errorf("%s %s: %v", logPrefix, msg, err)
		c.JSON(http.StatusBadRequest, gin.H{"message": msg})
		return
	}
	glog.Infof("%s %s", logPrefix, model)

	// Write to database
	result, err := a.DB.Exec(CreateMeasurementSQL, model.MeasurementTypeId, model.Value, model.Repetitions, model.StartTime, model.Duration, model.DataSource)
	if err != nil {
		msg := "Couldn't add measurement to database"
		glog.Errorf("%s %s: %v", logPrefix, msg, err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": msg})
		return
	}
	// Get the Id that was just auto-written to the database
	id, err := result.LastInsertId()
	if err != nil {
		// If the db doesn't support LastInsertId(), throw an error for now
		msg := "Internal configuration mishap"
		glog.Errorf("%s %s: %v", logPrefix, msg, err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": msg})
		return
	}
	model.Id = id

	glog.Infof("%s Success: %+v", logPrefix, model)
	c.JSON(http.StatusCreated, gin.H{"measurement": model})
	return
}

func (a *Api) updateMeasurement(c *gin.Context) {
	logPrefix := GetLogPrefix(c)
	glog.Infof("%s fcn start", logPrefix)

	// Parse input
	paramId, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		msg := "URL id should be an integer"
		glog.Errorf("%s %s: %v", logPrefix, msg, err)
		c.JSON(http.StatusNotFound, gin.H{"message": msg})
		return
	}
	model := &Measurement{}
	err = c.Bind(model)
	if err != nil {
		msg := "Couldn't parse input form"
		glog.Errorf("%s %s: %v", logPrefix, msg, err)
		c.JSON(http.StatusBadRequest, gin.H{"message": msg})
		return
	}
	model.Id = paramId
	glog.Infof("%s %s", logPrefix, model)

	// Write to database
	_, err = a.DB.Exec(UpdateMeasurementSQL, model.MeasurementTypeId, model.Value, model.Repetitions, model.StartTime, model.Duration, model.DataSource, model.Id)
	// TODO(gary): add 404
	if err != nil {
		msg := "Couldn't update measurement in database"
		glog.Errorf("%s %s: %v", logPrefix, msg, err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": msg})
		return
	}

	glog.Infof("%s Success: %+v", logPrefix, model)
	c.JSON(http.StatusOK, gin.H{"measurement": model})
	return
}

func (a *Api) deleteMeasurement(c *gin.Context) {
	logPrefix := GetLogPrefix(c)
	glog.Infof("%s fcn start", logPrefix)

	// Parse input
	paramId, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		msg := "URL id should be an integer"
		glog.Errorf("%s %s: %v", logPrefix, msg, err)
		c.JSON(http.StatusBadRequest, gin.H{"message": msg})
		return
	}

	// Write to database
	_, err = a.DB.Exec(DeleteMeasurementSQL, paramId)
	// TODO(gary): add 404
	if err != nil {
		msg := "Couldn't delete measurement in database"
		glog.Errorf("%s %s: %v", logPrefix, msg, err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": msg})
		return
	}

	glog.Infof("%s Success", logPrefix)
	c.JSON(http.StatusNoContent, nil)
	return
}

func (a *Api) getMeasurement(c *gin.Context) {
	logPrefix := GetLogPrefix(c)
	glog.Infof("%s fcn start", logPrefix)

	// Parse input
	paramId, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		msg := "URL id should be an integer"
		glog.Errorf("%s %s: %v", logPrefix, msg, err)
		c.JSON(http.StatusNotFound, gin.H{"message": msg})
		return
	}

	// Read from database
	model := Measurement{}
	err = a.DB.QueryRow(GetMeasurementSQL, paramId).Scan(&model.Id, &model.MeasurementTypeId, &model.Value, &model.Repetitions, &model.StartTime, &model.Duration, &model.DataSource)
	if err == sql.ErrNoRows {
		msg := "Couldn't find a measurement with that Id"
		glog.Errorf("%s %s: %v", logPrefix, msg, err)
		c.JSON(http.StatusNotFound, gin.H{"message": msg})
		return
	} else if err != nil {
		msg := "Couldn't get measurement from database"
		glog.Errorf("%s %s: %v", logPrefix, msg, err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": msg})
		return
	}

	glog.Infof("%s Success: %+v", logPrefix, model)
	c.JSON(http.StatusOK, gin.H{"measurement": model})
	return
}

func (a *Api) listMeasurement(c *gin.Context) {
	logPrefix := GetLogPrefix(c)
	glog.Infof("%s fcn start", logPrefix)

	// Read from database
	rows, err := a.DB.Query(ListMeasurementSQL)
	if err != nil {
		msg := "Couldn't get measurements from database"
		glog.Errorf("%s %s: %v", logPrefix, msg, err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": msg})
		return
	}
	defer rows.Close()
	models := []Measurement{}
	for rows.Next() {
		model := Measurement{}
		err = rows.Scan(&model.Id, &model.MeasurementTypeId, &model.Value, &model.Repetitions, &model.StartTime, &model.Duration, &model.DataSource)
		if err != nil {
			msg := "Couldn't scan row from database"
			glog.Infof("%s %s: %v", logPrefix, msg, err)
		}
		glog.Infof("%s %s", logPrefix, model)
		models = append(models, model)
	}
	err = rows.Err()
	if err != nil {
		msg := "Error scanning rows from database"
		glog.Errorf("%s %s: %v", logPrefix, msg, err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": msg})
		return
	}

	glog.Infof("%s Success: %+v", logPrefix, models)
	c.JSON(http.StatusOK, gin.H{"measurements": models})
	return
}
