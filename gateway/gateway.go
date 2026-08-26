package main

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
)

type Route struct {
	Prefix string
	Target *url.URL
}

func main() {
	var routes = []Route{
		{
			Prefix: "/users/",
			Target: &url.URL{
				Scheme: "http",
				Host:   "localhost:3001",
			},
		},
		{
			Prefix: "/products/",
			Target: &url.URL{
				Scheme: "http",
				Host:   "localhost:3002",
			},
		},
		{
			Prefix: "/health/",
			Target: &url.URL{
				Scheme: "http",
				Host:   "localhost:3003",
			},
		},
	}
	for _, route := range routes {
		http.HandleFunc(route.Prefix, func(responseWriter http.ResponseWriter, request *http.Request) {
			serviceResponse, err := proxyRequest(route.Target, request)
			if err != nil {
				http.Error(responseWriter, "Error fetching data from target service", http.StatusInternalServerError)
				return
			}

			defer serviceResponse.Body.Close()
			proxyResponse(responseWriter, serviceResponse)
		})
	}
	http.ListenAndServe(":3000", nil)
}

func proxyRequest(target *url.URL, request *http.Request) (*http.Response, error) {
	proxyURL := target.ResolveReference(request.URL)
	fmt.Printf("Proxying request to: %s\n", proxyURL.String())
	client := &http.Client{}
	newRequest, err := http.NewRequest("GET", proxyURL.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("error creating request to target service: %w", err)
	}

	for key, values := range request.Header {
		for _, value := range values {
			fmt.Printf("Adding Header %s with value %s to proxy request\n", key, value)
			newRequest.Header.Add(key, value)
		}
	}

	newRequest.URL.RawQuery = request.URL.Query().Encode()
	fmt.Printf("Proxying request to: %s\n", newRequest.URL.String())
	response, err := client.Do(newRequest)
	if err != nil {
		return nil, fmt.Errorf("error fetching data from target service: %w", err)
	}
	return response, nil
}

func proxyResponse(responseWriter http.ResponseWriter, response *http.Response) {

	proxyResponse := &http.Response{
		StatusCode: response.StatusCode,
		Body:       response.Body,
		Header:     response.Header,
	}

	for key, values := range proxyResponse.Header {
		for _, value := range values {
			fmt.Printf("Adding Header %s with value %s from target service\n", key, value)
			responseWriter.Header().Add(key, value)
		}
	}
	responseWriter.WriteHeader(proxyResponse.StatusCode)
	fmt.Printf("Response from target service: %d\n", proxyResponse.StatusCode)
	fmt.Printf("Response headers from target service: %v\n", proxyResponse.Header)
	fmt.Printf("Response body from target service: %s\n", proxyResponse.Body)

	reader, err := io.ReadAll(proxyResponse.Body)
	if err != nil {
		http.Error(responseWriter, "Error reading response body from target service", http.StatusInternalServerError)
		return
	}
	responseWriter.Write(reader)
}
