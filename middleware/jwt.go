package middleware

import (
	"go-backend-services/helpers"
	"go-backend-services/types"
	"net/http"
	"regexp"
	"strings"

	"github.com/labstack/echo/v4"
)

var response types.Response

func getToken(c echo.Context) string {
	auth := c.Request().Header.Get("Authorization")
	token := strings.TrimPrefix(auth, "Bearer ")
	return token
}

func JWTMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {

			if c.Path() == "/generate/global-token" {
				return next(c)
			}

			pattern := regexp.MustCompile(`^/crud.*|^/generate/global-token$|^/article.*|^/auth/login$|^/mail/test$|^/auth/register$`)

			if pattern.MatchString(c.Path()) {
				token := getToken(c)

				if token == "" {
					response = types.Response{
						Status:  http.StatusUnauthorized,
						Data:    struct{}{},
						Message: "Token is required",
					}
					return c.JSON(http.StatusUnauthorized, response)
				}

				err := helpers.VerifyToken(token, "SECRET_GLOBAL_TOKEN_KEY")

				if err != nil {
					response = types.Response{
						Status:  http.StatusUnauthorized,
						Data:    struct{}{},
						Message: "Invalid token",
					}
					return c.JSON(http.StatusUnauthorized, response)
				}

				return next(c)
			}

			token := getToken(c)

			if token == "" {
				response = types.Response{
					Status:  http.StatusUnauthorized,
					Data:    struct{}{},
					Message: "Token is required",
				}
				return c.JSON(http.StatusUnauthorized, response)
			}

			err := helpers.VerifyToken(token, "SECRET_TOKEN_KEY")

			if err != nil {
				response = types.Response{
					Status:  http.StatusUnauthorized,
					Data:    struct{}{},
					Message: "Invalid token",
				}
				return c.JSON(http.StatusUnauthorized, response)
			}

			return next(c)
		}
	}
}
