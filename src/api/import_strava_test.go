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

func TestImportStravaBasic(t *testing.T) {
	ResetTestApi()
	r := TestApi.GetRouter()

	resp := httptest.NewRecorder()

	path := "./test_data/strava.gpx"
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

	req, err := http.NewRequest("POST", "/api/v1/import/strava", reqBody)
	req.Header.Set("Content-Type", writer.FormDataContentType())

	r.ServeHTTP(resp, req)

	if resp.Code != http.StatusCreated {
		t.Fatalf("Expected status code %d, got %d. . .\n%v+", http.StatusCreated, resp.Code, resp)
	}

	// Verify MeasurementTypes
	resp = httptest.NewRecorder()

	req, _ = http.NewRequest("GET", "/api/v1/measurement_types/", nil)

	r.ServeHTTP(resp, req)

	if resp.Code != http.StatusOK {
		t.Fatalf("Expected status code %d, got %d. . .\n%v+", http.StatusOK, resp.Code, resp)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}
	if strings.TrimSpace(string(body)) != `{"measurement_types":[{"id":1,"name":"Road Cycling","units":"mi"}]}` {
		t.Fatal("ERROR: " + string(body))
	}

	// Verify Measurements
	resp = httptest.NewRecorder()

	req, _ = http.NewRequest("GET", "/api/v1/measurements/", nil)

	r.ServeHTTP(resp, req)

	if resp.Code != http.StatusOK {
		t.Fatalf("Expected status code %d, got %d. . .\n%v+", http.StatusOK, resp.Code, resp)
	}

	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}
	if strings.TrimSpace(string(body)) != `{"measurements":[{"id":1,"measurement_type_id":1,"value":6.1154513702093265,"repetitions":0,"start_time":1496243985,"duration":1886,"data_source":"strava"}]}` {
		t.Fatal("ERROR: " + string(body))
	}
}
