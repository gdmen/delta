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
	glog.Infof("%s form: %v+", logPrefix, form)
	fileHeaders := form.File["files"]

	// Read file and buffer models
	mtManager := &MeasurementTypeManager{DB: a.DB}
	mManager := &MeasurementManager{DB: a.DB}
	unitsRegex := regexp.MustCompile(`[^()]+\((?P<Units>.*)\)`)
	for _, fileHeader := range fileHeaders {
		glog.Infof("%s parsing file: %s", logPrefix, fileHeader.Filename)
		file, err := fileHeader.Open()
		if err != nil {
			msg := fmt.Sprintf("Error reading file: %s", fileHeader.Filename)
			glog.Errorf("%s %s: %v+", logPrefix, msg, err)
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
				glog.Errorf("%s %s: %v+", logPrefix, msg, err)
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

			//date = strings.Replace(date, "-", "/", -1)
			startDate, err := time.Parse("2006-01-02", date)
			if err != nil {
				msg := fmt.Sprintf("Error parsing file: %s", fileHeader.Filename)
				glog.Errorf("%s %s: %v+", logPrefix, msg, err)
				continue
			}
			startTime := startDate.Unix()

			var repetitions int64
			if reps != "" {
				repetitions, err = strconv.ParseInt(reps, 10, 16)
				if err != nil {
					msg := fmt.Sprintf("Error parsing file: %s", fileHeader.Filename)
					glog.Errorf("%s %s: %v+", logPrefix, msg, err)
					continue
				}
			}

			var float64Distance float64
			if distance != "" {
				float64Distance, err = strconv.ParseFloat(distance, 64)
				if err != nil {
					msg := fmt.Sprintf("Error parsing file: %s", fileHeader.Filename)
					glog.Errorf("%s %s: %v+", logPrefix, msg, err)
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
						glog.Errorf("%s %s: %v+", logPrefix, msg, err)
						continue
					}
				}
				value = float64Weight
				units = weight_units
			}

			// Convert '1:00:00', '1:00', or '' to seconds
			glog.Errorf("%s %s: %v", logPrefix, "time", row[7])
			split := strings.Split(row[7], ":")
			glog.Errorf("%s %s: %v", logPrefix, "split", split)
			duration := 0
			multiplier := 1
			for i := len(split) - 1; i >= 0; i-- {
				var addtlTime int
				if split[i] != "" {
					addtlTime, err = strconv.Atoi(split[i])
					if err != nil {
						msg := fmt.Sprintf("Error parsing file: %s", fileHeader.Filename)
						glog.Errorf("%s %s: %v+", logPrefix, msg, err)
						continue
					}
				}
				duration += addtlTime * multiplier
				multiplier *= 60
				glog.Errorf("%s %s: %v", logPrefix, "duration", duration)
			}

			// Save models
			mt := &MeasurementType{
				Name:  fitnotesGetUniformName(name),
				Units: units,
			}
			status, msg, err := mtManager.Create(mt)
			// TODO: create fcn should check for exact existence and return 201 if so. 400 error indicates that the Unique constraint was violated. If so, still add the Measurement.
			if err != nil && status != http.StatusBadRequest {
				glog.Errorf("%s %s: %v", logPrefix, msg, err)
				continue
			}
			m := &Measurement{
				MeasurementTypeId: mt.Id,
				Value:             value,
				Repetitions:       int16(repetitions),
				StartTime:         startTime,
				Duration:          int32(duration),
				DataSource:        "fitnotes",
			}
			status, msg, err = mManager.Create(m)
			if err != nil {
				glog.Errorf("%s %s: %v", logPrefix, msg, err)
				continue
			}
		}
	}

	// Write models

	// Write to database
	/*manager := &MeasurementTypeManager{DB: a.DB}
	status, msg, err := manager.Create(model)
	if err != nil {
		glog.Errorf("%s %s: %v", logPrefix, msg, err)
		c.JSON(status, gin.H{"message": msg})
		return
	}
	manager := &MeasurementManager{DB: a.DB}
	status, msg, err := manager.Create(model)
	if err != nil {
		glog.Errorf("%s %s: %v", logPrefix, msg, err)
		c.JSON(status, gin.H{"message": msg})
		return
	}*/

	//glog.Infof("%s Success: %+v", logPrefix, len(measurements))
	c.JSON(http.StatusCreated, gin.H{})
	return
}
