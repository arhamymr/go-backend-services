package db

import (
	"context"
	"fmt"
	"os"

	"github.com/jackc/pgx/v5"
)

type PSQLClient struct {
	DBConn *pgx.Conn
}

func NewConnectPsql() *PSQLClient {
	var err error
	var DBConn *pgx.Conn
	DBConn, err = pgx.Connect(context.Background(), os.Getenv("POSTGRES_URL"))

	fmt.Println("Database connected successfully")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect database %v \n", err)
		os.Exit(1)
	}

	return &PSQLClient{
		DBConn,
	}
}
