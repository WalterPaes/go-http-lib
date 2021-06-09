package http

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"time"
)

type Client struct {
	url          string
	headers      map[string]string
	client       *http.Client
	request      *http.Request
	responseBody io.ReadCloser
	StatusCode   int
}

func New(url string) *Client {
	return &Client{
		url:    url,
		client: &http.Client{},
	}
}

func (c *Client) AddHeader(key, value string) *Client {
	c.headers[key] = value
	return c
}

func (c *Client) SetTimeout(timeout int64) *Client {
	c.client.Timeout = time.Duration(timeout)
	return c
}

func (c *Client) Post(path string, body interface{}) (*Client, error) {
	data, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}

	request, err := http.NewRequest(http.MethodPost, c.url+path, bytes.NewBuffer(data))
	if err != nil {
		return nil, err
	}

	c.request = request
	c.setHeaders()
	return c.execute()
}

func (c *Client) Get(path string, params map[string]string) (*Client, error) {
	request, err := http.NewRequest(http.MethodGet, c.url+path, nil)
	if err != nil {
		return nil, err
	}

	c.request = request
	c.setHeaders()
	c.queryBuilder(params)
	return c.execute()
}

func (c *Client) setHeaders() {
	for key, value := range c.headers {
		c.request.Header.Set(key, value)
	}
}

func (c *Client) execute() (*Client, error) {
	response, err := c.client.Do(c.request)
	if err != nil {
		return nil, err
	}

	c.StatusCode = response.StatusCode
	c.responseBody = response.Body
	return c, nil
}

func (c *Client) Decode(i interface{}) (interface{}, error) {
	defer c.responseBody.Close()

	body, err := ioutil.ReadAll(c.responseBody)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(body, &i)
	if err != nil {
		return nil, err
	}

	return i, nil
}

func (c *Client) queryBuilder(params map[string]string) {
	q := c.request.URL.Query()
	for key, value := range params {
		q.Add(key, value)
	}
	c.request.URL.RawQuery = q.Encode()
}
