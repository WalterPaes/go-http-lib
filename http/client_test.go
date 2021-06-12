package http

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"reflect"
	"testing"
)

const (
	baseUrl         = "http://example.com"
	successPath     = "/success"
	failPath        = "/fail"
	responseSuccess = `{"message": "Success"}`
	responseFail    = `{"message": "Fail"}`
)

func TestNew(t *testing.T) {
	request := New(baseUrl, &http.Client{})

	expected := "*http.Request"
	got := fmt.Sprintf("%T", request)

	if got != expected {
		t.Errorf("Was expected '%s', but got '%s'", expected, got)
	}
}

func TestRequest_AddHeader(t *testing.T) {
	request := New(baseUrl, &http.Client{})
	headers := map[string]string{
		"abc":   "ced",
		"123":   "321",
		"hello": "goodbye",
	}

	for key, value := range headers {
		request.AddHeader(key, value)
	}

	if !reflect.DeepEqual(headers, request.headers) {
		t.Errorf("Was expected '%v', but got '%v'", headers, request.headers)
	}
}

func TestRequest_Post(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		client, mux, server := testServer(t)
		defer server.Close()
		handleFunc(t, mux, successPath, responseSuccess)

		request := New(baseUrl, client)
		err := request.Get(successPath, nil)
		if err != nil {
			t.Fatal(err.Error())
		}

		response, err := request.Json()
		if err != nil {
			t.Fatal(err.Error())
		}

		if response != responseSuccess {
			t.Errorf("Was expected '%s', but got '%s'", response, responseSuccess)
		}
	})
}

func testServer(t *testing.T) (*http.Client, *http.ServeMux, *httptest.Server) {
	t.Helper()
	mux := http.NewServeMux()
	server := httptest.NewServer(mux)
	transport := &http.Transport{
		Proxy: func(req *http.Request) (*url.URL, error) {
			return url.Parse(server.URL)
		},
	}
	client := &http.Client{Transport: transport}
	return client, mux, server
}

func handleFunc(t *testing.T, mux *http.ServeMux, path, message string) {
	t.Helper()
	mux.HandleFunc(path, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, message)
	})
}
