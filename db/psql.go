package db

import (
	"context"
	"fmt"
	"os"

	"github.com/jackc/pgx/v5"
)

var DBConn *pgx.Conn

func Connect_psql() {
	var err error
	DBConn, err = pgx.Connect(context.Background(), os.Getenv("POSTGRES_URL"))

	fmt.Println("Database connected successfully")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect database %v \n", err)
		os.Exit(1)
	}
}
