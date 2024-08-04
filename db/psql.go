package db

import (
	"database/sql"
	"fmt"
	"os"
	"sync"

	_ "github.com/lib/pq"
)

type PSQLClient struct {
	DBConn *sql.DB
}

var (
	instance *PSQLClient
	once     sync.Once
)

func NewConnectPsql() *PSQLClient {
	once.Do(func() {
		var err error
		var DBConn *sql.DB
		DBConn, err = sql.Open("postgres", os.Getenv("POSTGRES_URL"))

		if err != nil {
			fmt.Fprintf(os.Stderr, "Unable to connect database %v \n", err)
			os.Exit(1)
		}

		fmt.Println("Database connected successfully")

		instance = &PSQLClient{
			DBConn,
		}
	})
	return instance
}
