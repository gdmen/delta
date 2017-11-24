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
	CREATE TABLE measurements (
		id INT AUTO_INCREMENT PRIMARY KEY,
		measurement_type_id INT NOT NULL,
		value DOUBLE NOT NULL,
		repetitions SMALLINT NOT NULL,
		start_time BIGINT NOT NULL,
		duration INT NOT NULL,
		data_source VARCHAR(64) NOT NULL,
		FOREIGN KEY(measurement_type_id) REFERENCES measurement_types(id)
	);`
	CreateMeasurementSQL = `
	INSERT INTO measurements(measurement_type_id, value, repetitions, start_time, duration, data_source) VALUES(?, ?, ?, ?, ?, ?);`
	ExistsMeasurementSQL = `
	SELECT id FROM measurements WHERE measurement_type_id=? AND value=? AND repetitions=? AND start_time=? AND duration=? AND data_source=?;`
	UpdateMeasurementSQL = `
	UPDATE measurements SET measurement_type_id=?, value=?, repetitions=?, start_time=?, duration=?, data_source=? WHERE id=?;`
	DeleteMeasurementSQL = `
	DELETE FROM measurements WHERE id=?;`
	GetMeasurementSQL = `
	SELECT * FROM measurements WHERE id=?;`
	ListMeasurementSQL = `
	SELECT * FROM measurements;`
)

type Measurement struct {
	Id                int64   `json:"id"`
	MeasurementTypeId int64   `json:"measurement_type_id" form:"measurement_type_id"`
	Value             float64 `json:"value" form:"value"`
	Repetitions       int16   `json:"repetitions" form:"repetitions"`
	StartTime         int64   `json:"start_time" form:"start_time"`
	Duration          int32   `json:"duration" form:"duration"`
	DataSource        string  `json:"data_source" form:"data_source"`
}

func (m Measurement) String() string {
	return fmt.Sprintf("Id: %d, MeasurementTypeId: %d, Value: %f, Repetitions: %d, StartTime: %d, Duration: %d, DataSource: %s", m.Id, m.MeasurementTypeId, m.Value, m.Repetitions, m.StartTime, m.Duration, m.DataSource)
}

type MeasurementManager struct {
	DB *sql.DB
}

func (m *MeasurementManager) Create(model *Measurement) (int, string, error) {
	// Check for existence
	var id int64
	err := m.DB.QueryRow(ExistsMeasurementSQL, model.MeasurementTypeId, model.Value, model.Repetitions, model.StartTime, model.Duration, model.DataSource).Scan(&id)
	if err == nil {
		model.Id = id
		return http.StatusCreated, "", nil
	}
	// Doesn't exist, so try to add it
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
	id, _ = result.LastInsertId()
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

func (m *MeasurementManager) Custom(sql string) (*[]Measurement, int, string, error) {
	models := []Measurement{}
	rows, err := m.DB.Query(sql)
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
