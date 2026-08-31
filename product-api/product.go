package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type Product struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

func main() {
	fmt.Println("Starting product service on port 3002...")
	data := []Product{
		{ID: 1, Name: "Keyboard"},
	}
	http.HandleFunc("/products", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("Received request for /products")
		jsonEncoder := json.NewEncoder(w)
		jsonEncoder.Encode(data)
		fmt.Println("Response sent for /products : ", data)
	})
	http.ListenAndServe(":3002", nil)
}
