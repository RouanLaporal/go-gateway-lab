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
			Prefix: "/users",
			Target: &url.URL{
				Scheme: "http",
				Host:   "localhost:3001",
			},
		},
		{
			Prefix: "/products",
			Target: &url.URL{
				Scheme: "http",
				Host:   "localhost:3002",
			},
		},
		{
			Prefix: "/health",
			Target: &url.URL{
				Scheme: "http",
				Host:   "localhost:3003",
			},
		},
	}
	for _, route := range routes {
		http.HandleFunc(route.Prefix, func(responseWriter http.ResponseWriter, request *http.Request) {
			proxyURL := route.Target.ResolveReference(request.URL)
			fmt.Printf("Proxying request to: %s\n", proxyURL.String())
			response, err := http.Get(proxyURL.String())
			if err != nil {
				http.Error(responseWriter, "Error fetching data from target service", http.StatusInternalServerError)
				return
			}
			reader, err := io.ReadAll(response.Body)
			if err != nil {
				http.Error(responseWriter, "Error reading response body from target service", http.StatusInternalServerError)
				return
			}

			defer response.Body.Close()

			responseWriter.Header().Set("Content-Type", "application/json")
			responseWriter.WriteHeader(response.StatusCode)
			fmt.Printf("Response from target service: %d\n", response.StatusCode)
			fmt.Printf("Response body from target service: %s\n", reader)
			responseWriter.Write(reader)
		})
	}
	http.ListenAndServe(":3000", nil)
}
