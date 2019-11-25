package dummyendpoint

import (
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis/v7"
	"github.com/xeipuuv/gojsonschema"
	"strings"
)

const REDIS_KEY_PREFIX = "DUMMY"

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
            "headers": {"type": "object"}
        },
        "required": ["pathPattern", "body", "statusCode", "name", "httpMethod"],
        "additionalProperties": false
}
`

var DummyEndpointSchemaLoader = gojsonschema.NewStringLoader(dummyEndpointSchema)

func Validate(data []byte, schema gojsonschema.JSONLoader) (bool, []string) {
	valid := true
	dataDoc := gojsonschema.NewStringLoader(string(data))
	result, err := gojsonschema.Validate(schema, dataDoc)

	fmt.Println(err)

	if err != nil {
		return false, nil
	}

	errs := []string{}

	if !result.Valid() {
		valid = false
		fmt.Println(result.Errors())
		for _, desc := range result.Errors() {
			errs = append(errs, fmt.Sprintf("%s", desc))
		}
	}

	return valid, errs
}

func buildKey(parts []string) string {
	// prepend
	parts = append([]string{REDIS_KEY_PREFIX}, parts...)
	return strings.Join(parts, ":")
}

type DummyEndpoint struct {
	Name        string                 `json:"name"`
	PathPattern string                 `json:"pathPattern"`
	HttpMethod  string                 `json:"httpMethod"`
	Body        string                 `json:"body"`
	StatusCode  float64                `json:"statusCode"`
	Headers     map[string]interface{} `json:"headers"`
}

func (de *DummyEndpoint) Save(db *redis.Client) error {
	// Store name of route keyed by path pattern + method
	fmt.Println(de)
	keyPattern := buildKey([]string{de.PathPattern, de.HttpMethod})
	fmt.Println(keyPattern)
	err := db.Set(keyPattern, de.Name, 0).Err()

	if err != nil {
		return err
	}

	// Store dummy endpoint object by name
	keyName := buildKey([]string{de.Name})
	marshalled, _ := json.Marshal(de)
	err2 := db.Set(keyName, marshalled, 0).Err()

	return err2
}

func (de *DummyEndpoint) PathPatternExists(db *redis.Client) bool {
	pathPatternExistsKey := buildKey([]string{de.PathPattern, de.HttpMethod})
	pathPatternExists, _ := db.Exists(pathPatternExistsKey).Result()

	if pathPatternExists == 1 {
		return true
	}

	return false

}

func (de *DummyEndpoint) NameExists(db *redis.Client) bool {
	nameExists, _ := db.Exists(buildKey([]string{de.Name})).Result()

	if nameExists == 1 {
		return true
	}

	return false

}

func (de *DummyEndpoint) IsUnique(db *redis.Client) (bool, string) {
	if de.PathPatternExists(db) {
		return false, "pathPattern + httpMethod is not unique."

	}

	if de.NameExists(db) {
		return false, "name is not unique."
	}

	return true, ""

}

func MatchEndpoint(pathPattern string, db *redis.Client) (*DummyEndpoint, error) {
	// get all patterns
	// if pattern in patterns
	// if re.search match, return
	return &DummyEndpoint{}, nil
}
