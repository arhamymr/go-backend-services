package middleware

import (
	"github.com/go-playground/validator"
	"github.com/labstack/echo/v4"
)

func Validator(val *validator.Validate) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			c.Set("validator", val)
			return next(c)
		}
	}
}
