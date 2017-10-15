package api

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang/glog"

	"github.com/gdmen/delta/src/common"
)

type attr int

const (
	_          = iota
	VALUE attr = iota
	DURATION
)

type DataField struct {
	Name string `json:"name" form:"name" binding:"required"`
	Attr attr   `json:"attr" form:"attr" binding:"required"`
}

type increment int

const (
	_             = iota
	DAY increment = iota
	WEEK
	MONTH
)

type DataForm struct {
	Fields    string    `json:"fields" form:"fields" binding:"required"`
	Increment increment `json:"fields" form:"increment" binding:"required"`
}

type HCData struct {
	Name      string  `json:"name"`
	Drilldown string  `json:"drilldown"`
	Y         float64 `json:"y"`
}
type HCDrilldownData struct {
	Name string          `json:"name"`
	Id   string          `json:"id"`
	Data [][]interface{} `json:"data"`
}

func (a *Api) getDrilldown(c *gin.Context) {
	logPrefix := common.GetLogPrefix(c)
	glog.Infof("%s fcn start", logPrefix)

	// Parse input
	form := &DataForm{}
	err := c.Bind(form)
	var dataFields []DataField
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
	measurementTypeIdToAttr := map[int64]attr{}
	for _, field := range dataFields {
		if id, ok := measurementTypeNameToId[field.Name]; ok {
			measurementTypeIds = append(measurementTypeIds, strconv.FormatInt(id, 10))
			measurementTypeIdToAttr[id] = field.Attr
			if field.Attr < 0 || field.Attr > DURATION {
				msg := "Couldn't parse input, unknown attr: " + string(field.Attr)
				glog.Errorf("%s %s", logPrefix, msg)
				c.JSON(http.StatusBadRequest, gin.H{"message": msg})
				return
			}
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
	glog.Infof("%s Measurements: %+v", logPrefix, measurements)

	// Generate output x values (with 0 values, a.k.a. not sparse)
	data := map[int64]float64{}
	drilldownData := map[int64]map[int64]float64{}
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
			data[newTs] = 0
			drilldownData[newTs] = map[int64]float64{}
			currTs = newTs
		}
		dayTs := tDay.Unix()
		if _, ok := drilldownData[currTs][dayTs]; !ok {
			drilldownData[currTs][dayTs] = 0
		}
		switch measurementTypeIdToAttr[m.MeasurementTypeId] {
		case VALUE:
			data[currTs] += m.Value
			drilldownData[currTs][dayTs] += m.Value
		case DURATION:
			data[currTs] += float64(m.Duration) / 3600
			drilldownData[currTs][dayTs] += float64(m.Duration) / 3600
		}
	}
	/*
	   type HCData struct {
	   	Name      string  `json:"name"`
	   	Drilldown string  `json:"drilldown"`
	   	Y         float64 `json:"y"`
	   }
	   type HCDrilldownData struct {
	   	Name string                 `json:"name"`
	   	Id   string                 `json:"id"`
	   	Data [][]interface{} `json:"data"`
	   }
	*/
	var hcData []HCData
	var hcDrilldownData []HCDrilldownData
	currTime := time.Unix(minTs, 0).In(loc)
	for currTime.Before(time.Now().In(loc)) {
		// get y value (default to 0)
		var y float64
		if val, ok := data[currTime.Unix()]; ok {
			y = val
		}
		// get x value (modify from first to last day of the period)
		var x string
		var nextTime time.Time
		switch form.Increment {
		case DAY:
			x = currTime.Format("2006-01-02")
			nextTime = currTime.AddDate(0, 0, 1)
		case WEEK:
			x = currTime.AddDate(0, 0, 6).Format("2006-01-02")
			nextTime = currTime.AddDate(0, 0, 7)
		case MONTH:
			x = currTime.AddDate(0, 1, -1).Format("Jan 2006")
			nextTime = currTime.AddDate(0, 1, 0)
		}
		hcData = append(hcData, HCData{Name: x, Drilldown: x, Y: y})

		// set drilldown data
		if dayData, ok := drilldownData[currTime.Unix()]; ok {
			hcDD := HCDrilldownData{
				Name: x,
				Id:   x,
				Data: [][]interface{}{},
			}
			drilldownTime := currTime
			for drilldownTime.Before(nextTime) {
				y = 0
				if val, ok := dayData[drilldownTime.Unix()]; ok {
					y = val
				}
				x = drilldownTime.Format("2006-01-02")
				hcDD.Data = append(hcDD.Data, []interface{}{x, y})
				// increment drilldownTime
				drilldownTime = drilldownTime.AddDate(0, 0, 1)
			}
			hcDrilldownData = append(hcDrilldownData, hcDD)
		}

		// increment currTime
		currTime = nextTime
	}

	c.JSON(http.StatusOK, gin.H{"data": hcData, "drilldown": hcDrilldownData})
	return
}
