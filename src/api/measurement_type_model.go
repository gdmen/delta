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
	CreateMeasurementTypeTableSQL = `
create table measurement_types (
	id integer primary key,
	name varchar not null,
	units varchar not null,
	constraint name unique (name)
);`
	CreateMeasurementTypeSQL = `
insert into measurement_types(name, units) values(?, ?);`
	UpdateMeasurementTypeSQL = `
update measurement_types set name=?, units=? where id=?;`
	DeleteMeasurementTypeSQL = `
delete from measurement_types where id=?;`
	GetMeasurementTypeSQL = `
select * from measurement_types where id=?;`
	ListMeasurementTypeSQL = `
select * from measurement_types;`
)

type MeasurementType struct {
	Id    int64  `json:"id"`
	Name  string `json:"name" form:"name" binding:"required"`
	Units string `json:"units" form:"units" binding:"required"`
}

func (m *MeasurementType) String() string {
	return fmt.Sprintf("Id: %d, Name: %s, Units: %s", m.Id, m.Name, m.Units)
}

type MeasurementTypeManager struct {
	DB *sql.DB
}

func (m *MeasurementTypeManager) Create(model *MeasurementType) (int, string, error) {
	result, err := m.DB.Exec(CreateMeasurementTypeSQL, model.Name, model.Units)
	if err != nil {
		if strings.Contains(err.Error(), "UNIQUE constraint failed") {
			msg := "That measurement_type already exits"
			return http.StatusBadRequest, msg, err
		}
		msg := "Couldn't add measurement_type to database"
		return http.StatusInternalServerError, msg, err
	}
	// Get the Id that was just auto-written to the database
	// Ignore errors (if the database doesn't support LastInsertId)
	id, _ := result.LastInsertId()
	model.Id = id
	return http.StatusCreated, "", nil
}

func (m *MeasurementTypeManager) Update(model *MeasurementType) (int, string, error) {
	result, err := m.DB.Exec(UpdateMeasurementTypeSQL, model.Name, model.Units, model.Id)
	if err != nil {
		msg := "Couldn't update measurement_type in database"
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

func (m *MeasurementTypeManager) Delete(id int64) (int, string, error) {
	result, err := m.DB.Exec(DeleteMeasurementTypeSQL, id)
	if err != nil {
		msg := "Couldn't delete measurement_type in database"
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

func (m *MeasurementTypeManager) Get(id int64) (*MeasurementType, int, string, error) {
	model := &MeasurementType{}
	err := m.DB.QueryRow(GetMeasurementTypeSQL, id).Scan(&model.Id, &model.Name, &model.Units)
	if err == sql.ErrNoRows {
		msg := "Couldn't find a measurement_type with that Id"
		return nil, http.StatusNotFound, msg, err
	} else if err != nil {
		msg := "Couldn't get measurement_type from database"
		return nil, http.StatusInternalServerError, msg, err
	}
	return model, http.StatusOK, "", nil
}

func (m *MeasurementTypeManager) List() (*[]MeasurementType, int, string, error) {
	models := []MeasurementType{}
	rows, err := m.DB.Query(ListMeasurementTypeSQL)
	defer rows.Close()
	if err != nil {
		msg := "Couldn't get measurement_types from database"
		return nil, http.StatusInternalServerError, msg, err
	}
	for rows.Next() {
		model := MeasurementType{}
		err = rows.Scan(&model.Id, &model.Name, &model.Units)
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
