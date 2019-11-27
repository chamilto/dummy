package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/chamilto/dummy/internal/db"
)

func HealthCheck(db *db.DB, w http.ResponseWriter, r *http.Request) {
	_, err := db.Ping().Result()

	ok := true

	if err != nil {
		ok = false
	}

	json.NewEncoder(w).Encode(map[string]bool{"ok": ok})
}
