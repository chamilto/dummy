package handlers

import (
	"encoding/json"
	"github.com/go-redis/redis/v7"
	"net/http"
)

func HealthCheck(db *redis.Client, w http.ResponseWriter, r *http.Request) {
	_, err := db.Ping().Result()

	ok := true

	if err != nil {
		ok = false
	}

	json.NewEncoder(w).Encode(map[string]bool{"ok": ok})
}
