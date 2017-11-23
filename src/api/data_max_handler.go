package api

import (
	"encoding/json"
	"math"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang/glog"

	"github.com/gdmen/delta/src/common"
)

type MaxesDataField struct {
	Name string `json:"name" form:"name" binding:"required"`
}

type MaxesDataForm struct {
	Fields    string    `json:"fields" form:"fields" binding:"required"`
	MaxOnly   bool      `json:"maxOnly" form:"maxOnly"`
	Increment increment `json:"increment" form:"increment" binding:"required"`
}

type MaxesHCData struct {
	Name string  `json:"name"`
	Y    float64 `json:"y"`
}

type DecayDetails struct {
	Y    float64
	Days int64 // # of days that have passed since this max
}

const MaxDecayDays = 365

func getDecayedMax(m *DecayDetails) float64 {
	decayed := m.Y * math.Log(math.Max(1, float64(MaxDecayDays-m.Days))) / math.Log(MaxDecayDays)
	return decayed
}

func (a *Api) getMaxes(c *gin.Context) {
	logPrefix := common.GetLogPrefix(c)
	glog.Infof("%s fcn start", logPrefix)

	// Parse input
	form := &MaxesDataForm{}
	err := c.Bind(form)
	var dataFields []MaxesDataField
	if err == nil {
		err = json.Unmarshal([]byte(form.Fields), &dataFields)
	}
	if err != nil {
		msg := "Couldn't parse input"
		glog.Errorf("%s %s: %v", logPrefix, msg, err)
		c.JSON(http.StatusBadRequest, gin.H{"message": msg})
		return
	}
	glog.Infof("%s form: %+v", logPrefix, form)
	glog.Infof("%s dataFields: %+v", logPrefix, dataFields)

	if form.Increment < 0 || form.Increment > MONTH {
		msg := "Couldn't parse input, unknown time increment"
		glog.Errorf("%s %s: %d", logPrefix, msg, form.Increment)
		c.JSON(http.StatusBadRequest, gin.H{"message": msg})
		return
	}

	// Read MeasurementTypes from database
	mtManager := &MeasurementTypeManager{DB: a.DB}
	measurementTypes, status, msg, err := mtManager.List()
	if err != nil {
		glog.Errorf("%s %s: %v", logPrefix, msg, err)
		c.JSON(status, gin.H{"message": msg})
		return
	}
	glog.Infof("%s MeasurementTypes: %+v", logPrefix, measurementTypes)
	measurementTypeNameToId := map[string]int64{}
	for _, mt := range *measurementTypes {
		measurementTypeNameToId[mt.Name] = mt.Id
	}

	// Read filtered Measurements from database
	measurementTypeIds := []string{}
	for _, field := range dataFields {
		if id, ok := measurementTypeNameToId[field.Name]; ok {
			measurementTypeIds = append(measurementTypeIds, strconv.FormatInt(id, 10))
		} else {
			msg := "Couldn't parse input, unknown measurement name: " + field.Name
			glog.Errorf("%s %s", logPrefix, msg)
			c.JSON(http.StatusBadRequest, gin.H{"message": msg})
			return
		}
	}
	mManager := &MeasurementManager{DB: a.DB}
	sql := "select * from measurements where measurement_type_id in (" + strings.Join(measurementTypeIds, ",") + ") order by start_time ASC;"
	measurements, status, msg, err := mManager.Custom(sql)
	if err != nil {
		glog.Errorf("%s %s: %v", logPrefix, msg, err)
		c.JSON(status, gin.H{"message": msg})
		return
	}
	glog.Infof("%s %d measurements", logPrefix, len(*measurements))

	// Generate output x values (with 0 values, a.k.a. not sparse)
	// First level of the map is <x> (timestamp)
	// Second level of the map is measurement_type_id
	data := map[int64]map[int64]float64{}
	var minTs int64
	var currTs int64
	// TODO: don't use hardcoded time zones. . .
	loc, _ := time.LoadLocation("America/Los_Angeles")
	for _, m := range *measurements {
		t := time.Unix(m.StartTime, 0).In(loc)
		// round to day
		tDay := time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
		switch form.Increment {
		case DAY:
			t = tDay
		case WEEK:
			// round to day and then iterate down to Sunday
			t = time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
			for t.Weekday() != time.Sunday {
				t = t.AddDate(0, 0, -1)
			}
		case MONTH:
			// get first of month
			t = time.Date(t.Year(), t.Month(), 1, 0, 0, 0, 0, t.Location())
		}
		newTs := t.Unix()
		if newTs != currTs {
			if currTs == 0 {
				// if this is the first newTs we've seen, initialize minTs
				minTs = newTs
			}
			if newTs < minTs {
				minTs = newTs
			}
			if _, ok := data[newTs]; !ok {
				data[newTs] = map[int64]float64{}
			}
			data[newTs][m.MeasurementTypeId] = 0
			currTs = newTs
		}
		// Potential new max for currTs
		computedMax := m.Value
		if m.Repetitions > 1 {
			computedMax = math.Floor(m.Value * (1 + (float64(m.Repetitions) / 30)))
		}
		data[currTs][m.MeasurementTypeId] = math.Max(data[currTs][m.MeasurementTypeId], computedMax)
	}
	/*
	   type MaxesHCData struct {
	   	Name      string  `json:"name"`
	   	Y         float64 `json:"y"`
	   }
	*/
	var hcData []MaxesHCData
	// map of measurementType id to the most significant recent max
	var mTypeDecay = map[int64]*DecayDetails{}
	currTime := time.Unix(minTs, 0).In(loc)
	for currTime.Before(time.Now().In(loc)) {
		// get y value
		var y float64
		var mTypeVals map[int64]float64
		mTypeVals, ok := data[currTime.Unix()]
		if !ok {
			mTypeVals = map[int64]float64{}
		}
		for _, mTypeIdStr := range measurementTypeIds {
			mTypeId, _ := strconv.ParseInt(mTypeIdStr, 10, 64)
			if _, ok := mTypeVals[mTypeId]; !ok {
				// if not okay, we first add a 0 max for this measurement
				mTypeVals[mTypeId] = 0
				// try to use the stored decaying max
				if decayDetails, ok := mTypeDecay[mTypeId]; ok {
					mTypeVals[mTypeId] = math.Max(mTypeVals[mTypeId], getDecayedMax(decayDetails))
				}
			}
			mTypeDecay[mTypeId] = &DecayDetails{
				Y:    mTypeVals[mTypeId],
				Days: 0,
			}
			y += mTypeVals[mTypeId]
		}

		// get x value (modify from first to last day of the period)
		var x string
		var nextTime time.Time
		var decayDaysIncr int64
		switch form.Increment {
		case DAY:
			x = currTime.Format("2006-01-02")
			nextTime = currTime.AddDate(0, 0, 1)
			decayDaysIncr = 1
		case WEEK:
			x = currTime.AddDate(0, 0, 6).Format("2006-01-02")
			nextTime = currTime.AddDate(0, 0, 7)
			decayDaysIncr = 7
		case MONTH:
			x = currTime.AddDate(0, 1, -1).Format("Jan '06")
			nextTime = currTime.AddDate(0, 1, 0)
			decayDaysIncr = int64(nextTime.Sub(currTime).Hours() / 24)
		}
		hcData = append(hcData, MaxesHCData{Name: x, Y: y})

		// increment currTime
		currTime = nextTime
		// increment decay days
		for _, mTypeIdStr := range measurementTypeIds {
			mTypeId, _ := strconv.ParseInt(mTypeIdStr, 10, 64)
			mTypeDecay[mTypeId].Days += decayDaysIncr
		}
	}

	if form.MaxOnly {
		c.JSON(http.StatusOK, gin.H{"value": hcData[len(hcData)-1].Y})
	} else {
		c.JSON(http.StatusOK, gin.H{"data": hcData})
	}
	return
}
