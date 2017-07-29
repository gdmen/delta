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

func resetTestDB(t *testing.T) {
	var err error
	users := []User{
		// username, password
		User{Username: "username", Password: "$2a$10$UMBNySrXiZgARiK1l9m/F.ACV2MBOQPAglYluAHdBqsZBdahmMCTm"},
	}
	if _, err = api.DB.Exec(`delete from users;`); err != nil {
		t.Fatalf("Failed to delete users table: %v", err)
	}
	for _, u := range users {
		if _, err = api.DB.Exec(InsertUserSQL, u.Username, u.Password); err != nil {
			t.Fatalf("Failed to insert user(%s, %s): %v", u.Username, u.Password, err)
		}
	}
}

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
	//os.Remove(TestDB)
	os.Exit(ret)
}
