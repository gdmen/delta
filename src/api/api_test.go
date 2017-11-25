package api

import (
	"database/sql"
	"flag"
	"fmt"
	"github.com/gdmen/delta/src/common"
	_ "github.com/go-sql-driver/mysql"
	"os"
	"os/exec"
	"testing"
)

const (
	TestDataDir          = "./test_data/"
	TestDataInsertScript = "insert_data.sh"
)

var TestApi *Api

// Set up a global test db and clean up after running all tests
func TestMain(m *testing.M) {
	flag.Set("alsologtostderr", "true")
	flag.Set("v", "100")
	flag.Parse()
	c, err := common.ReadConfig("../../test_conf.json")
	if err != nil {
		fmt.Printf("Couldn't read config: %v", err)
		os.Exit(1)
	}
	ResetTestApi(c)
	defer TestApi.DB.Close()
	ret := m.Run()
	os.Exit(ret)
}

func ResetTestApi(c *common.Config) {
	connectStr := fmt.Sprintf("%s:%s@tcp(%s:%s)/?charset=utf8", c.MySQLUser, c.MySQLPass, c.MySQLHost, c.MySQLPort)
	db, err := sql.Open("mysql", connectStr)
	if err != nil {
		fmt.Printf("Couldn't connect to db: %v", err)
		os.Exit(1)
	}
	db.Exec("DROP DATABASE delta_test;")
	db.Exec("CREATE DATABASE delta_test;")
	db.Close()
	// Reconnect specifically to the test database
	connectStr = fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8", c.MySQLUser, c.MySQLPass, c.MySQLHost, c.MySQLPort, c.MySQLDatabase)
	db, err = sql.Open("mysql", connectStr)
	if err != nil {
		fmt.Printf("Couldn't connect to db: %v", err)
		os.Exit(1)
	}
	TestApi, err = NewApi(db)
	if err != nil {
		fmt.Printf("Couldn't init Api: %v", err)
		os.Exit(1)
	}
}

func InsertTestMeasurementTypes(c *common.Config) {
	insertTestData(c, "measurement_types")
}

func InsertTestMeasurements(c *common.Config) {
	insertTestData(c, "measurements")
}

func insertTestData(c *common.Config, tableName string) {
	cmd := exec.Command(
		TestDataDir+TestDataInsertScript, c.MySQLUser, c.MySQLPass, c.MySQLDatabase,
		TestDataDir+tableName+".sql")
	err := cmd.Run()
	if err != nil {
		fmt.Printf("Couldn't populate %s in db: %v", tableName, err)
		os.Exit(1)
	}
}
