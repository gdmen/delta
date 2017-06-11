package api

import (
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"strconv"
	"strings"
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
		User{Id: 1, Username: "username", Password: "$2a$10$UMBNySrXiZgARiK1l9m/F.ACV2MBOQPAglYluAHdBqsZBdahmMCTm"},
	}
	if _, err = api.DB.Exec(`delete from users;`); err != nil {
		t.Fatalf("Failed to delete users table: %v", err)
	}
	for _, u := range users {
		if _, err = api.DB.Exec(`insert into users(id, username, password) values(?, ?, ?);`, u.Id, u.Username, u.Password); err != nil {
			t.Fatalf("Failed to insert user(%d, %s, %s): %v", u.Id, u.Username, u.Password, err)
		}
	}
}

// Set up a global test db and clean up after running all tests
func TestMain(m *testing.M) {
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
	os.Exit(ret)
}

func TestRegisterBasic(t *testing.T) {
	resetTestDB(t)
	resp := httptest.NewRecorder()

	r := api.GetRouter()

	values := url.Values{}
	values.Add("username", "newusername")
	values.Add("password", "pw")
	paramString := values.Encode()

	req, _ := http.NewRequest("POST", "/api/v1/u/register", strings.NewReader(paramString))
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Content-Length", strconv.Itoa(len(paramString)))

	r.ServeHTTP(resp, req)

	if resp.Code != http.StatusCreated {
		t.Fatal(resp)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}
	if strings.TrimSpace(string(body)) != "{\"user\":{\"id\":2,\"username\":\"newusername\"}}" {
		t.Fatal(string(body))
	}
}

func TestRegisterUnavailableUsername(t *testing.T) {
	resetTestDB(t)
	resp := httptest.NewRecorder()

	r := api.GetRouter()

	values := url.Values{}
	values.Add("username", "username")
	values.Add("password", "pw")
	paramString := values.Encode()

	req, _ := http.NewRequest("POST", "/api/v1/u/register", strings.NewReader(paramString))
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Content-Length", strconv.Itoa(len(paramString)))

	r.ServeHTTP(resp, req)

	if resp.Code != http.StatusBadRequest {
		t.Fatal(resp)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil || !strings.Contains(string(body), "Username isn't available") {
		t.Fatal(string(body))
	}
}
