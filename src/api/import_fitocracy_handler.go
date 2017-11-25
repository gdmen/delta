package api

import (
	"encoding/csv"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang/glog"

	"github.com/gdmen/delta/src/common"
)

var fitocracy_name_override_map = map[string]string{
	"Ab Wheel (kneeling)":                                     "Kneeling Ab Wheel",
	"Barbell Bench Press":                                     "Flat Barbell Bench Press",
	"Barbell Deadlift":                                        "Conventional Barbell Deadlift",
	"Barbell Squat":                                           "Barbell Back Squat",
	"Chin-Up":                                                 "Chin Up",
	"Cycling":                                                 "Road Cycling",
	"Dips - Triceps Version":                                  "Parallel Bar Triceps Dip",
	"General Yoga":                                            "Yoga",
	"Indoor Volleyball":                                       "Volleyball",
	"Light Walking (secondary e.g. commute, on the job, etc)": "Walking",
	"One-Arm Dumbbell Row":                                    "One Arm Dumbbell Row",
	"Parallel-Grip Pull-Up":                                   "Neutral Grip Pull Up",
	"Pull-Up":                                                 "Pull Up",
	"Push-Up":                                                 "Push Up",
	"Standing Military Press":                                 "Standing Barbell Shoulder Press (OHP)",
	"Wide-Grip Pull-Up":                                       "Wide Grip Pull Up",
}

func fitocracyGetUniformName(name string) string {
	if val, ok := fitocracy_name_override_map[name]; ok {
		return val
	}
	return name
}

func (a *Api) importFitocracy(c *gin.Context) {
	logPrefix := common.GetLogPrefix(c)
	glog.Infof("%s fcn start", logPrefix)

	c.JSON(a.importMeasurements(c, parseFitocracyFile, "fitocracy"))
	return
}

func parseFitocracyFile(c *gin.Context, file io.Reader, filename string, mtManager *MeasurementTypeManager, measurementTypesCreated map[string]int64) (*[]*Measurement, error) {
	logPrefix := common.GetLogPrefix(c)
	glog.Infof("%s fcn start", logPrefix)

	measurements := []*Measurement{}
	r := csv.NewReader(file)
	isHeaderRow := true
	for {
		// ["Activity", "Date (YYYYMMDD)", "Set", "", "unit", "Combined", "Points"]
		// or
		// ["Activity","Date (YYYYMMDD)","Session",,"unit","Combined","Points"]
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
			isHeaderRow = false
			continue
		}
		// Process row
		name := row[0]
		date := row[1]

		var repetitions int16
		var duration int32
		var value float64
		var units string

		// Parse combined set/rep/etc field
		// e.g.
		//  5 reps
		//  5 reps || weighted || 10 lb
		//  60 min || 3.2 mi
		//  45 lb || 5 reps
		//  45 min
		//  120 min || practice
		combined := row[len(row)-2]
		pieces := strings.Split(combined, "||")

		for _, piece := range pieces {
			// e.g.
			// 60 min
			// weighted
			piece = strings.TrimSpace(piece)
			split_piece := strings.Split(piece, " ")
			if len(split_piece) < 2 {
				// ignore things like "weighted" and "practice"
				glog.Infof("%s Skipping piece %s", logPrefix, piece)
				continue
			}
			_value, err := strconv.ParseFloat(split_piece[0], 64)
			if err != nil {
				glog.Infof("%s Couldn't parse value from piece %s: %+v", logPrefix, piece, err)
				continue
			}
			_units := split_piece[1]
			switch _units {
			case "lb":
				value = _value
				units = "lbs"
			case "reps":
				repetitions = int16(_value)
			case "mi":
				value = _value
				units = "mi"
			case "min":
				duration = int32(_value * 60)
			default:
				glog.Infof("%s Couldn't parse piece %s", logPrefix, piece)
			}
		}

		// TODO: don't use hardcoded time zones. . .
		loc, _ := time.LoadLocation("America/Los_Angeles")
		startDate, err := time.ParseInLocation("2006-01-02", date, loc)
		if err != nil {
			msg := fmt.Sprintf("Error parsing file: %s", filename)
			glog.Errorf("%s %s: %+v", logPrefix, msg, err)
			return nil, err
		}
		startTime := startDate.Unix()

		// Save models
		mt := &MeasurementType{
			Name:  fitocracyGetUniformName(name),
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
			Repetitions:       repetitions,
			StartTime:         startTime,
			Duration:          duration,
			DataSource:        "fitocracy",
		}
		measurements = append(measurements, m)
	}
	return &measurements, nil
}
