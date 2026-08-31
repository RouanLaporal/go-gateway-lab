package main

import (
	"fmt"
	"net/http"
	"net/http/httputil"
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
		fmt.Println("Proxying request for ", route.Prefix, " to ", route.Target)
		proxy := httputil.NewSingleHostReverseProxy(route.Target)
		http.Handle(route.Prefix, proxy)
		http.Handle(route.Prefix+"/", proxy)
		http.HandleFunc(route.Prefix+"/{$}", func(w http.ResponseWriter, r *http.Request) {
			http.Error(w, "404 Not Found", http.StatusNotFound)
		})
	}

	http.ListenAndServe(":3000", nil)
}
