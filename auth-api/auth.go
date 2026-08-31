package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type HealthResponse struct {
	Status string `json:"status"`
}

func main() {
	fmt.Println("Starting health service on port 3003...")
	data := HealthResponse{
		Status: "success",
	}
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("Handling request for /health")
		jsonEncoder := json.NewEncoder(w)
		jsonEncoder.Encode(data)
		fmt.Println("Response sent for /health : ", data)
	})
	http.ListenAndServe(":3003", nil)
}
