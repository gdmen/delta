package api

import (
	_ "github.com/mattn/go-sqlite3"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strconv"
	"strings"
	"testing"
)

func TestMeasurementTypeBasic(t *testing.T) {
	r := api.GetRouter()

	// Create
	resp := httptest.NewRecorder()

	values := url.Values{}
	values.Add("name", "barbell back squat")
	values.Add("units", "lbs")
	paramString := values.Encode()

	req, _ := http.NewRequest("POST", "/api/v1/measurement_types/", strings.NewReader(paramString))
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Content-Length", strconv.Itoa(len(paramString)))

	r.ServeHTTP(resp, req)

	if resp.Code != http.StatusCreated {
		t.Fatalf("Expected status code %d, got %d. . .\n%v+", http.StatusCreated, resp.Code, resp)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}
	if strings.TrimSpace(string(body)) != "{\"measurement_type\":{\"id\":1,\"name\":\"barbell back squat\",\"units\":\"lbs\"}}" {
		t.Fatal(string(body))
	}

	// List
	resp = httptest.NewRecorder()

	req, _ = http.NewRequest("GET", "/api/v1/measurement_types/", nil)

	r.ServeHTTP(resp, req)

	if resp.Code != http.StatusOK {
		t.Fatalf("Expected status code %d, got %d. . .\n%v+", http.StatusOK, resp.Code, resp)
	}

	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}
	if strings.TrimSpace(string(body)) != "{\"measurement_types\":[{\"id\":1,\"name\":\"barbell back squat\",\"units\":\"lbs\"}]}" {
		t.Fatal(string(body))
	}

	// Update
	resp = httptest.NewRecorder()

	values = url.Values{}
	values.Add("name", "deadlift")
	values.Add("units", "kg")
	paramString = values.Encode()

	req, _ = http.NewRequest("POST", "/api/v1/measurement_types/1", strings.NewReader(paramString))
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Content-Length", strconv.Itoa(len(paramString)))

	r.ServeHTTP(resp, req)

	if resp.Code != http.StatusOK {
		t.Fatalf("Expected status code %d, got %d. . .\n%v+", http.StatusOK, resp.Code, resp)
	}

	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}
	if strings.TrimSpace(string(body)) != "{\"measurement_type\":{\"id\":1,\"name\":\"deadlift\",\"units\":\"kg\"}}" {
		t.Fatal(string(body))
	}

	// Get
	resp = httptest.NewRecorder()

	req, _ = http.NewRequest("GET", "/api/v1/measurement_types/1", nil)

	r.ServeHTTP(resp, req)

	if resp.Code != http.StatusOK {
		t.Fatalf("Expected status code %d, got %d. . .\n%v+", http.StatusOK, resp.Code, resp)
	}

	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}
	if strings.TrimSpace(string(body)) != "{\"measurement_type\":{\"id\":1,\"name\":\"deadlift\",\"units\":\"kg\"}}" {
		t.Fatal(string(body))
	}

	// Delete
	resp = httptest.NewRecorder()

	req, _ = http.NewRequest("DELETE", "/api/v1/measurement_types/1", nil)

	r.ServeHTTP(resp, req)

	if resp.Code != http.StatusNoContent {
		t.Fatalf("Expected status code %d, got %d. . .\n%v+", http.StatusNoContent, resp.Code, resp)
	}

	// List
	resp = httptest.NewRecorder()

	req, _ = http.NewRequest("GET", "/api/v1/measurement_types/", nil)

	r.ServeHTTP(resp, req)

	if resp.Code != http.StatusOK {
		t.Fatalf("Expected status code %d, got %d. . .\n%v+", http.StatusOK, resp.Code, resp)
	}

	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}
	if strings.TrimSpace(string(body)) != "{\"measurement_types\":[]}" {
		t.Fatal(string(body))
	}
}
