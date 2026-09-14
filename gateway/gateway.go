package main

import (
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
)

type Route struct {
	Prefix string
	Target *url.URL
}

func main() {
	routes, err := LoadRoutesFromConfig()
	if err != nil {
		log.Fatal(err)
	}

	gateway := NewGateway(routes)
	gateway.Start(":3000")
}

type Gateway struct {
	handler http.Handler
}

func NewGateway(routes []Route) *Gateway {
	mux := http.NewServeMux()
	gateway := &Gateway{}
	for _, route := range routes {
		proxy := httputil.NewSingleHostReverseProxy(route.Target)
		mux.Handle(route.Prefix, proxy)
		mux.Handle(route.Prefix+"/", proxy)
		mux.HandleFunc(route.Prefix+"/{$}", func(w http.ResponseWriter, r *http.Request) {
			http.Error(w, "404 Not Found", http.StatusNotFound)
		})
	}
	gateway.handler = mux
	return gateway
}

func (g *Gateway) Start(ports string) {
	log.Fatal(http.ListenAndServe(ports, g))
}

func (g *Gateway) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	g.handler.ServeHTTP(w, r)
}
