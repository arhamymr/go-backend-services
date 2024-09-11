package handlers

import (
	"database/sql"
	"fmt"
	"go-backend-services/types"
	"net/http"
	"sync"

	"github.com/labstack/echo/v4"
)

// you can using syncRWMutex for better performance
var lock sync.Mutex

func SaveData(c echo.Context) error {
	lock.Lock()
	defer lock.Unlock()

	data := new(types.Crud)

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

	query := "INSERT INTO crud ( name, description) VALUES ($1, $2)"
	db := c.Get("db").(*sql.DB)
	_, err = db.Exec(query, data.Name, data.Description)

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

func GetData(c echo.Context) error {
	lock.Lock()
	defer lock.Unlock()

	id := c.Param("id")
	fmt.Print("id:", id)

	query := "SELECT id, name, description FROM crud WHERE id = $1"
	db := c.Get("db").(*sql.DB)
	row := db.QueryRow(query, id)

	var item types.Crud

	err := row.Scan(&item.Id, &item.Name, &item.Description)

	var response types.Response

	if err != nil {
		if err == sql.ErrNoRows {
			response = types.Response{
				Status:  http.StatusNotFound,
				Data:    struct{}{},
				Message: fmt.Sprintf("Data not found: %v", err),
			}
			return c.JSON(http.StatusNotFound, response)
		}
	}

	response = types.Response{
		Status:  http.StatusOK,
		Data:    item,
		Message: "Ok",
	}
	return c.JSON(http.StatusOK, response)
}

func GetAllData(c echo.Context) error {
	lock.Lock()
	defer lock.Unlock()

	query := "SELECT id, name, description FROM crud"
	db := c.Get("db").(*sql.DB)

	rows, err := db.Query(query)

	var response types.Response

	if err != nil {

		response = types.Response{
			Status:  http.StatusBadRequest,
			Data:    []struct{}{},
			Message: fmt.Sprintf("Failed to query: %v", err),
		}

		return c.JSON(http.StatusBadRequest, response)
	}

	defer rows.Close()

	var data []types.Crud

	for rows.Next() {
		item := types.Crud{}
		err := rows.Scan(&item.Id, &item.Name, &item.Description)

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

func UpdateData(c echo.Context) error {

	lock.Lock()
	defer lock.Unlock()

	id := c.Param("id")
	var response types.Response

	data := new(types.Crud)
	err := c.Bind(data)

	if err != nil {
		response = types.Response{
			Status:  http.StatusBadRequest,
			Data:    struct{}{},
			Message: err.Error(),
		}
		return c.JSON(http.StatusBadRequest, response)
	}

	query := "UPDATE crud SET"
	params := []interface{}{}
	paramIndex := 1
	updatedData := make(map[string]interface{})

	if data.Name != "" {
		query += fmt.Sprintf(" name = $%d,", paramIndex)
		params = append(params, data.Name)
		updatedData["name"] = data.Name
		paramIndex++
	}

	if data.Description != "" {
		query += fmt.Sprintf(" description = $%d,", paramIndex)
		params = append(params, data.Description)
		updatedData["description"] = data.Description
		paramIndex++
	}

	// Remove the trailing comma
	query = query[:len(query)-1]
	query += fmt.Sprintf(" WHERE id = $%d", paramIndex)
	db := c.Get("db").(*sql.DB)

	// Prepare the query
	stmt, err := db.Prepare(query)

	if err != nil {
		response = types.Response{
			Status:  http.StatusBadRequest,
			Data:    struct{}{},
			Message: fmt.Sprintf("Failed to prepare query: %v", err),
		}

		return c.JSON(http.StatusBadRequest, response)
	}

	defer stmt.Close()

	params = append(params, id)
	_, err = stmt.Exec(params...)

	if err != nil {
		response = types.Response{
			Status:  http.StatusBadRequest,
			Data:    struct{}{},
			Message: fmt.Sprintf("Failed update to database: %v", err),
		}
		return c.JSON(http.StatusBadRequest, response)
	}

	response = types.Response{
		Status:  http.StatusOK,
		Data:    updatedData,
		Message: "Ok",
	}

	return c.JSON(http.StatusOK, response)
}

func DeleteData(c echo.Context) error {
	lock.Lock()
	defer lock.Unlock()

	id := c.Param("id")

	query := "DELETE FROM crud WHERE id = $1"
	db := c.Get("db").(*sql.DB)

	// Prepare the query
	stmt, err := db.Prepare(query)
	var response types.Response

	if err != nil {
		response = types.Response{
			Status:  http.StatusBadRequest,
			Data:    struct{}{},
			Message: fmt.Sprintf("Failed to prepare query: %v", err),
		}

		return c.JSON(http.StatusBadRequest, response)
	}

	defer stmt.Close()

	result, err := stmt.Exec(id)

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
