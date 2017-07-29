package api

import (
	"database/sql"
	"fmt"
	"github.com/golang/glog"
	"net/http"
	"strconv"
	"strings"

	"gopkg.in/gin-gonic/gin.v1"
)

func (a *Api) createMeasurementType(c *gin.Context) {
	logPrefix := GetLogPrefix(c)
	glog.Infof("%s fcn start", logPrefix)

	// Parse input
	model := &MeasurementType{}
	err := c.Bind(model)
	if err != nil {
		msg := "Couldn't parse input form"
		glog.Errorf("%s %s: %v", logPrefix, msg, err)
		c.JSON(http.StatusBadRequest, gin.H{"message": msg})
		return
	}
	glog.Infof("%s Name: %s, Units: %s", logPrefix, model.Name, model.Units)

	// Write to database
	result, err := a.DB.Exec(CreateMeasurementTypeSQL, model.Name, model.Units)
	if err != nil {
		if strings.Contains(err.Error(), "UNIQUE constraint failed") {
			msg := fmt.Sprintf("The Name '%s' already exists", model.Name)
			glog.Errorf("%s %s: %v", logPrefix, msg, err)
			c.JSON(http.StatusBadRequest, gin.H{"message": msg})
			return
		}
		msg := "Couldn't add measurement_type to database"
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
	c.JSON(http.StatusCreated, gin.H{"measurement_type": model})
	return
}

func (a *Api) updateMeasurementType(c *gin.Context) {
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
	model := &MeasurementType{}
	err = c.Bind(model)
	if err != nil {
		msg := "Couldn't parse input form"
		glog.Errorf("%s %s: %v", logPrefix, msg, err)
		c.JSON(http.StatusBadRequest, gin.H{"message": msg})
		return
	}
	// model.Id will be 0 if it wasn't set in the request
	// TODO(gary): Is there a way to check if this was set?
	if model.Id == 0 {
		model.Id = paramId
	}
	if model.Id != paramId {
		msg := fmt.Sprintf("URL Id (%s) and JSON Id (%d) mismatch", paramId, model.Id)
		glog.Errorf("%s %s", logPrefix, msg)
		c.JSON(http.StatusBadRequest, gin.H{"message": msg})
		return
	}
	glog.Infof("%s Name: %s, Units: %s", logPrefix, model.Name, model.Units)

	// Write to database
	_, err = a.DB.Exec(UpdateMeasurementTypeSQL, model.Name, model.Units, paramId)
	// TODO(gary): add 404
	if err != nil {
		msg := "Couldn't update measurement_type in database"
		glog.Errorf("%s %s: %v", logPrefix, msg, err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": msg})
		return
	}

	glog.Infof("%s Success: %+v", logPrefix, model)
	c.JSON(http.StatusOK, gin.H{"measurement_type": model})
	return
}

func (a *Api) deleteMeasurementType(c *gin.Context) {
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
	_, err = a.DB.Exec(DeleteMeasurementTypeSQL, paramId)
	// TODO(gary): add 404
	if err != nil {
		msg := "Couldn't delete measurement_type in database"
		glog.Errorf("%s %s: %v", logPrefix, msg, err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": msg})
		return
	}

	glog.Infof("%s Success", logPrefix)
	c.JSON(http.StatusNoContent, nil)
	return
}

func (a *Api) getMeasurementType(c *gin.Context) {
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
	model := MeasurementType{}
	err = a.DB.QueryRow(GetMeasurementTypeSQL, paramId).Scan(&model.Id, &model.Name, &model.Units)
	if err == sql.ErrNoRows {
		msg := "Couldn't find a measurement_type with that Id"
		glog.Errorf("%s %s: %v", logPrefix, msg, err)
		c.JSON(http.StatusNotFound, gin.H{"message": msg})
		return
	} else if err != nil {
		msg := "Couldn't get measurement_type from database"
		glog.Errorf("%s %s: %v", logPrefix, msg, err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": msg})
		return
	}

	glog.Infof("%s Success: %+v", logPrefix, model)
	c.JSON(http.StatusOK, gin.H{"measurement_type": model})
	return
}

func (a *Api) listMeasurementType(c *gin.Context) {
	logPrefix := GetLogPrefix(c)
	glog.Infof("%s fcn start", logPrefix)

	// Read from database
	rows, err := a.DB.Query(ListMeasurementTypeSQL)
	if err != nil {
		msg := "Couldn't get measurement_types from database"
		glog.Errorf("%s %s: %v", logPrefix, msg, err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": msg})
		return
	}
	defer rows.Close()
	models := []MeasurementType{}
	for rows.Next() {
		model := MeasurementType{}
		err = rows.Scan(&model.Id, &model.Name, &model.Units)
		if err != nil {
			msg := "Couldn't scan row from database"
			glog.Infof("%s %s: %v", logPrefix, msg, err)
		}
		glog.Infof("%s Name: %s, Units: %s", logPrefix, model.Name, model.Units)
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
	c.JSON(http.StatusOK, gin.H{"measurement_types": models})
	return
}
