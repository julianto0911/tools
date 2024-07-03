package tools

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

const (
	MockUrl = "http://127.0.0.1:80"
)

type BrokenReader struct{}

func (br *BrokenReader) Read(p []byte) (n int, err error) {
	return 0, fmt.Errorf("failed reading")
}

func (br *BrokenReader) Close() error {
	return fmt.Errorf("failed closing")
}

func TestRequestWithJSON(t *testing.T) {
	//mock response
	mockResp := map[string]interface{}{
		"code":    0,
		"message": "success",
		"balance": 1000,
	}

	//init logger
	logger, _ := MockLogs()

	//mock handler
	mockHandler := func(req *http.Request) *http.Response {
		//assert route
		//check , if url is "faulty then return init request error"
		if req.URL.String() == "faulty" {
			//return faulty
			return nil
		}
		if req.URL.String() == "badbody" {
			reader := BrokenReader{}
			//return bad reading response
			return &http.Response{
				ContentLength: 1,
				StatusCode:    500,
				Body:          &reader,
				Header:        make(http.Header),
			}
		}

		//assert header
		//assert.Equal(t, "application/json", req.Header.Get("Content-Type"))

		//get response body
		var input map[string]interface{}
		err := json.NewDecoder(req.Body).Decode(&input)
		if err != nil {
			return &http.Response{
				StatusCode: 500,
				Body:       io.NopCloser(bytes.NewBufferString(err.Error())),
				Header:     make(http.Header),
			}
		}

		//by checking if token is "1" return positive login
		return &http.Response{
			StatusCode: 200,
			Body:       io.NopCloser(bytes.NewBufferString(ToJSONString(mockResp))),
			Header:     make(http.Header),
		}

	}

	//mock client and net connection
	client := MockClient(mockHandler)
	net := NewRequestAdaptor(client.Transport, 1*time.Second, logger, true)

	//mock success
	mockInput := map[string]string{
		"member_code": "1",
		"token":       "1",
	}

	//mock success with json same of mock response
	mockInput["token"] = "1"
	b, err := net.RequestWithJSON("POST", MockUrl, mockInput, nil)
	assert.Nil(t, err, "error should nil")
	assert.Equal(t, b, ToBytes(mockResp), "object should be same")

	//test with json with fault input
	type Dummy struct {
		Name string
		Next *Dummy
	}
	dummy := Dummy{Name: "Dummy"}
	dummy.Next = &dummy
	_, err = net.RequestWithJSON("POST", MockUrl, dummy, nil)
	if assert.NotNil(t, err, "error should exist") {
		assert.Contains(t, err.Error(), "marshal", "error should report marshaling error")
	}

	//test with json url parse error
	_, err = net.RequestWithJSON("POST", "postgres://user:abc{DEf1=ghi@example.com:5432/db?sslmode=require", nil, nil)
	if assert.NotNil(t, err, "error should exist") {
		assert.Contains(t, err.Error(), "parse", "error should report parse url error")
	}

	//test with json init request error
	_, err = net.RequestWithJSON("bad method", "faulty", nil, nil)
	if assert.NotNil(t, err, "error should exist") {
		assert.Contains(t, err.Error(), "init request", "error should report init request error")
	}

	//test with json do request error
	_, err = net.RequestWithJSON("POST", "faulty", nil, nil)
	if assert.NotNil(t, err, "error should exist") {
		assert.Contains(t, err.Error(), "do request", "error should report init request error")
	}

	//test with json do request error
	_, err = net.RequestWithJSON("POST", "badbody", nil, nil)
	if assert.NotNil(t, err, "error should exist") {
		assert.Contains(t, err.Error(), "read body", "error should report init request error")
	}

}

func TestRequestWithQuery(t *testing.T) {
	//mock response
	mockResp := map[string]interface{}{
		"code":    0,
		"message": "success",
		"balance": 1000,
	}

	//init logger
	logger, _ := MockLogs()

	//mock handler
	mockHandler := func(req *http.Request) *http.Response {
		//assert route
		//check , if url is "faulty then return init request error"
		if req.URL.String() == "faulty" {
			//return faulty
			return nil
		}
		if req.URL.String() == "badbody" {
			reader := BrokenReader{}
			//return bad reading response
			return &http.Response{
				ContentLength: 1,
				StatusCode:    500,
				Body:          &reader,
				Header:        make(http.Header),
			}
		}

		//assert header
		//assert.Equal(t, "application/json", req.Header.Get("Content-Type"))

		//by checking if token is "1" return positive login
		return &http.Response{
			StatusCode: 200,
			Body:       io.NopCloser(bytes.NewBufferString(ToJSONString(mockResp))),
			Header:     make(http.Header),
		}

	}

	//mock client and net connection
	client := MockClient(mockHandler)
	net := NewRequestAdaptor(client.Transport, 1*time.Second, logger, true)

	//test with query for success request
	b, err := net.RequestWithQuery("GET", "https://www.google.com", nil, nil)
	assert.Nil(t, err, "error should nil")
	assert.Equal(t, b, ToBytes(mockResp), "object should be same")

	//test with json url parse error
	_, err = net.RequestWithQuery("POST", "postgres://user:abc{DEf1=ghi@example.com:5432/db?sslmode=require", nil, nil)
	if assert.NotNil(t, err, "error should exist") {
		assert.Contains(t, err.Error(), "parse", "error should report parse url error")
	}

	//test with query init request error
	_, err = net.RequestWithQuery("bad method", "faulty", nil, nil)
	if assert.NotNil(t, err, "error should exist") {
		assert.Contains(t, err.Error(), "init request", "error should report init request error")
	}

	//test with query do request error
	_, err = net.RequestWithQuery("POST", "faulty", nil, nil)
	if assert.NotNil(t, err, "error should exist") {
		assert.Contains(t, err.Error(), "do request", "error should report init request error")
	}

	//test with json do request error
	_, err = net.RequestWithQuery("POST", "badbody", nil, nil)
	if assert.NotNil(t, err, "error should exist") {
		assert.Contains(t, err.Error(), "read body", "error should report init request error")
	}
}

func TestMockClient(t *testing.T) {
	//mock handler
	mockHandler := func(req *http.Request) *http.Response {
		//assert route
		assert.Equal(t, MockUrl, req.URL.String())

		//assert method
		assert.Equal(t, req.Method, http.MethodGet)

		//get response body
		return &http.Response{
			StatusCode: 200,
			Body:       io.NopCloser(bytes.NewBufferString(ToJSONString(""))),
			Header:     make(http.Header),
		}
	}

	//mock client and net connection
	client := MockClient(mockHandler)
	_, err := client.Get(MockUrl)
	assert.Nil(t, err)
}
