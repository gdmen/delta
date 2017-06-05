package api

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strconv"
	"strings"
	"testing"
)

func TestRegisterBasic(t *testing.T) {
	resp := httptest.NewRecorder()

	r := GetRouterV1()

	values := url.Values{}
	values.Add("username", "u1")
	values.Add("password", "p1")
	paramString := values.Encode()

	req, _ := http.NewRequest("POST", "/api/v1/u/register", strings.NewReader(paramString))
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Content-Length", strconv.Itoa(len(paramString)))

	r.ServeHTTP(resp, req)

	if resp.Code != http.StatusCreated {
		t.Fail()
	}
}

func TestRegisterUnavailableUsername(t *testing.T) {
	resp := httptest.NewRecorder()

	r := GetRouterV1()

	values := url.Values{}
	values.Add("username", "username")
	values.Add("password", "password")
	paramString := values.Encode()

	req, _ := http.NewRequest("POST", "/api/v1/u/register", strings.NewReader(paramString))
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Content-Length", strconv.Itoa(len(paramString)))

	r.ServeHTTP(resp, req)

	if resp.Code != http.StatusBadRequest {
		t.Fail()
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil || !strings.Contains(string(body), UsernameUnavailableUserErrMsg) {
		t.Fail()
	}
}
