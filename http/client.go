package http

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
)

// HttpClient manage http requests
type HttpClient interface {
	Get(path string, params map[string]string) error
	Post(path string, body interface{}) error
}

// Request represents a http client
type Request struct {
	url          string
	headers      map[string]string
	client       *http.Client
	request      *http.Request
	responseBody io.ReadCloser
	StatusCode int
}

// New creates a new Request instance
func New(url string, client *http.Client) *Request {
	return &Request{
		url:     url,
		client:  client,
		headers: map[string]string{},
	}
}

// AddHeader increment http request headers
func (c *Request) AddHeader(key, value string) *Request {
	c.headers[key] = value
	return c
}

func (c *Request) setHeaders() {
	for key, value := range c.headers {
		c.request.Header.Set(key, value)
	}
}

// Post do a http post request
func (c *Request) Post(path string, body interface{}) error {
	data, err := json.Marshal(body)
	if err != nil {
		return err
	}

	request, err := http.NewRequest(http.MethodPost, c.url+path, bytes.NewBuffer(data))
	if err != nil {
		return err
	}

	c.request = request
	c.setHeaders()
	return c.execute()
}

// Get do a http get request
func (c *Request) Get(path string, params map[string]string) error {
	request, err := http.NewRequest(http.MethodGet, c.url+path, nil)
	if err != nil {
		return err
	}

	c.request = request
	c.setHeaders()
	c.queryBuilder(params)
	return c.execute()
}

func (c *Request) queryBuilder(params map[string]string) {
	q := c.request.URL.Query()
	for key, value := range params {
		q.Add(key, value)
	}
	c.request.URL.RawQuery = q.Encode()
}

func (c *Request) execute() error {
	response, err := c.client.Do(c.request)
	if err != nil {
		return err
	}

	c.StatusCode = response.StatusCode
	c.responseBody = response.Body
	return nil
}

func (c *Request) readBody() ([]byte, error) {
	defer c.responseBody.Close()

	body, err := ioutil.ReadAll(c.responseBody)
	if err != nil {
		return nil, err
	}
	return body, nil
}

// Json parse response body to json
func (c *Request) Json() (string, error) {
	body, err := c.readBody()
	if err != nil {
		return "", err
	}
	return string(body), nil
}

// Decode parse response body to a specific interface
func (c *Request) Decode(i interface{}) (interface{}, error) {
	body, err := c.readBody()
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(body, &i)
	if err != nil {
		return nil, err
	}

	return i, nil
}
