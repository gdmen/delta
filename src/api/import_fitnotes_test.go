package api

import (
	"bytes"
	_ "github.com/mattn/go-sqlite3"
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

func TestImportFitnotesBasic(t *testing.T) {
	ResetTestApi()
	r := TestApi.GetRouter()

	resp := httptest.NewRecorder()

	path := "./test_data/fitnotes.csv"
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

	req, err := http.NewRequest("POST", "/api/v1/import/fitnotes", reqBody)
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
	if strings.TrimSpace(string(body)) != `{"measurement_types":[{"id":1,"name":"Brazilian Jiu-Jitsu","units":""},{"id":2,"name":"Gymnastics","units":""},{"id":3,"name":"Road Cycling","units":"mi"},{"id":4,"name":"Barbell Back Squat","units":"lbs"},{"id":5,"name":"Pull Up","units":"lbs"}]}` {
		t.Fatal(string(body))
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
	if strings.TrimSpace(string(body)) != `{"measurements":[{"id":1,"measurement_type_id":1,"value":0,"repetitions":0,"start_time":1479369600,"duration":9000,"data_source":"fitnotes"},{"id":2,"measurement_type_id":2,"value":0,"repetitions":0,"start_time":1481011200,"duration":4800,"data_source":"fitnotes"},{"id":3,"measurement_type_id":3,"value":3.3,"repetitions":0,"start_time":1481097600,"duration":0,"data_source":"fitnotes"},{"id":4,"measurement_type_id":4,"value":95,"repetitions":10,"start_time":1484553600,"duration":0,"data_source":"fitnotes"},{"id":5,"measurement_type_id":4,"value":145,"repetitions":5,"start_time":1484553600,"duration":0,"data_source":"fitnotes"},{"id":6,"measurement_type_id":5,"value":0,"repetitions":7,"start_time":1484553600,"duration":0,"data_source":"fitnotes"},{"id":7,"measurement_type_id":5,"value":0,"repetitions":3,"start_time":1484553600,"duration":0,"data_source":"fitnotes"}]}` {
		t.Fatal(string(body))
	}
}
