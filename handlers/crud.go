package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func SaveData(c echo.Context) error {
	id := c.Param("id")
	return c.String(http.StatusOK, id)
}

func GetData(c echo.Context) error {
	id := c.Param("id")
	return c.String(http.StatusOK, id)
}

func UpdateData(c echo.Context) error {
	id := c.Param("id")
	return c.String(http.StatusOK, id)
}

func DeleteData(c echo.Context) error {
	id := c.Param("id")
	return c.String(http.StatusOK, id)
}
