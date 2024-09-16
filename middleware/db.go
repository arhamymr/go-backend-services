package middleware

import (
	"database/sql"
	"go-backend-services/db"

	"github.com/labstack/echo/v4"
)

func DBConn(dbSQL *sql.DB, dbRedis *db.RedisClient) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			c.Set("db", dbSQL)
			c.Set("db-redis", dbRedis)
			return next(c)
		}
	}
}
