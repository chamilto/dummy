package api

import (
	"encoding/json"
	"github.com/go-redis/redis/v7"
	"net/http"
)

type DummyEndpoint struct {
	pathPattern string
	statusCode  int64
	body        string
	name        string
	httpMethod  string
	headers     struct{}
}

func (*DummyEndpoint) save() {

}

func MatchEndpoint(pathPattern string, db *redis.Client) (*DummyEndpoint, error) {
	// get all patterns
	// if pattern in patterns
	// if re.search match, return
	return &DummyEndpoint{}, nil
}

func healthCheckHandler(db *redis.Client, w http.ResponseWriter, r *http.Request) {
	_, err := db.Ping().Result()

	ok := true

	if err != nil {
		ok = false
	}

	json.NewEncoder(w).Encode(map[string]bool{"ok": ok})
}
