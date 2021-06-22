package api

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
)

const (
	baseUrl = "https://api.fing.ir/v1/"
)

type Client struct {
	dev     bool
	client  *http.Client
	BaseURL *url.URL

	headers     map[string]string
	accessToken string
}

func NewClient(token string, options ...Option) *Client {

	baseUrl, _ := url.Parse(baseUrl)

	c := &Client{
		client: &http.Client{
			Transport: http.DefaultTransport,
		},
		BaseURL:     baseUrl,
		headers:     make(map[string]string),
		accessToken: token,
	}

	for _, opt := range options {
		opt(c)
	}

	return c
}

func (c *Client) doRequest(req *http.Request, res interface{}) error {
	resp, err := c.client.Do(req)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	return json.NewDecoder(resp.Body).Decode(res)
}

func (c *Client) NewRequest(method, path string, body interface{}) (*http.Request, error) {
	u, err := c.BaseURL.Parse(path)
	if err != nil {
		return nil, err
	}

	var req *http.Request

	if method == http.MethodGet {
		req, err = http.NewRequest(method, u.String(), nil)
		if err != nil {
			return nil, err
		}
	} else {
		buf := new(bytes.Buffer)
		if body == nil {
			req, err = http.NewRequest(method, u.String(), nil)
			if err != nil {
				return nil, err
			}
		} else if reader, ok := body.(io.Reader); ok {
			req, err = http.NewRequest(method, u.String(), reader)
			if err != nil {
				return nil, err
			}
		} else {
			if err := json.NewEncoder(buf).Encode(body); err != nil {
				return nil, err
			}

			req, err = http.NewRequest(method, u.String(), buf)
			if err != nil {
				return nil, err
			}
			req.Header.Add("Content-Type", "application/json")
		}
	}

	for k, v := range c.headers {
		req.Header.Add(k, v)
	}

	req.Header.Set("User-Agent", "fingcli")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.accessToken))

	return req, nil
}

var (
	ErrNotFound = errors.New("not found")
)

type Response struct {
	*http.Response
	Files   []*FileInfo `json:"files"`
	Message string      `json:"error"`
}

func (r *Response) Error() string {
	return r.Message
	// return fmt.Sprintf("%v %v: %d %v", r.Response.Request.Method, r.Response.Request.URL, r.Response.StatusCode, r.Message)
}

func (r *Response) IsUploadChangesErr() bool {
	return r.StatusCode == 404 && r.Message == "upload changed files"
}

func CheckResponse(r *Response) error {
	if r.StatusCode >= 200 && r.StatusCode <= 299 {
		return nil
	}

	data, err := ioutil.ReadAll(r.Body)
	if err == nil && len(data) > 0 {
		err := json.Unmarshal(data, r)
		if err != nil {
			r.Message = string(data)
		}
	}

	return r
}

func (c *Client) Do(req *http.Request, v interface{}) (*Response, error) {
	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}

	response := &Response{Response: resp}
	err = CheckResponse(response)
	if err != nil {
		return response, err
	}

	defer func() {
		resp.Body.Close()
	}()

	if v != nil {
		err = json.NewDecoder(resp.Body).Decode(v)
		if err != nil {
			return nil, err
		}
	}

	return response, err
}
