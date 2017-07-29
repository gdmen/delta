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

func TestRegisterBasic(t *testing.T) {
	resetTestDB(t)
	resp := httptest.NewRecorder()

	r := api.GetRouter()

	values := url.Values{}
	values.Add("username", "newusername")
	values.Add("password", "pw")
	paramString := values.Encode()

	req, _ := http.NewRequest("POST", "/api/v1/users/register", strings.NewReader(paramString))
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

	req, _ := http.NewRequest("POST", "/api/v1/users/register", strings.NewReader(paramString))
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
