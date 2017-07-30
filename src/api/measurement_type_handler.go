package api

import (
	"github.com/golang/glog"
	"net/http"
	"strconv"

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
	glog.Infof("%s %s", logPrefix, model)

	// Write to database
	manager := &MeasurementTypeManager{DB: a.DB}
	status, msg, err := manager.Create(model)
	if err != nil {
		glog.Errorf("%s %s: %v", logPrefix, msg, err)
		c.JSON(status, gin.H{"message": msg})
		return
	}

	glog.Infof("%s Success: %+v", logPrefix, model)
	c.JSON(status, gin.H{"measurement_type": model})
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
	model.Id = paramId
	glog.Infof("%s %s", logPrefix, model)

	// Write to database
	manager := &MeasurementTypeManager{DB: a.DB}
	status, msg, err := manager.Update(model)
	if err != nil {
		glog.Errorf("%s %s: %v", logPrefix, msg, err)
		c.JSON(status, gin.H{"message": msg})
		return
	}

	glog.Infof("%s Success: %+v", logPrefix, model)
	c.JSON(status, gin.H{"measurement_type": model})
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
		c.JSON(http.StatusNotFound, gin.H{"message": msg})
		return
	}

	// Write to database
	manager := &MeasurementTypeManager{DB: a.DB}
	status, msg, err := manager.Delete(paramId)
	if err != nil {
		glog.Errorf("%s %s: %v", logPrefix, msg, err)
		c.JSON(status, gin.H{"message": msg})
		return
	}

	glog.Infof("%s Success", logPrefix)
	c.JSON(status, nil)
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
	manager := &MeasurementTypeManager{DB: a.DB}
	model, status, msg, err := manager.Get(paramId)
	if err != nil {
		glog.Errorf("%s %s: %v", logPrefix, msg, err)
		c.JSON(status, gin.H{"message": msg})
		return
	}

	glog.Infof("%s Success: %+v", logPrefix, model)
	c.JSON(status, gin.H{"measurement_type": model})
	return
}

func (a *Api) listMeasurementType(c *gin.Context) {
	logPrefix := GetLogPrefix(c)
	glog.Infof("%s fcn start", logPrefix)

	// Read from database
	manager := &MeasurementTypeManager{DB: a.DB}
	models, status, msg, err := manager.List()
	if err != nil {
		glog.Errorf("%s %s: %v", logPrefix, msg, err)
		c.JSON(status, gin.H{"message": msg})
		return
	}

	glog.Infof("%s Success: %+v", logPrefix, models)
	c.JSON(status, gin.H{"measurement_types": models})
	return
}
