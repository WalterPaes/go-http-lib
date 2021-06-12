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

func TestRequest_GetJsonDecode(t *testing.T) {
	cases := []struct {
		title    string
		url      string
		path     string
		response string
	}{
		{
			"Success request: response json decode",
			baseUrl,
			successPath,
			responseSuccess,
		},
		{
			"Fail request: response json decode",
			baseUrl,
			failPath,
			responseFail,
		},
	}

	assert := func(t *testing.T, title, baseUrl, path, expectedResponse string) {
		t.Run(title, func(t *testing.T) {
			request := createServerAndRequest(t, path, expectedResponse, baseUrl)
			response, err := request.Json()
			if err != nil {
				t.Fatal(err.Error())
			}

			if response != expectedResponse {
				t.Errorf("Was expected '%s', but got '%s'", response, expectedResponse)
			}
		})
	}

	for _, c := range cases {
		assert(t, c.title, c.url, c.path, c.response)
	}
}

//func TestRequest_GetInterfaceDecode(t *testing.T) {
//	cases := []struct {
//		title    string
//		url      string
//		path     string
//		response string
//	}{
//		{
//			"Success request: response interface decode",
//			baseUrl,
//			successPath,
//			responseSuccess,
//		},
//		{
//			"Fail request: response interface decode",
//			baseUrl,
//			failPath,
//			responseFail,
//		},
//	}
//
//	assert := func(t *testing.T, title, baseUrl, path, expectedResponse string) {
//		t.Run(title, func(t *testing.T) {
//			data := make(map[string]string)
//			request := createServerAndRequest(t, path, expectedResponse, baseUrl)
//			response, err := request.Decode(data)
//			if err != nil {
//				t.Fatal(err.Error())
//			}
//
//			err = json.Unmarshal([]byte(expectedResponse), &data)
//			if err != nil {
//				t.Fatal(err.Error())
//			}
//
//			if !reflect.DeepEqual(response, data) {
//				t.Errorf("Was expected '%+v', but got '%+v'", data, response)
//			}
//		})
//	}
//
//	for _, c := range cases {
//		assert(t, c.title, c.url, c.path, c.response)
//	}
//}

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

func createServerAndRequest(t *testing.T, path, expectedResponse, baseUrl string) *Request {
	t.Helper()
	client, mux, server := testServer(t)
	defer server.Close()
	handleFunc(t, mux, path, expectedResponse)

	request := New(baseUrl, client)
	err := request.Get(path, nil)
	if err != nil {
		t.Fatal(err.Error())
	}
	return request
}
