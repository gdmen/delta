package api

import (
	"fmt"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/golang/glog"

	"github.com/gdmen/delta/src/common"
)

// Parse an input file and return Measurements
// As a side effect, it's okay to create MeasurementTypes
// It is also expected that *measurementTypesCreated is modified in place
type ParseImportFile func(*gin.Context, io.Reader, string, *MeasurementTypeManager, map[string]int64) (*[]*Measurement, error)

func (a *Api) importMeasurements(c *gin.Context, parseFcn ParseImportFile, dataSource string) (int, gin.H) {
	logPrefix := common.GetLogPrefix(c)
	glog.Infof("%s fcn start", logPrefix)

	// Parse input
	form, err := c.MultipartForm()
	if err != nil {
		msg := "Couldn't parse input form"
		glog.Errorf("%s %s: %v", logPrefix, msg, err)
		return http.StatusBadRequest, gin.H{"message": msg}
	}
	glog.Infof("%s form: %+v", logPrefix, form)
	fileHeaders := form.File["files"]

	// Read file and buffer models
	mtManager := &MeasurementTypeManager{DB: a.DB}
	mManager := &MeasurementManager{DB: a.DB}
	// Buffer measurements to make a single write to the db at the end
	measurements := []*Measurement{}
	// Store measurement type name:id we've already created so we don't need to make db requests to check for pre-existence
	measurementTypesCreated := map[string]int64{}
	for _, fileHeader := range fileHeaders {
		glog.Infof("%s parsing file: %s", logPrefix, fileHeader.Filename)
		file, err := fileHeader.Open()
		if err != nil {
			msg := fmt.Sprintf("Error reading file: %s", fileHeader.Filename)
			glog.Errorf("%s %s: %+v", logPrefix, msg, err)
		}
		// Parse

		parsedMeasurements, err := parseFcn(c, file, fileHeader.Filename, mtManager, measurementTypesCreated)
		if err != nil {
			msg := fmt.Sprintf("Failed parsing file: %s", fileHeader.Filename)
			glog.Errorf("%s %s: %+v", logPrefix, msg, err)
		}
		measurements = append(measurements, *parsedMeasurements...)
	}

	// The start and end times of the full range of measurements being imported in this call
	var importStart, importEnd int64
	for _, m := range measurements {
		if importStart == 0 || m.StartTime < importStart {
			importStart = m.StartTime
		}
		if m.StartTime > importEnd {
			importEnd = m.StartTime
		}
	}

	// Delete all entries from the database that are from this source & during the time range covered by the current import
	// This avoids duplicating measurements
	glog.Infof("Clearing time range (%d - %d) for %s from database", importStart, importEnd, dataSource)
	status, msg, err := mManager.DeleteTimeRangeForSource(importStart, importEnd, dataSource)
	if err != nil {
		glog.Errorf("%s %s: %v", logPrefix, msg, err)
	} else {
		// Save all models we just parsed
		status, msg, err = mManager.CreateMultiple(measurements)
		if err != nil {
			glog.Errorf("%s %s: %v", logPrefix, msg, err)
		}
	}

	return status, gin.H{}
}
