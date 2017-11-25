package api

import (
	"encoding/csv"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang/glog"

	"github.com/gdmen/delta/src/common"
)

var fitnotes_name_override_map = map[string]string{
	"Barbell Squat":                       "Barbell Back Squat",
	"BJJ":                                 "Brazilian Jiu-Jitsu",
	"Deadlift":                            "Conventional Barbell Deadlift",
	"Dumbbell Overhead Triceps Extension": "Lying Dumbbell Triceps Extension",
	"Lying Triceps Extension":             "Lying Barbell Triceps Extension",
	"Overhead Press":                      "Standing Barbell Shoulder Press (OHP)",
	"Running (Outdoor)":                   "Running",
	"Stationary Bike":                     "Road Cycling",
	"Cycling":                             "Road Cycling",
}

var fitnotesUnitRegex = regexp.MustCompile(`[^()]+\((?P<Units>.*)\)`)

func fitnotesGetUniformName(name string) string {
	if val, ok := fitnotes_name_override_map[name]; ok {
		return val
	}
	return name
}

func (a *Api) importFitnotes(c *gin.Context) {
	logPrefix := common.GetLogPrefix(c)
	glog.Infof("%s fcn start", logPrefix)

	c.JSON(a.importMeasurements(c, parseFitnotesFile, "fitnotes"))
	return
}

func parseFitnotesFile(c *gin.Context, file io.Reader, filename string, mtManager *MeasurementTypeManager, measurementTypesCreated map[string]int64) (*[]*Measurement, error) {
	logPrefix := common.GetLogPrefix(c)
	glog.Infof("%s fcn start", logPrefix)

	measurements := []*Measurement{}
	r := csv.NewReader(file)
	isHeaderRow := true
	var weight_units string
	for {
		// ["Date", "Exercise", "Category", "Weight (lbs)", "Reps", "Distance", "Distance Unit", "Time"]
		row, err := r.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			msg := fmt.Sprintf("Error reading file: %s", filename)
			glog.Errorf("%s %s: %+v", logPrefix, msg, err)
			return nil, err
		}
		glog.Infof("%s read row: %s", logPrefix, row)
		if isHeaderRow {
			weight_units = fitnotesUnitRegex.FindStringSubmatch(row[3])[1]
			isHeaderRow = false
			continue
		}
		// Process row
		date := row[0]
		name := row[1]
		weight := row[3]
		reps := row[4]
		distance := row[5]
		distance_units := row[6]

		// TODO: don't use hardcoded time zones. . .
		loc, _ := time.LoadLocation("America/Los_Angeles")
		startDate, err := time.ParseInLocation("2006-01-02", date, loc)
		if err != nil {
			msg := fmt.Sprintf("Error parsing file: %s", filename)
			glog.Errorf("%s %s: %+v", logPrefix, msg, err)
			return nil, err
		}
		startTime := startDate.Unix()

		var repetitions int64
		if reps != "" {
			repetitions, err = strconv.ParseInt(reps, 10, 16)
			if err != nil {
				msg := fmt.Sprintf("Error parsing file: %s", filename)
				glog.Errorf("%s %s: %+v", logPrefix, msg, err)
				return nil, err
			}
		}

		var float64Distance float64
		if distance != "" {
			float64Distance, err = strconv.ParseFloat(distance, 64)
			if err != nil {
				msg := fmt.Sprintf("Error parsing file: %s", filename)
				glog.Errorf("%s %s: %+v", logPrefix, msg, err)
				return nil, err
			}
		}
		var value float64
		var units string
		if distance != "" && float64Distance > 0 {
			value = float64Distance
			units = distance_units
		} else if weight != "" {
			var float64Weight float64
			if weight != "" {
				float64Weight, err = strconv.ParseFloat(weight, 64)
				if err != nil {
					msg := fmt.Sprintf("Error parsing file: %s", filename)
					glog.Errorf("%s %s: %+v", logPrefix, msg, err)
					return nil, err
				}
			}
			value = float64Weight
			units = weight_units
		}

		// Convert '1:00:00', '1:00', or '' to seconds
		split := strings.Split(row[7], ":")
		duration := 0
		multiplier := 1
		for i := len(split) - 1; i >= 0; i-- {
			var addtlTime int
			if split[i] != "" {
				addtlTime, err = strconv.Atoi(split[i])
				if err != nil {
					msg := fmt.Sprintf("Error parsing file: %s", filename)
					glog.Errorf("%s %s: %+v", logPrefix, msg, err)
					return nil, err
				}
			}
			duration += addtlTime * multiplier
			multiplier *= 60
		}

		// Save models
		mt := &MeasurementType{
			Name:  fitnotesGetUniformName(name),
			Units: units,
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
		m := &Measurement{
			MeasurementTypeId: mt.Id,
			Value:             value,
			Repetitions:       int16(repetitions),
			StartTime:         startTime,
			Duration:          int32(duration),
			DataSource:        "fitnotes",
		}
		measurements = append(measurements, m)
	}
	return &measurements, nil
}
