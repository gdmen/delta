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

func fitnotesGetUniformName(name string) string {
	if val, ok := fitnotes_name_override_map[name]; ok {
		return val
	}
	return name
}

func (a *Api) importFitnotes(c *gin.Context) {
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
	glog.Infof("%s form: %+v", logPrefix, form)
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
	unitsRegex := regexp.MustCompile(`[^()]+\((?P<Units>.*)\)`)
	for _, fileHeader := range fileHeaders {
		glog.Infof("%s parsing file: %s", logPrefix, fileHeader.Filename)
		file, err := fileHeader.Open()
		if err != nil {
			msg := fmt.Sprintf("Error reading file: %s", fileHeader.Filename)
			glog.Errorf("%s %s: %+v", logPrefix, msg, err)
		}
		// Read as CSV
		r := csv.NewReader(file)
		isHeaderRow := true
		var weight_units string
		for {
			// ["Date", "Exercise", "Category", "Weight (lbs)", "Reps", "Distance", "Distance Unit", "Time"]
			row, err := r.Read()
			if err == io.EOF {
				break
			} else if err != nil {
				msg := fmt.Sprintf("Error reading file: %s", fileHeader.Filename)
				glog.Errorf("%s %s: %+v", logPrefix, msg, err)
			}
			glog.Infof("%s read row: %s", logPrefix, row)
			if isHeaderRow {
				weight_units = unitsRegex.FindStringSubmatch(row[3])[1]
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
				msg := fmt.Sprintf("Error parsing file: %s", fileHeader.Filename)
				glog.Errorf("%s %s: %+v", logPrefix, msg, err)
				continue
			}
			startTime := startDate.Unix()

			var repetitions int64
			if reps != "" {
				repetitions, err = strconv.ParseInt(reps, 10, 16)
				if err != nil {
					msg := fmt.Sprintf("Error parsing file: %s", fileHeader.Filename)
					glog.Errorf("%s %s: %+v", logPrefix, msg, err)
					continue
				}
			}

			var float64Distance float64
			if distance != "" {
				float64Distance, err = strconv.ParseFloat(distance, 64)
				if err != nil {
					msg := fmt.Sprintf("Error parsing file: %s", fileHeader.Filename)
					glog.Errorf("%s %s: %+v", logPrefix, msg, err)
					continue
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
						msg := fmt.Sprintf("Error parsing file: %s", fileHeader.Filename)
						glog.Errorf("%s %s: %+v", logPrefix, msg, err)
						continue
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
						msg := fmt.Sprintf("Error parsing file: %s", fileHeader.Filename)
						glog.Errorf("%s %s: %+v", logPrefix, msg, err)
						continue
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
	glog.Infof("Clearing time range (%d - %d) for %s from database", importStart, importEnd, "fitnotes")
	status, msg, err := mManager.DeleteTimeRangeForSource(importStart, importEnd, "fitnotes")
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
