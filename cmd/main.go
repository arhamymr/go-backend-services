package main

import (
	"fmt"

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

	// custom validator
	v := validator.New()
	v.RegisterValidation("custom-pass", helpers.ValidatePassword)
	e.Validator = &helpers.AppValidator{Validator: v}
	fmt.Print(RedisClient)

	// middleware
	e.Use(middleware.DBConn(PSQLClient.DBConn, RedisClient))
	e.Use(middleware.JWTMiddleware())

	// CRUD
	e.POST("/crud", handlers.SaveData)
	e.GET("/crud", handlers.GetAllData)
	e.GET("/crud/:uuid", handlers.GetData)
	e.PUT("/crud/:uuid", handlers.UpdateData)
	e.DELETE("/crud/:uuid", handlers.DeleteData)

	// Auth
	e.POST("/auth/register", handlers.AuthRegister)
	e.POST("/auth/login", handlers.AuthLogin)

	// Articles
	e.POST("/article", handlers.CreateArticle)
	e.GET("/article", handlers.GetAllArticle)
	e.GET("/article/:uuid", handlers.GetArticle)
	e.PUT("/article/:uuid", handlers.UpdateArticle)
	e.DELETE("/article/:uuid", handlers.DeleteArticle)

	// Token
	e.POST("/generate/global-token", handlers.GlobalToken)

	// Test
	e.POST("/mail/test", handlers.TestMessaging)
	e.Logger.Fatal(e.Start(":1323"))

}
