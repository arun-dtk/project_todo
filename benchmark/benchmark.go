package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type BenchmarkResult struct {
	Endpoint   string
	Method     string
	StatusCode int
	Duration   time.Duration
	Error      error
}

type EndPoint struct {
	EndPoint string
	Method   string
	Body     interface{}
	Token    string
}

func main() {
	token := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJlbWFpbCI6ImFydW5kYW1vZGFydGtAZ21haWwuY29tIiwidXNlcklkIjoxMn0.aaqeUY4h3DlH5msngfy47PmzxmxiWDci1fC3gIOkpRI"

	endpoints := []EndPoint{
		EndPoint{
			EndPoint: "http://localhost:8080/todos/10",
			Method:   "GET",
			Body:     map[string]interface{}{},
			Token:    token,
		},
		EndPoint{
			EndPoint: "http://localhost:8080/todos",
			Method:   "GET",
			Body:     map[string]interface{}{},
			Token:    token,
		},
		EndPoint{
			EndPoint: "http://localhost:8080/todos/11",
			Method:   "PUT",
			Body: map[string]interface{}{
				"title": "Updated Todo",
				"list": []map[string]interface{}{
					{"item": "updated item1", "checked": true},
				},
			},
			Token: token,
		},
		EndPoint{
			EndPoint: "http://localhost:8080/todos",
			Method:   "POST",
			Body: map[string]interface{}{
				"title": "New Todo",
				"list": []map[string]interface{}{
					{"item": "checklist item1", "checked": true},
				},
			},
			Token: token,
		},
	}
	resultChannels := make([]chan BenchmarkResult, len(endpoints))

	for index, endpoint := range endpoints {
		resultChannels[index] = make(chan BenchmarkResult)
		go invokeApi(endpoint, resultChannels[index])
	}

	// Use select to listen for results from any of the channels
	for range endpoints {
		select {
		case result := <-resultChannels[0]:
			printBenchmarkResult(result)
		case result := <-resultChannels[1]:
			printBenchmarkResult(result)
		case result := <-resultChannels[2]:
			printBenchmarkResult(result)
		case result := <-resultChannels[3]:
			printBenchmarkResult(result)
		}
	}
}

func printBenchmarkResult(result BenchmarkResult) {
	fmt.Printf("API %s %s completed with status %d in time %v err %v\n", result.Method, result.Endpoint, result.StatusCode, result.Duration, result.Error)
}

func invokeApi(endpoint EndPoint, resultChan chan BenchmarkResult) {
	start := time.Now()

	// uncomment to check concurrent execution of invokeApi.
	// if endpoint.Method == "GET" {
	// 	time.Sleep(4 * time.Second)
	// }

	// Marshal the body to JSON if necessary
	jsonBody, err := json.Marshal(endpoint.Body)
	if err != nil {
		resultChan <- BenchmarkResult{Endpoint: endpoint.EndPoint, Method: endpoint.Method, StatusCode: 500, Duration: 0, Error: err}
		return // Exit early if there's an error
	}

	// Create the request
	req, err := http.NewRequest(endpoint.Method, endpoint.EndPoint, bytes.NewBuffer(jsonBody))

	if err != nil {
		resultChan <- BenchmarkResult{Endpoint: endpoint.EndPoint, Method: endpoint.Method, StatusCode: 500, Duration: 0, Error: err}
		return // Exit early if there's an error
	}

	// Set content type and authorization headers if necessary
	if len(jsonBody) > 0 {
		req.Header.Set("Content-Type", "application/json")
	}
	if endpoint.Token != "" {
		req.Header.Set("Authorization", endpoint.Token)
	}

	// Execute the request
	client := &http.Client{}
	resp, err := client.Do(req)
	duration := time.Since(start)

	if err != nil {
		resultChan <- BenchmarkResult{Endpoint: endpoint.EndPoint, Method: endpoint.Method, StatusCode: 500, Duration: duration, Error: err}
		return // Exit early if there's an error
	}
	defer resp.Body.Close()
	// Send the result to the channel
	resultChan <- BenchmarkResult{Endpoint: endpoint.EndPoint, Method: endpoint.Method, StatusCode: resp.StatusCode, Duration: duration, Error: nil}
}
