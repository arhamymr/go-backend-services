package main

import (
	"fmt"
	"net/http"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"

	"go-backend-services/db"
	"go-backend-services/handlers"
	"go-backend-services/middleware"
)

var (
	PSQLClient  *db.PSQLClient
	RedisClient *db.RedisClient
)

func init() {
	// Load the .env file in the current directory
	if err := godotenv.Load(); err != nil {
		fmt.Println("No .env file found")
	}

	PSQLClient = db.NewConnectPsql()

	// init redis client
	db.InitRedisClient()
	RedisClient = db.GetRedisClient()
}

func main() {
	// close when program done
	defer PSQLClient.DBConn.Close()

	// start server
	e := echo.New()
	fmt.Print(RedisClient)
	e.Use(middleware.DBConn(PSQLClient.DBConn, RedisClient))

	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})

	// CRUD
	e.POST("/crud", handlers.SaveData)
	e.GET("/crud", handlers.GetAllData)
	e.GET("/crud/:uuid", handlers.GetData)
	e.PUT("/crud/:uuid", handlers.UpdateData)
	e.DELETE("/crud/:uuid", handlers.DeleteData)

	e.Logger.Fatal(e.Start(":1323"))

}
