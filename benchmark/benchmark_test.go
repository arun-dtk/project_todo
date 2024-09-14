package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"reflect"
	"testing"
	"time"

	"bou.ke/monkey"
	"github.com/stretchr/testify/assert"
)

func TestInvokeApiSuccess(t *testing.T) {
	// Mock time.Now to return a fixed time
	startTime := time.Now()
	monkey.Patch(time.Now, func() time.Time {
		return startTime
	})
	defer monkey.Unpatch(time.Now)

	// Mock http.NewRequest to always return a valid request
	monkey.Patch(http.NewRequest, func(method, url string, body io.Reader) (*http.Request, error) {
		return &http.Request{}, nil
	})
	defer monkey.Unpatch(http.NewRequest)

	// Mock client.Do to return a successful response
	monkey.PatchInstanceMethod(reflect.TypeOf(&http.Client{}), "Do", func(client *http.Client, req *http.Request) (*http.Response, error) {
		return &http.Response{
			StatusCode: http.StatusOK,
			Body:       http.NoBody,
		}, nil
	})
	defer monkey.UnpatchInstanceMethod(reflect.TypeOf(&http.Client{}), "Do")

	resultChan := make(chan BenchmarkResult)
	endpoint := EndPoint{
		EndPoint: "http://localhost:8080/todos/10",
		Method:   "GET",
		Body:     map[string]interface{}{},
		Token:    "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJlbWFpbCI6ImFydW5kYW1vZGFydGtAZ21haWwuY29tIiwidXNlcklkIjoxMn0.aaqeUY4h3DlH5msngfy47PmzxmxiWDci1fC3gIOkpRI",
	}

	go invokeApi(endpoint, resultChan)

	// Verify the result
	result := <-resultChan
	assert.Equal(t, http.StatusOK, result.StatusCode)
	assert.Nil(t, result.Error)
	assert.WithinDuration(t, time.Now(), startTime, 10*time.Millisecond)
}

func TestInvokeApiJsonMarshalError(t *testing.T) {
	// Mock time.Now to return a fixed time
	startTime := time.Now()
	monkey.Patch(time.Now, func() time.Time {
		return startTime
	})
	defer monkey.Unpatch(time.Now)

	// Mock json.Marshal to return an error
	monkey.Patch(json.Marshal, func(v interface{}) ([]byte, error) {
		return nil, errors.New("json marshal error")
	})
	defer monkey.Unpatch(json.Marshal)

	resultChan := make(chan BenchmarkResult)
	endpoint := EndPoint{
		EndPoint: "http://localhost:8080/todos",
		Method:   "POST",
		Body:     map[string]interface{}{}, // Will fail during marshaling
		Token:    "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJlbWFpbCI6ImFydW5kYW1vZGFydGtAZ21haWwuY29tIiwidXNlcklkIjoxMn0.aaqeUY4h3DlH5msngfy47PmzxmxiWDci1fC3gIOkpRI",
	}

	go invokeApi(endpoint, resultChan)

	// Verify the result
	result := <-resultChan
	assert.Equal(t, 500, result.StatusCode) // Expecting internal server error due to JSON marshal error
	assert.EqualError(t, result.Error, "json marshal error")
}

func TestInvokeApi_FailedHttpCall(t *testing.T) {
	// Simulate a server that doesn't respond, causing an HTTP error
	endpoint := EndPoint{
		EndPoint: "http://nonexistent.url",
		Method:   "GET",
		Body:     nil,
		Token:    "test_token",
	}

	resultChan := make(chan BenchmarkResult)

	// Call the invokeApi function
	go invokeApi(endpoint, resultChan)

	// Get the result and assert
	result := <-resultChan
	assert.NotNil(t, result.Error)
	assert.Contains(t, result.Error.Error(), "no such host")
	assert.Equal(t, 500, result.StatusCode)
}

func TestInvokeApi_Error(t *testing.T) {
	// Mock a server that returns a 405 Method Not Allowed error
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusMethodNotAllowed)
	}))
	defer server.Close()

	// Define test endpoint with an invalid method
	endpoint := EndPoint{
		EndPoint: server.URL,
		Method:   "INVALID_METHOD", // Simulating an invalid method
		Body:     nil,
		Token:    "invalid_token",
	}

	resultChan := make(chan BenchmarkResult)

	// Call the invokeApi function
	go invokeApi(endpoint, resultChan)

	// Get the result and assert
	result := <-resultChan

	// Validate the result
	assert.Equal(t, 405, result.StatusCode) // Expecting 405 Method Not Allowed
	assert.Nil(t, result.Error)             // The error should be nil because the request went through, but the response indicates an invalid method
}
func captureOutput(f func()) string {
	// Create a pipe to capture stdout
	r, w, _ := os.Pipe()
	// Save original stdout
	origStdout := os.Stdout
	// Set stdout to the write end of the pipe
	os.Stdout = w

	// Execute the function
	f()

	// Close the writer and restore stdout
	w.Close()
	os.Stdout = origStdout

	// Read the captured output
	var buf bytes.Buffer
	io.Copy(&buf, r)
	r.Close()

	return buf.String()
}

func TestPrintBenchmarkResult_WithoutError(t *testing.T) {
	result := BenchmarkResult{
		Endpoint:   "http://localhost:8080/todos/1",
		Method:     "GET",
		StatusCode: 200,
		Duration:   500 * time.Millisecond,
		Error:      nil,
	}

	// Capture the output of printBenchmarkResult
	output := captureOutput(func() {
		printBenchmarkResult(result)
	})

	expectedOutput := "API GET http://localhost:8080/todos/1 completed with status 200 in time 500ms err <nil>\n"
	assert.Contains(t, output, expectedOutput)
}

func TestPrintBenchmarkResult_WithError(t *testing.T) {
	result := BenchmarkResult{
		Endpoint:   "http://localhost:8080/todos/1",
		Method:     "GET",
		StatusCode: 500,
		Duration:   500 * time.Millisecond,
		Error:      fmt.Errorf("test error"),
	}

	// Capture the output of printBenchmarkResult
	output := captureOutput(func() {
		printBenchmarkResult(result)
	})

	expectedOutput := "API GET http://localhost:8080/todos/1 completed with status 500 in time 500ms err test error\n"
	assert.Contains(t, output, expectedOutput)
}

func TestMainLogic(t *testing.T) {
	// Mock all API calls to return success
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	// Replace the endpoints with test server URLs
	endpoints := []EndPoint{
		{
			EndPoint: server.URL + "/todos/10",
			Method:   "GET",
			Body:     nil,
			Token:    "test_token",
		},
		{
			EndPoint: server.URL + "/todos",
			Method:   "POST",
			Body: map[string]interface{}{
				"title": "Test Todo",
				"list": []map[string]interface{}{
					{"item": "test item", "checked": true},
				},
			},
			Token: "test_token",
		},
	}

	resultChannels := make([]chan BenchmarkResult, len(endpoints))

	for index, endpoint := range endpoints {
		resultChannels[index] = make(chan BenchmarkResult)
		go invokeApi(endpoint, resultChannels[index])
	}

	// Simulate receiving results from all channels
	for range endpoints {
		select {
		case result := <-resultChannels[0]:
			assert.Equal(t, 200, result.StatusCode)
		case result := <-resultChannels[1]:
			assert.Equal(t, 200, result.StatusCode)
		}
	}
}
