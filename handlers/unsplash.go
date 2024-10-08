package handlers

import (
	"encoding/json"
	"go-backend-services/types"
	"io"
	"net/http"
	"os"

	"github.com/labstack/echo/v4"
)

func SearchUnplash(c echo.Context) error {
	lock.Lock()
	defer lock.Unlock()

	keyword := c.Param("keyword")
	accessKey := os.Getenv("UNSPLASH_ACCESS_KEY")

	url := "https://api.unsplash.com/search/photos?query=" + keyword + "&client_id=" + accessKey

	resp, err := http.Get(url)

	if err != nil {
		response = types.Response{
			Status:  http.StatusBadRequest,
			Data:    struct{}{},
			Message: err.Error(),
		}
		return c.JSON(http.StatusBadRequest, response)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)

	if err != nil {
		response = types.Response{
			Status:  http.StatusBadRequest,
			Data:    struct{}{},
			Message: err.Error(),
		}
		return c.JSON(http.StatusInternalServerError, response)
	}

	var result types.UnsplashResponse

	err = json.Unmarshal(body, &result)

	if err != nil {
		response = types.Response{
			Status:  http.StatusInternalServerError,
			Data:    struct{}{},
			Message: err.Error(),
		}
	}

	response = types.Response{
		Status:  http.StatusOK,
		Data:    result,
		Message: "Ok keyword: " + keyword,
	}

	return c.JSON(http.StatusOK, response)
}
