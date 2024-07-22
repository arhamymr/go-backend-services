package main

import (
	"context"
	"fmt"
	"net/http"

	"github.com/joho/godotenv"

	"go-restfull-api/db"
	"go-restfull-api/handlers"
)

func init() {
	// Load the .env file in the current directory
	if err := godotenv.Load(); err != nil {
		fmt.Println("No .env file found")
	}
}

func main() {
	http.HandleFunc("/", handlers.HelloWorldHandler)

	PSQLClient := db.NewConnectPsql()
	// close when program done
	defer PSQLClient.DBConn.Close(context.Background())

	RdsClient := db.NewRedisClient()
	err := RdsClient.Rdb.Set(context.Background(), "key", "value", 0).Err()
	if err != nil {
		panic(err)
	}
	val, err := RdsClient.Rdb.Get(context.Background(), "key").Result()
	fmt.Println("key", val)

	if err != nil {
		panic(err)
	}

	serverAddress := "localhost:8080"
	fmt.Printf("Server starting at http://%s\n", serverAddress)

	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Printf("Error starting server: %s\n", err)
	}
}
