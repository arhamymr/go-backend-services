package handlers

import (
	"fmt"
	"net/http"
	"os"
)

func HelloWorldHandler(w http.ResponseWriter, r *http.Request) {
	envVar := os.Getenv("POSTGRES_URL")
	fmt.Fprintf(w, "Hello, World!")
	fmt.Fprintf(w, "Env %s \n", envVar)
	fmt.Fprintf(w, "waw, World!")
}
