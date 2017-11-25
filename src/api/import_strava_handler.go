package api

import (
	"fmt"
	"io"
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

	c.JSON(a.importMeasurements(c, parseStravaFile, "strava"))
	return
}

func parseStravaFile(c *gin.Context, file io.Reader, filename string, mtManager *MeasurementTypeManager, measurementTypesCreated map[string]int64) (*[]*Measurement, error) {
	logPrefix := common.GetLogPrefix(c)
	glog.Infof("%s fcn start", logPrefix)

	measurements := []*Measurement{}
	g, err := gpx.Parse(file)
	if err != nil {
		msg := fmt.Sprintf("Error parsing GPX file: %s", filename)
		glog.Errorf("%s %s: %v+", logPrefix, msg, err)
		return nil, err
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
	}
	return &measurements, nil
}
