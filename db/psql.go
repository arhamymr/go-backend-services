package db

import (
	"context"
	"fmt"
	"os"
	"sync"

	"github.com/jackc/pgx/v5"
)

type PSQLClient struct {
	DBConn *pgx.Conn
}

var (
	instance *PSQLClient
	once     sync.Once
)

func NewConnectPsql() *PSQLClient {
	once.Do(func() {
		var err error
		var DBConn *pgx.Conn
		DBConn, err = pgx.Connect(context.Background(), os.Getenv("POSTGRES_URL"))

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
