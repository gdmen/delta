package api

import (
	"fmt"
)

type Measurement struct {
	Id                int64  `json:"id"`
	MeasurementTypeId int64  `json:"measurement_type_id" form:"measurement_type_id"`
	Value             int64  `json:"value" form:"value"`
	Repetitions       int16  `json:"repetitions" form:"repetitions"`
	StartTime         int64  `json:"start_time" form:"start_time"`
	Duration          int32  `json:"duration" form:"duration"`
	DataSource        string `json:"data_source" form:"data_source"`
}

func (m Measurement) String() string {
	return fmt.Sprintf("Id: %d, MeasurementTypeId: %d, Value: %d, Repetitions: %d, StartTime: %d, Duration: %d, DataSource: %s", m.Id, m.MeasurementTypeId, m.Value, m.Repetitions, m.StartTime, m.Duration, m.DataSource)
}

const (
	// The id lines should be 'bigint' instead of 'integer'
	// but sqlite3 has a fucky primary key system.
	CreateMeasurementTableSQL = `
create table measurements (
	id integer primary key,
	measurement_type_id integer not null,
	value bigint not null,
	repetitions smallint not null,
	start_time bigint not null,
	duration int not null,
	data_source varchar not null

);`
	CreateMeasurementSQL = `
insert into measurements(measurement_type_id, value, repetitions, start_time, duration, data_source) values(?, ?, ?, ?, ?, ?);`
	UpdateMeasurementSQL = `
update measurements set measurement_type_id=?, value=?, repetitions=?, start_time=?, duration=?, data_source=? where id=?;`
	DeleteMeasurementSQL = `
delete from measurements where id=?;`
	GetMeasurementSQL = `
select * from measurements where id=?;`
	ListMeasurementSQL = `
select * from measurements;`
)
