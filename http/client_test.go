package http

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

func TestNew(t *testing.T) {
	client, mux, server := testServer()
	defer server.Close()
	mux.HandleFunc("/success", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, `{"text": "Some text", "favorite_count": 24}`)
	})

	request := New("http://example.com", client)
	err := request.Get("/success", nil)
	if err != nil {
		t.Fatal(err.Error())
	}

	result, err := request.Decode(map[string]string{})
	if err != nil {
		t.Fatal(err.Error())
	}

	fmt.Println(result)
}

func testServer() (*http.Client, *http.ServeMux, *httptest.Server) {
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
