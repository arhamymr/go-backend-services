package handlers

import (
	"database/sql"
	"fmt"
	"go-backend-services/types"
	"net/http"

	"github.com/labstack/echo/v4"
)

func CreateCategory(c echo.Context) error {
	lock.Lock()
	defer lock.Unlock()

	data := new(types.CreateCategoryDTO)
	var response types.Response

	err := c.Bind(data)

	if err != nil {
		response = types.Response{
			Status:  http.StatusBadRequest,
			Data:    struct{}{},
			Message: err.Error(),
		}
		return c.JSON(http.StatusBadRequest, response)
	}

	query := "INSERT INTO categories (name) VALUES ($1)"

	db := c.Get("db").(*sql.DB)

	stmt, err := db.Prepare(query)

	if err != nil {
		response = types.Response{
			Status:  http.StatusBadRequest,
			Data:    struct{}{},
			Message: fmt.Sprintf("Failed to prepare this query: %v", err),
		}
		return c.JSON(http.StatusBadRequest, response)
	}

	defer stmt.Close()

	_, err = stmt.Exec(data.Name)

	if err != nil {
		response = types.Response{
			Status:  http.StatusBadRequest,
			Data:    struct{}{},
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

func GetAllCategory(c echo.Context) error {
	lock.Lock()
	defer lock.Unlock()

	query := "SELECT * FROM categories"
	db := c.Get("db").(*sql.DB)

	stmt, err := db.Prepare(query)

	if err != nil {
		response = types.Response{
			Status:  http.StatusBadRequest,
			Data:    struct{}{},
			Message: fmt.Sprintf("Failed to prepare query %v", err),
		}

		return c.JSON(http.StatusOK, response)
	}

	defer stmt.Close()

	rows, err := stmt.Query()
	if err != nil {
		response = types.Response{
			Status:  http.StatusBadRequest,
			Data:    []struct{}{},
			Message: fmt.Sprintf("Failed to query: %v", err),
		}

		return c.JSON(http.StatusBadRequest, response)
	}

	defer rows.Close()

	var data []types.Category

	for rows.Next() {
		item := types.Category{}
		err := rows.Scan(&item.Uuid, &item.Name)

		if err != nil {
			response = types.Response{
				Status:  http.StatusBadRequest,
				Data:    []struct{}{},
				Message: fmt.Sprintf("Failed to scan: %v", err),
			}
			return c.JSON(http.StatusBadRequest, response)
		}
		data = append(data, item)

		if err := rows.Err(); err != nil {
			response = types.Response{
				Status:  http.StatusInternalServerError,
				Data:    []struct{}{},
				Message: fmt.Sprintf("Failed to iterating rows: %v", err),
			}
			return c.JSON(http.StatusInternalServerError, response)
		}
	}

	response = types.Response{
		Status:  http.StatusOK,
		Data:    data,
		Message: "Ok",
	}

	return c.JSON(http.StatusOK, response)
}

func DeleteCategory(c echo.Context) error {
	lock.Lock()
	defer lock.Unlock()

	uuid := c.Param("uuid")

	query := "DELETE FROM categories WHERE uuid = $1"
	dbPsql := c.Get("db").(*sql.DB)

	// Prepare the query
	stmt, err := dbPsql.Prepare(query)

	if err != nil {
		response = types.Response{
			Status:  http.StatusBadRequest,
			Data:    struct{}{},
			Message: fmt.Sprintf("Failed to prepare query: %v", err),
		}

		return c.JSON(http.StatusBadRequest, response)
	}

	defer stmt.Close()

	result, err := stmt.Exec(uuid)

	if err != nil {
		response = types.Response{
			Status:  http.StatusBadRequest,
			Data:    struct{}{},
			Message: fmt.Sprintf("Failed delete to database: %v", err),
		}
		return c.JSON(http.StatusBadRequest, response)
	}

	// check if not found
	rowsAffected, err := result.RowsAffected()

	if err != nil {
		response = types.Response{
			Status:  http.StatusBadRequest,
			Data:    struct{}{},
			Message: fmt.Sprintf("Failed to get rows affected: %v", err),
		}
		return c.JSON(http.StatusBadRequest, response)
	}

	if rowsAffected == 0 {
		response = types.Response{
			Status:  http.StatusNotFound,
			Data:    struct{}{},
			Message: fmt.Sprintf("Data not found: %v", err),
		}
		return c.JSON(http.StatusNotFound, response)
	}

	response = types.Response{
		Status:  http.StatusOK,
		Data:    struct{}{},
		Message: "Data deleted successfully",
	}

	return c.JSON(http.StatusOK, response)
}
