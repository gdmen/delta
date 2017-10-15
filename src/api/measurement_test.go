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

func TestMeasurementBasic(t *testing.T) {
	ResetTestApi()
	r := TestApi.GetRouter()

	// Create
	resp := httptest.NewRecorder()

	values := url.Values{}
	values.Add("measurement_type_id", "1")
	values.Add("value", "405")
	values.Add("repetitions", "1")
	values.Add("data_source", "fitnotes")
	paramString := values.Encode()

	req, _ := http.NewRequest("POST", "/api/v1/measurements/", strings.NewReader(paramString))
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Content-Length", strconv.Itoa(len(paramString)))

	r.ServeHTTP(resp, req)

	if resp.Code != http.StatusCreated {
		t.Fatalf("Expected status code %d, got %d. . .\n%+v", http.StatusCreated, resp.Code, resp)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}
	if strings.TrimSpace(string(body)) != `{"measurement":{"id":1,"measurement_type_id":1,"value":405,"repetitions":1,"start_time":0,"duration":0,"data_source":"fitnotes"}}` {
		t.Fatal(string(body))
	}

	// List
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
	if strings.TrimSpace(string(body)) != `{"measurements":[{"id":1,"measurement_type_id":1,"value":405,"repetitions":1,"start_time":0,"duration":0,"data_source":"fitnotes"}]}` {
		t.Fatal(string(body))
	}

	// Update
	resp = httptest.NewRecorder()

	values.Set("value", "315")
	values.Set("repetitions", "20")
	paramString = values.Encode()

	req, _ = http.NewRequest("POST", "/api/v1/measurements/1", strings.NewReader(paramString))
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Content-Length", strconv.Itoa(len(paramString)))

	r.ServeHTTP(resp, req)

	if resp.Code != http.StatusOK {
		t.Fatalf("Expected status code %d, got %d. . .\n%+v", http.StatusOK, resp.Code, resp)
	}

	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}
	if strings.TrimSpace(string(body)) != `{"measurement":{"id":1,"measurement_type_id":1,"value":315,"repetitions":20,"start_time":0,"duration":0,"data_source":"fitnotes"}}` {
		t.Fatal(string(body))
	}

	// Get
	resp = httptest.NewRecorder()

	req, _ = http.NewRequest("GET", "/api/v1/measurements/1", nil)

	r.ServeHTTP(resp, req)

	if resp.Code != http.StatusOK {
		t.Fatalf("Expected status code %d, got %d. . .\n%+v", http.StatusOK, resp.Code, resp)
	}

	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}
	if strings.TrimSpace(string(body)) != `{"measurement":{"id":1,"measurement_type_id":1,"value":315,"repetitions":20,"start_time":0,"duration":0,"data_source":"fitnotes"}}` {
		t.Fatal(string(body))
	}

	// Delete
	resp = httptest.NewRecorder()

	req, _ = http.NewRequest("DELETE", "/api/v1/measurements/1", nil)

	r.ServeHTTP(resp, req)

	if resp.Code != http.StatusNoContent {
		t.Fatalf("Expected status code %d, got %d. . .\n%+v", http.StatusNoContent, resp.Code, resp)
	}

	// List
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
	if strings.TrimSpace(string(body)) != "{\"measurements\":[]}" {
		t.Fatal(string(body))
	}
}
