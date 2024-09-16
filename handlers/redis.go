package handlers

import (
	"encoding/json"
	"fmt"
	"go-backend-services/db"
	"net/http"
	"os"
)

type RedisRequest struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

func RedisHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		RedisGet(w, r)
	case http.MethodPatch:
		RedisPatch(w, r)
	default:
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}

func RedisSet(w http.ResponseWriter, r *http.Request) {
	envVar := os.Getenv("POSTGRES_URL")
	fmt.Fprintf(w, "Hello, World!")
	fmt.Fprintf(w, "Env %s \n", envVar)
	fmt.Fprintf(w, "waw, World!")
}

func RedisGet(w http.ResponseWriter, r *http.Request) {
	redisClient := db.GetRedisClient()

	key := r.URL.Query().Get("key")

	if key == "" {
		http.Error(w, "Missing 'key' paramater", http.StatusBadRequest)
		return
	}

	val, err := redisClient.Get(key)

	if err != nil {
		http.Error(w, "Error getting value from redis", http.StatusInternalServerError)
		return
	}

	w.Write([]byte(val))
}

func RedisPatch(w http.ResponseWriter, r *http.Request) {
	redisClient := db.GetRedisClient()

	var req RedisRequest
	err := json.NewDecoder(r.Body).Decode(&req)

	if err != nil {
		http.Error(w, "Error decoding request body", http.StatusBadRequest)
		return
	}

	err = redisClient.Set(req.Key, req.Value)
	if err != nil {
		http.Error(w, "Failed to set 'value' in redis", http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Value set successfully"))
}
