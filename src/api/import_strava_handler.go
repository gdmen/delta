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
	// Buffer measurements to make a single write to the db at the end
	measurements := []*Measurement{}
	// Store measurement type name:id we've already created so we don't need to make db requests to check for pre-existence
	measurementTypesCreated := map[string]int64{}
	// The start and end times of the full range of measurements being imported in this call
	var importStart, importEnd int64
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
			if mtId, ok := measurementTypesCreated[mt.Name]; ok {
				mt.Id = mtId
			} else {
				status, msg, err := mtManager.Create(mt)
				// TODO: create fcn should check for exact existence and return 201 if so. 400 error indicates that the Unique constraint was violated. If so, still add the Measurement.
				if err != nil && status != http.StatusBadRequest {
					glog.Errorf("%s %s: %v", logPrefix, msg, err)
				}
				measurementTypesCreated[mt.Name] = mt.Id
			}
			startTimeObj, _ := trk.TimeBounds()
			startTime := startTimeObj.Unix()
			m := &Measurement{
				MeasurementTypeId: mt.Id,
				Value:             trk.Length3D() / 1609.34, // Length3D is in meters, convert it to miles
				Repetitions:       0,
				StartTime:         startTime,
				Duration:          int32(trk.Duration()),
				DataSource:        "strava",
			}
			measurements = append(measurements, m)
			// Update the total time range
			if importStart == 0 || startTime < importStart {
				importStart = startTime
			}
			if startTime > importEnd {
				importEnd = startTime
			}
		}
	}

	// Delete all entries from the database that are from this source & during the time range covered by the current import
	// This avoids duplicating measurements
	glog.Infof("Clearing time range (%d - %d) for %s from database", importStart, importEnd, "strava")
	status, msg, err := mManager.DeleteTimeRangeForSource(importStart, importEnd, "strava")
	if err != nil {
		glog.Errorf("%s %s: %v", logPrefix, msg, err)
	} else {
		// Save all models we just parsed
		status, msg, err = mManager.CreateMultiple(measurements)
		if err != nil {
			glog.Errorf("%s %s: %v", logPrefix, msg, err)
		}
	}

	c.JSON(status, gin.H{})
	return
}
