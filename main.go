package main

import (
	"fmt"
	"go-restfull-api/handlers"
	"net/http"
)

func main() {
	http.HandleFunc("/", handlers.HelloWorldHandler)

	serverAddress := "localhost:8080"
	fmt.Printf("Server starting at http://%s\n", serverAddress)

	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Printf("Error starting server: %s\n", err)
	}
}
