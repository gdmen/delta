package api

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang/glog"
	gpx "github.com/ptrv/go-gpx"

	"github.com/gdmen/delta/src/common"
)

var strava_name_override_map = map[string]string{
	"Run":  "Running",
	"Ride": "Road Cycling",
}

func stravaGetUniformName(name string) string {
	// Strava names are e.g. "Morning Run", "Evening Run"
	split_name := strings.Fields(name)
	name = split_name[len(split_name)-1]
	if val, ok := strava_name_override_map[name]; ok {
		return val
	}
	return name
}

func (a *Api) importStrava(c *gin.Context) {
	logPrefix := common.GetLogPrefix(c)
	glog.Infof("%s fcn start", logPrefix)

	// Parse input
	form, err := c.MultipartForm()
	if err != nil {
		msg := "Couldn't parse input form"
		glog.Errorf("%s %s: %v", logPrefix, msg, err)
		c.JSON(http.StatusBadRequest, gin.H{"message": msg})
		return
	}
	glog.Infof("%s form: %v+", logPrefix, form)
	fileHeaders := form.File["files"]

	// Read file and buffer models
	mtManager := &MeasurementTypeManager{DB: a.DB}
	mManager := &MeasurementManager{DB: a.DB}
	for _, fileHeader := range fileHeaders {
		glog.Infof("%s parsing file: %s", logPrefix, fileHeader.Filename)
		file, err := fileHeader.Open()
		if err != nil {
			msg := fmt.Sprintf("Error reading file: %s", fileHeader.Filename)
			glog.Errorf("%s %s: %v+", logPrefix, msg, err)
		}
		// Read as GPX
		g, err := gpx.Parse(file)
		if err != nil {
			msg := fmt.Sprintf("Error parsing GPX file: %s", fileHeader.Filename)
			glog.Errorf("%s %s: %v+", logPrefix, msg, err)
		}

		for _, trk := range g.Tracks {
			// Save models
			mt := &MeasurementType{
				Name:  stravaGetUniformName(trk.Name),
				Units: "mi",
			}
			status, msg, err := mtManager.Create(mt)
			// TODO: create fcn should check for exact existence and return 201 if so. 400 error indicates that the Unique constraint was violated. If so, still add the Measurement.
			if err != nil && status != http.StatusBadRequest {
				glog.Errorf("%s %s: %v", logPrefix, msg, err)
				continue
			}
			startTime, _ := trk.TimeBounds()
			m := &Measurement{
				MeasurementTypeId: mt.Id,
				Value:             trk.Length3D() / 1609.34, // Length3D is in meters, convert it to miles
				Repetitions:       0,
				StartTime:         startTime.Unix(),
				Duration:          int32(trk.Duration()),
				DataSource:        "strava",
			}
			status, msg, err = mManager.Create(m)
			if err != nil {
				glog.Errorf("%s %s: %v", logPrefix, msg, err)
				continue
			}
		}
	}

	c.JSON(http.StatusCreated, gin.H{})
	return
}
