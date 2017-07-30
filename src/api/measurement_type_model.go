package api

import (
	"fmt"
)

type MeasurementType struct {
	Id    int64  `json:"id"`
	Name  string `json:"name" form:"name" binding:"required"`
	Units string `json:"units" form:"units" binding:"required"`
}

func (m MeasurementType) String() string {
	return fmt.Sprintf("Id: %d, Name: %s, Units: %s", m.Id, m.Name, m.Units)
}

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
