package api

import (
	"bytes"
	"github.com/gdmen/delta/src/common"
	_ "github.com/go-sql-driver/mysql"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestImportFitocracyLifts(t *testing.T) {
	c, err := common.ReadConfig("../../test_conf.json")
	if err != nil {
		t.Fatalf("Couldn't read config: %v", err)
	}
	ResetTestApi(c)
	r := TestApi.GetRouter()

	resp := httptest.NewRecorder()

	path := "./test_data/fitocracy_lifts.csv"
	file, err := os.Open(path)
	if err != nil {
		t.Fatalf("Couldn't open test file: %v", err)
	}
	defer file.Close()

	reqBody := &bytes.Buffer{}
	writer := multipart.NewWriter(reqBody)
	part, err := writer.CreateFormFile("files", filepath.Base(path))
	if err != nil {
		t.Fatalf("Couldn't create form: %v", err)
	}
	_, err = io.Copy(part, file)
	if err != nil {
		t.Fatalf("Couldn't copy file to form: %v", err)
	}
	err = writer.Close()
	if err != nil {
		t.Fatalf("Couldn't close writer: %v", err)
	}

	req, err := http.NewRequest("POST", "/api/v1/import/fitocracy", reqBody)
	req.Header.Set("Content-Type", writer.FormDataContentType())

	r.ServeHTTP(resp, req)

	if resp.Code != http.StatusCreated {
		t.Fatalf("Expected status code %d, got %d. . .\n%+v", http.StatusCreated, resp.Code, resp)
	}

	// Verify MeasurementTypes
	resp = httptest.NewRecorder()

	req, _ = http.NewRequest("GET", "/api/v1/measurement_types/", nil)

	r.ServeHTTP(resp, req)

	if resp.Code != http.StatusOK {
		t.Fatalf("Expected status code %d, got %d. . .\n%+v", http.StatusOK, resp.Code, resp)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}
	if strings.TrimSpace(string(body)) != `{"measurement_types":[{"id":1,"name":"Running","units":""},{"id":2,"name":"Barbell Back Squat","units":"lbs"},{"id":3,"name":"Pull Up","units":""}]}` {
		t.Fatal("ERROR: " + string(body))
	}

	// Verify Measurements
	resp = httptest.NewRecorder()

	req, _ = http.NewRequest("GET", "/api/v1/measurements/", nil)

	r.ServeHTTP(resp, req)

	if resp.Code != http.StatusOK {
		t.Fatalf("Expected status code %d, got %d. . .\n%+v", http.StatusOK, resp.Code, resp)
	}

	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}
	if strings.TrimSpace(string(body)) != `{"measurements":[{"id":1,"measurement_type_id":1,"value":0,"repetitions":0,"start_time":1318057200,"duration":1800,"data_source":"fitocracy"},{"id":2,"measurement_type_id":1,"value":0.8,"repetitions":0,"start_time":1332658800,"duration":450,"data_source":"fitocracy"},{"id":3,"measurement_type_id":2,"value":45,"repetitions":5,"start_time":1337670000,"duration":0,"data_source":"fitocracy"},{"id":4,"measurement_type_id":2,"value":65,"repetitions":5,"start_time":1337670000,"duration":0,"data_source":"fitocracy"},{"id":5,"measurement_type_id":3,"value":0,"repetitions":5,"start_time":1345532400,"duration":0,"data_source":"fitocracy"},{"id":6,"measurement_type_id":3,"value":10,"repetitions":2,"start_time":1345705200,"duration":0,"data_source":"fitocracy"}]}` {
		t.Fatal("ERROR: " + string(body))
	}
}

func TestImportFitocracyTimed(t *testing.T) {
	c, err := common.ReadConfig("../../test_conf.json")
	if err != nil {
		t.Fatalf("Couldn't read config: %v", err)
	}
	ResetTestApi(c)
	r := TestApi.GetRouter()

	resp := httptest.NewRecorder()

	path := "./test_data/fitocracy_timed.csv"
	file, err := os.Open(path)
	if err != nil {
		t.Fatalf("Couldn't open test file: %v", err)
	}
	defer file.Close()

	reqBody := &bytes.Buffer{}
	writer := multipart.NewWriter(reqBody)
	part, err := writer.CreateFormFile("files", filepath.Base(path))
	if err != nil {
		t.Fatalf("Couldn't create form: %v", err)
	}
	_, err = io.Copy(part, file)
	if err != nil {
		t.Fatalf("Couldn't copy file to form: %v", err)
	}
	err = writer.Close()
	if err != nil {
		t.Fatalf("Couldn't close writer: %v", err)
	}

	req, err := http.NewRequest("POST", "/api/v1/import/fitocracy", reqBody)
	req.Header.Set("Content-Type", writer.FormDataContentType())

	r.ServeHTTP(resp, req)

	if resp.Code != http.StatusCreated {
		t.Fatalf("Expected status code %d, got %d. . .\n%+v", http.StatusCreated, resp.Code, resp)
	}

	// Verify MeasurementTypes
	resp = httptest.NewRecorder()

	req, _ = http.NewRequest("GET", "/api/v1/measurement_types/", nil)

	r.ServeHTTP(resp, req)

	if resp.Code != http.StatusOK {
		t.Fatalf("Expected status code %d, got %d. . .\n%+v", http.StatusOK, resp.Code, resp)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}
	if strings.TrimSpace(string(body)) != `{"measurement_types":[{"id":1,"name":"Judo","units":""},{"id":2,"name":"Yoga","units":""}]}` {
		t.Fatal("ERROR: " + string(body))
	}

	// Verify Measurements
	resp = httptest.NewRecorder()

	req, _ = http.NewRequest("GET", "/api/v1/measurements/", nil)

	r.ServeHTTP(resp, req)

	if resp.Code != http.StatusOK {
		t.Fatalf("Expected status code %d, got %d. . .\n%+v", http.StatusOK, resp.Code, resp)
	}

	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}
	if strings.TrimSpace(string(body)) != `{"measurements":[{"id":1,"measurement_type_id":1,"value":0,"repetitions":0,"start_time":1361606400,"duration":3600,"data_source":"fitocracy"},{"id":2,"measurement_type_id":1,"value":0,"repetitions":0,"start_time":1361865600,"duration":7200,"data_source":"fitocracy"},{"id":3,"measurement_type_id":2,"value":0,"repetitions":0,"start_time":1365750000,"duration":3600,"data_source":"fitocracy"}]}` {
		t.Fatal("ERROR: " + string(body))
	}
}
