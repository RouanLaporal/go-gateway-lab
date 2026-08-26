package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
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

	http.HandleFunc("/users/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("Handling request for /users")
		w.Header().Set("Content-Type", "application/json")
		jsonEncoder := json.NewEncoder(w)
		jsonEncoder.Encode(data)
		fmt.Println("Response sent for /users : ", data)
	})

	http.HandleFunc("/users/{id}", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("Handling request for /users/{id}")
		id := r.URL.Path[len("/users/"):]
		i, err := strconv.Atoi(id)
		if err != nil || i < 1 || i > len(data) {
			http.Error(w, "User not found", http.StatusNotFound)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		jsonEncoder := json.NewEncoder(w)
		jsonEncoder.Encode(data[i-1])
		fmt.Println("Response sent for /users/{id} : ", data)
	})
	http.ListenAndServe(":3001", nil)
}
