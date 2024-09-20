package main

import (
	"fmt"
	"net/http"

	"github.com/go-playground/validator"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"

	"go-backend-services/db"
	"go-backend-services/handlers"
	"go-backend-services/helpers"
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

type AppValidator struct {
	Validator *validator.Validate
}

func (cv *AppValidator) Validate(i interface{}) error {
	return cv.Validator.Struct(i)
}

func main() {
	// close when program done
	defer PSQLClient.DBConn.Close()
	defer RedisClient.Rdb.Close()

	// start server
	e := echo.New()
	v := validator.New()
	v.RegisterValidation("custom-pass", helpers.ValidatePassword)
	e.Validator = &helpers.AppValidator{Validator: v}
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

	// Auth
	e.POST("/auth/register", handlers.AuthRegister)
	e.Logger.Fatal(e.Start(":1323"))

}
