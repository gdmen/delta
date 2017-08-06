package api

import (
	"database/sql"
	"fmt"
	"net/http"
	"strings"
)

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

type MeasurementManager struct {
	DB *sql.DB
}

func (m *MeasurementManager) Create(model *Measurement) (int, string, error) {
	result, err := m.DB.Exec(CreateMeasurementSQL, model.MeasurementTypeId, model.Value, model.Repetitions, model.StartTime, model.Duration, model.DataSource)
	if err != nil {
		if strings.Contains(err.Error(), "UNIQUE constraint failed") {
			msg := "That measurement already exits"
			return http.StatusBadRequest, msg, err
		}
		msg := "Couldn't add measurement to database"
		return http.StatusInternalServerError, msg, err
	}
	// Get the Id that was just auto-written to the database
	// Ignore errors (if the database doesn't support LastInsertId)
	id, _ := result.LastInsertId()
	model.Id = id
	return http.StatusCreated, "", nil
}

func (m *MeasurementManager) Update(model *Measurement) (int, string, error) {
	result, err := m.DB.Exec(UpdateMeasurementSQL, model.MeasurementTypeId, model.Value, model.Repetitions, model.StartTime, model.Duration, model.DataSource, model.Id)
	if err != nil {
		msg := "Couldn't update measurement in database"
		return http.StatusInternalServerError, msg, err
	}
	// Check for 404s
	// Ignore errors (if the database doesn't support RowsAffected)
	affected, _ := result.RowsAffected()
	if affected == 0 {
		return http.StatusNotFound, "", nil
	}
	return http.StatusOK, "", nil
}

func (m *MeasurementManager) Delete(id int64) (int, string, error) {
	result, err := m.DB.Exec(DeleteMeasurementSQL, id)
	if err != nil {
		msg := "Couldn't delete measurement in database"
		return http.StatusInternalServerError, msg, err
	}
	// Check for 404s
	// Ignore errors (if the database doesn't support RowsAffected)
	affected, _ := result.RowsAffected()
	if affected == 0 {
		return http.StatusNotFound, "", nil
	}
	return http.StatusNoContent, "", nil
}

func (m *MeasurementManager) Get(id int64) (*Measurement, int, string, error) {
	model := &Measurement{}
	err := m.DB.QueryRow(GetMeasurementSQL, id).Scan(&model.Id, &model.MeasurementTypeId, &model.Value, &model.Repetitions, &model.StartTime, &model.Duration, &model.DataSource)
	if err == sql.ErrNoRows {
		msg := "Couldn't find a measurement with that Id"
		return nil, http.StatusNotFound, msg, err
	} else if err != nil {
		msg := "Couldn't get measurement from database"
		return nil, http.StatusInternalServerError, msg, err
	}
	return model, http.StatusOK, "", nil
}

func (m *MeasurementManager) List() (*[]Measurement, int, string, error) {
	models := []Measurement{}
	rows, err := m.DB.Query(ListMeasurementSQL)
	defer rows.Close()
	if err != nil {
		msg := "Couldn't get measurements from database"
		return nil, http.StatusInternalServerError, msg, err
	}
	for rows.Next() {
		model := Measurement{}
		err = rows.Scan(&model.Id, &model.MeasurementTypeId, &model.Value, &model.Repetitions, &model.StartTime, &model.Duration, &model.DataSource)
		if err != nil {
			msg := "Couldn't scan row from database"
			return nil, http.StatusInternalServerError, msg, err
		}
		models = append(models, model)
	}
	err = rows.Err()
	if err != nil {
		msg := "Error scanning rows from database"
		return nil, http.StatusInternalServerError, msg, err
	}
	return &models, http.StatusOK, "", nil
}
