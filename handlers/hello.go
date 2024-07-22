package handlers

import (
	"fmt"
	"net/http"
	"os"
)

func HelloWorldHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello, World!")
	fmt.Fprintf(w, os.Getenv(""))
}
