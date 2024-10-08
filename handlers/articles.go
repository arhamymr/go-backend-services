package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"go-backend-services/db"
	"go-backend-services/types"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
)

func CreateArticle(c echo.Context) error {
	lock.Lock()
	defer lock.Unlock()

	data := new(types.CreateArticleDTO)
	var response types.Response

	err := c.Bind(data)

	if err != nil {
		response = types.Response{
			Status:  http.StatusBadRequest,
			Data:    struct{}{},
			Message: fmt.Sprintf("Failed to bind: %v", err.Error()),
		}
		return c.JSON(http.StatusBadRequest, response)
	}

	err = c.Validate(data)
	if err != nil {
		response = types.Response{
			Status:  http.StatusBadRequest,
			Data:    struct{}{},
			Message: fmt.Sprintf("Failed to validate: %v", err.Error()),
		}
		return c.JSON(http.StatusBadRequest, response)
	}

	imageJSON, err := json.Marshal(data.Image)

	if err != nil {
		response = types.Response{
			Status:  http.StatusBadRequest,
			Data:    struct{}{},
			Message: fmt.Sprintf("Failed to marshall image: %v", err),
		}
		return c.JSON(http.StatusBadRequest, response)
	}

	query := "INSERT INTO articles (title, content, author, excerpt, image, slug, category_id) VALUES ($1, $2, $3, $4, $5, $6, $7)"

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

	_, err = stmt.Exec(data.Title, data.Content, data.Author, data.Excerpt, string(imageJSON), data.Slug, data.CategoryId)

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

func GetArticle(c echo.Context) error {
	lock.Lock()
	defer lock.Unlock()

	uuid := c.Param("uuid")

	query := "SELECT articles.uuid, articles.created_at, articles.updated_at, articles.title, articles.content, articles.author, articles.image, articles.slug, articles.excerpt, categories.name AS category FROM articles LEFT JOIN categories ON articles.category_id = categories.uuid WHERE articles.uuid = $1"

	dbPsql := c.Get("db").(*sql.DB)
	dbRedis := c.Get("db-redis").(*db.RedisClient)

	var data types.Article

	// check redis before go to database
	result, err := dbRedis.Get(uuid)

	// if any data found in redis
	if err == nil {
		err = json.Unmarshal([]byte(result), &data)

		if err != nil {
			fmt.Println("Failed to unmarshall result continue to go to database")
		}

		if err == nil {
			fmt.Println("Get From Redis")
			response = types.Response{
				Status:  http.StatusOK,
				Data:    data,
				Message: "OK from Redis",
			}
			return c.JSON(http.StatusOK, response)
		}
	}

	stmt, err := dbPsql.Prepare(query)

	if err != nil {
		response = types.Response{
			Status:  http.StatusBadRequest,
			Data:    struct{}{},
			Message: fmt.Sprintf("Failed to prepare this query: %v", err),
		}

		return c.JSON(http.StatusOK, response)
	}

	defer stmt.Close()

	row := stmt.QueryRow(uuid)

	err = row.Scan(&data.Uuid, &data.CreatedAt, &data.UpdatedAt, &data.Title, &data.Content, &data.Author, &data.Image, &data.Slug, &data.Excerpt, &data.Category)

	if err != nil {
		if err == sql.ErrNoRows {
			response = types.Response{
				Status:  http.StatusNotFound,
				Data:    struct{}{},
				Message: fmt.Sprintf("Data not found: %v", err),
			}
			return c.JSON(http.StatusNotFound, response)
		}

		response = types.Response{
			Status:  http.StatusInternalServerError,
			Data:    struct{}{},
			Message: fmt.Sprintf("Failed to scan this row: %v", err),
		}

		return c.JSON(http.StatusInternalServerError, response)
	}

	dataJSON, err := json.Marshal(&data)

	if err != nil {
		fmt.Println("Failed to marshall json data not saved to redis")
	} else {
		dbRedis.SetWithExpired(uuid, string(dataJSON), 60*time.Minute)
	}

	response = types.Response{
		Status:  http.StatusOK,
		Data:    data,
		Message: "OK from Database",
	}

	return c.JSON(http.StatusOK, response)
}

func GetAllArticle(c echo.Context) error {
	lock.Lock()
	defer lock.Unlock()

	query := "SELECT articles.uuid, articles.created_at, articles.updated_at, articles.title, articles.content, articles.author, articles.image, articles.slug, articles.excerpt, categories.name AS category FROM articles LEFT JOIN categories ON articles.category_id = categories.uuid LIMIT $1 OFFSET $2"

	limit := c.QueryParam("limit")
	offset := c.QueryParam("offset")

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

	rows, err := stmt.Query(limit, offset)
	if err != nil {
		response = types.Response{
			Status:  http.StatusBadRequest,
			Data:    []struct{}{},
			Message: fmt.Sprintf("Failed to query: %v", err),
		}

		return c.JSON(http.StatusBadRequest, response)
	}

	defer rows.Close()

	var data []types.Article

	for rows.Next() {
		item := types.Article{}

		var imageJSON string

		err := rows.Scan(&item.Uuid, &item.CreatedAt, &item.UpdatedAt, &item.Title, &item.Content, &item.Author, &imageJSON, &item.Slug, &item.Excerpt, &item.Category)

		if err != nil {
			response = types.Response{
				Status:  http.StatusBadRequest,
				Data:    []struct{}{},
				Message: fmt.Sprintf("Failed to scan: %v", err),
			}
			return c.JSON(http.StatusBadRequest, response)
		}

		err = json.Unmarshal([]byte(imageJSON), &item.Image)

		if err != nil {
			response = types.Response{
				Status:  http.StatusBadRequest,
				Data:    []struct{}{},
				Message: fmt.Sprintf("Failed to unmarshall image: %v", err),
			}
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

func UpdateArticle(c echo.Context) error {

	lock.Lock()
	defer lock.Unlock()

	uuid := c.Param("uuid")

	data := new(types.CreateArticleDTO)
	err := c.Bind(data)

	if err != nil {
		response = types.Response{
			Status:  http.StatusBadRequest,
			Data:    struct{}{},
			Message: err.Error(),
		}
		return c.JSON(http.StatusBadRequest, response)
	}

	query := "UPDATE articles SET"
	params := []interface{}{}
	paramIndex := 1
	updatedData := make(map[string]interface{})

	if data.Title != "" {
		query += fmt.Sprintf(" title = $%d,", paramIndex)
		params = append(params, data.Title)
		updatedData["title"] = data.Title
		paramIndex++
	}

	if data.Content != "" {
		query += fmt.Sprintf(" content = $%d,", paramIndex)
		params = append(params, data.Content)
		updatedData["content"] = data.Content
		paramIndex++
	}

	if data.Author != "" {
		query += fmt.Sprintf(" author = $%d,", paramIndex)
		params = append(params, data.Author)
		updatedData["author"] = data.Author
		paramIndex++
	}

	if data.Excerpt != "" {
		query += fmt.Sprintf(" excerpt = $%d,", paramIndex)
		params = append(params, data.Excerpt)
		updatedData["excerpt"] = data.Excerpt
		paramIndex++
	}

	// if data.Image != "" {
	// 	query += fmt.Sprintf(" thumbnail = $%d,", paramIndex)
	// 	params = append(params, data.Thumbnail)
	// 	updatedData["thumbnail"] = data.Thumbnail
	// 	paramIndex++
	// }

	if data.Slug != "" {
		query += fmt.Sprintf(" slug = $%d,", paramIndex)
		params = append(params, data.Slug)
		updatedData["slug"] = data.Slug
		paramIndex++
	}

	// Remove the trailing comma
	query = query[:len(query)-1]
	query += fmt.Sprintf(" WHERE uuid = $%d", paramIndex)
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

	params = append(params, uuid)
	_, err = stmt.Exec(params...)

	if err != nil {
		response = types.Response{
			Status:  http.StatusBadRequest,
			Data:    struct{}{},
			Message: fmt.Sprintf("Failed update to database: %v", err),
		}
		return c.JSON(http.StatusBadRequest, response)
	}

	// if success update to database, delete cache in redis
	dbRedis := c.Get("db-redis").(*db.RedisClient)
	dbRedis.Del(uuid)

	response = types.Response{
		Status:  http.StatusOK,
		Data:    updatedData,
		Message: "Ok",
	}

	return c.JSON(http.StatusOK, response)
}

func DeleteArticle(c echo.Context) error {
	lock.Lock()
	defer lock.Unlock()

	uuid := c.Param("uuid")

	query := "DELETE FROM articles WHERE uuid = $1"
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

	// if success delete to database, delete cache in redis
	dbRedis := c.Get("db-redis").(*db.RedisClient)
	dbRedis.Del(uuid)

	response = types.Response{
		Status:  http.StatusOK,
		Data:    struct{}{},
		Message: "Data deleted successfully",
	}

	return c.JSON(http.StatusOK, response)
}
