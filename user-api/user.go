package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type User struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

func main() {
	fmt.Println("Starting user service on port 3001...")
	data := []User{
		{ID: 1, Name: "Alice"},
		{ID: 2, Name: "Bob"},
	}
	http.HandleFunc("/users", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("Handling request for /users")
		jsonEncoder := json.NewEncoder(w)
		jsonEncoder.Encode(data)
		fmt.Println("Response sent for /users : ", data)
	})
	http.ListenAndServe(":3001", nil)
}
