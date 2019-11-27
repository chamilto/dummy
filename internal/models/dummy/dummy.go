package dummy

import (
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/xeipuuv/gojsonschema"

	"github.com/chamilto/dummy/internal/db"
)

const (
	PATTERNS_HMAP = "patterns"
	NAME_HMAP     = "names"
)

type DummyEndpoint struct {
	Name        string            `json:"name"`
	PathPattern string            `json:"pathPattern"`
	HttpMethod  string            `json:"httpMethod"`
	Body        string            `json:"body"`
	StatusCode  float64           `json:"statusCode"`
	Headers     map[string]string `json:"headers"`
	Delay       float64           `json:"delay"`
}

const dummyEndpointSchema = `
{
        "$schema": "http://json-schema.org/schema#",
        "title": "DummyEndpoint",
        "type": "object",
        "properties": {
            "pathPattern": {"type": "string"},
            "statusCode": {"type": "integer"},
            "body": {"type": "string"},
            "name": {"type": "string"},
            "httpMethod": {"type": "string", "enum": ["GET", "POST", "PUT", "PATCH", "DELETE"]},
            "headers": {"type": "object"},
	    "delay": {"type": "integer", "description": "Delay in milliseconds"}
        },
        "required": ["pathPattern", "body", "statusCode", "name", "httpMethod"],
        "additionalProperties": false
}
`

var DummyEndpointSchemaLoader = gojsonschema.NewStringLoader(dummyEndpointSchema)

func (de *DummyEndpoint) Save(db *db.DB) error {
	// Store name of route keyed by path pattern + method
	fmt.Println(de)
	fieldPattern := strings.Join([]string{de.PathPattern, de.HttpMethod}, ":")
	err := db.HSet(db.BuildKey([]string{PATTERNS_HMAP}), fieldPattern, de.Name).Err()

	if err != nil {
		return err
	}

	// Store dummy endpoint object by name
	marshalled, _ := json.Marshal(de)
	err2 := db.HSet(db.BuildKey([]string{NAME_HMAP}), de.Name, marshalled).Err()

	return err2
}

func (de *DummyEndpoint) PathPatternExists(db *db.DB) bool {
	hm := db.BuildKey([]string{PATTERNS_HMAP})
	exists, _ := db.HExists(hm, strings.Join([]string{de.PathPattern, de.HttpMethod}, ":")).Result()

	return exists
}

func (de *DummyEndpoint) NameExists(db *db.DB) bool {
	hm := db.BuildKey([]string{NAME_HMAP})
	exists, _ := db.HExists(hm, de.Name).Result()

	return exists

}

func (de *DummyEndpoint) IsUnique(db *db.DB) (bool, string) {
	if de.PathPatternExists(db) {
		return false, "pathPattern + httpMethod is not unique."

	}

	if de.NameExists(db) {
		return false, "name is not unique."
	}

	return true, ""

}

// Write the response headers, status code, and body from the DummyEndpoint
func (de *DummyEndpoint) SetResponseData(w http.ResponseWriter) {
	de.setResponseHeaders(w)
	w.WriteHeader(int(de.StatusCode))
	w.Write([]byte(de.Body))

}

func (de *DummyEndpoint) RunDelay() {
	time.Sleep(time.Duration(de.Delay) * time.Millisecond)
}

func (de *DummyEndpoint) setResponseHeaders(w http.ResponseWriter) {
	for k, v := range de.Headers {
		w.Header().Set(k, v)
	}
}

func LoadFromName(db *db.DB, name string) (*DummyEndpoint, error) {
	hm := db.BuildKey([]string{NAME_HMAP})
	v, err := db.HGet(hm, name).Result()

	if v == "" {
		return nil, err
	}

	de := &DummyEndpoint{}
	json.Unmarshal([]byte(v), de)

	return de, err
}

func GetAllDummyEndpoints(db *db.DB) (map[string]string, error) {
	hm := db.BuildKey([]string{NAME_HMAP})
	allEndpoints, err := db.HGetAll(hm).Result()

	return allEndpoints, err
}

func MatchEndpoint(db *db.DB, r *http.Request) (*DummyEndpoint, error) {
	hm := db.BuildKey([]string{PATTERNS_HMAP})

	requestPattern := strings.Join([]string{r.URL.Path, r.Method}, ":")
	allPatterns, err := db.HGetAll(hm).Result()

	for pattern, name := range allPatterns {
		if regexp.MustCompile(pattern).MatchString(requestPattern) {
			return LoadFromName(db, name)
		}

	}

	return nil, err
}
