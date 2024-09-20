package handlers

import (
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"fmt"
	"go-backend-services/types"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/lib/pq"
)

func AuthLogin(c echo.Context) error {
	lock.Lock()
	defer lock.Unlock()

	// data := new(types.RegisterDTO)

	// todo soon

	var response types.Response

	response = types.Response{
		Status:  http.StatusOK,
		Data:    struct{}{},
		Message: "Ok",
	}

	return c.JSON(http.StatusOK, response)
}

func AuthRegister(c echo.Context) error {
	lock.Lock()
	defer lock.Unlock()

	data := new(types.RegisterDTO)
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

	err = c.Validate(data)
	if err != nil {
		response = types.Response{
			Status:  http.StatusBadRequest,
			Data:    struct{}{},
			Message: err.Error(),
		}
		return c.JSON(http.StatusBadRequest, response)
	}

	query := "INSERT INTO users (name, email, password) VALUES ($1, $2, $3)"
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

	// hash password before insert to database
	hashedPassword := sha256.Sum256([]byte(data.Password))
	hashedPasswordStr := hex.EncodeToString(hashedPassword[:])

	_, err = stmt.Exec(data.Name, data.Email, hashedPasswordStr)

	if err != nil {
		pqErr := err.(*pq.Error)
		fmt.Println(pqErr.Code)

		switch pqErr.Code {
		case "23505":
			response = types.Response{
				Status:  http.StatusBadRequest,
				Data:    struct{}{},
				Message: fmt.Sprintf("Email already exists: %v", err),
			}
			return c.JSON(http.StatusBadRequest, response)
		default:
			response = types.Response{
				Status:  http.StatusBadRequest,
				Data:    struct{}{},
				Message: fmt.Sprintf("Failed insert to database: %v", err),
			}
			return c.JSON(http.StatusBadRequest, response)
		}
	}

	response = types.Response{
		Status:  http.StatusOK,
		Data:    struct{}{},
		Message: "Ok",
	}

	return c.JSON(http.StatusOK, response)
}
