package main

import (
	"encoding/json"
	"errors"
	"net/url"
	"os"
	"strings"
)

type RouteConfig struct {
	Prefix string `json:"prefix"`
	Target string `json:"target"`
}

func LoadRoutesFromConfig() ([]Route, error) {
	file, err := os.Open("routes.json") //File implements the io.Reader interface cause it has a Read method, so we can use it with json.NewDecoder
	if err != nil {
		return nil, err
	}
	defer file.Close()
	var routesConfig []RouteConfig
	var routes []Route
	var existingPrefixes = make(map[string]bool)

	err = json.NewDecoder(file).Decode(&routesConfig)
	if err != nil {
		return nil, err
	}
	for _, route := range routesConfig {
		if route.Prefix == "" || route.Target == "" {
			return nil, errors.New("Invalid route configuration")
		}
		if !strings.HasPrefix(route.Prefix, "/") || route.Prefix == "/" {
			return nil, errors.New("Invalid prefix")
		}
		targetURL, err := url.Parse(route.Target)
		if err != nil {
			return nil, err
		}
		if targetURL.Scheme != "http" && targetURL.Scheme != "https" {
			return nil, errors.New("Invalid url scheme")
		}
		if targetURL.Host == "" {
			return nil, errors.New("Invalid url host")
		}
		if existingPrefixes[route.Prefix] {
			return nil, errors.New("Duplicate prefix found for route")
		}
		existingPrefixes[route.Prefix] = true
		routes = append(routes, Route{Prefix: route.Prefix, Target: targetURL})
	}
	return routes, nil
}
