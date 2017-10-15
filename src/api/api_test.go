package api

import (
	"database/sql"
	"flag"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"io/ioutil"
	"os"
	"testing"
)

const (
	TestDB      = "./test.db"
	TestDataSQL = "./test_data/populate.sql"
)

var TestApi *Api

// Set up a global test db and clean up after running all tests
func TestMain(m *testing.M) {
	flag.Set("alsologtostderr", "true")
	flag.Set("v", "100")
	flag.Parse()
	ResetTestApi()
	defer TestApi.DB.Close()
	ret := m.Run()
	os.Exit(ret)
}

func ResetTestApi() {
	os.Remove(TestDB)
	db, err := sql.Open("sqlite3", TestDB)
	if err != nil {
		fmt.Errorf("Couldn't create db: %v", err)
		os.Exit(1)
	}
	TestApi, err = NewApi(db)
	if err != nil {
		fmt.Errorf("Couldn't init Api: %v", err)
		os.Exit(1)
	}
}

func PopulateTestApi() {
	sqlBytes, err := ioutil.ReadFile(TestDataSQL)
	if err != nil {
		fmt.Errorf("Couldn't read test data SQL: %v", err)
		os.Exit(1)
	}
	fmt.Printf(string(sqlBytes))
	_, err = TestApi.DB.Exec(string(sqlBytes))
	if err != nil {
		fmt.Errorf("Couldn't populate db: %v", err)
		os.Exit(1)
	}
}
