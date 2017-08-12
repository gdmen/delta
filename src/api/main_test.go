package api

import (
	"database/sql"
	"flag"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"os"
	"testing"
)

const (
	TestDB = "./test.db"
)

var api *Api

// Set up a global test db and clean up after running all tests
func TestMain(m *testing.M) {
	flag.Set("alsologtostderr", "true")
	flag.Set("v", "3")
	flag.Parse()
	os.Remove(TestDB)
	db, err := sql.Open("sqlite3", TestDB)
	if err != nil {
		fmt.Errorf("Couldn't create db: %v", err)
		os.Exit(1)
	}
	api, err = NewApi(db)
	if err != nil {
		fmt.Errorf("Couldn't init Api: %v", err)
		os.Exit(1)
	}
	ret := m.Run()
	db.Close()
	os.Remove(TestDB)
	os.Exit(ret)
}
