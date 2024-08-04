package handlers

import (
	"database/sql"
	"fmt"
	"go-backend-services/types"
	"net/http"

	"github.com/labstack/echo/v4"
)

func SaveData(c echo.Context) error {
	data := new(types.Crud)

	var response types.Response

	if err := c.Bind(data); err != nil {
		response = types.Response{
			Status:  http.StatusBadRequest,
			Message: err.Error(),
		}
		return c.JSON(http.StatusBadRequest, response)
	}

	query := "INSERT INTO crud (name, description) VALUES ($1, $2)"
	db := c.Get("db").(*sql.DB)
	_, err := db.Exec(query, data.Name, data.Description)

	if err != nil {
		response = types.Response{
			Status:  http.StatusBadRequest,
			Message: fmt.Sprintf("Failed insert to database: %v", err),
		}
		return c.JSON(http.StatusBadRequest, response)
	}

	response = types.Response{
		Status:  http.StatusOK,
		Data:    data,
		Message: "Ok",
	}
	return c.JSON(http.StatusOK, response)
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
