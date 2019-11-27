package utils

import (
	"fmt"

	"github.com/xeipuuv/gojsonschema"
)

func ValidateJson(data []byte, schema gojsonschema.JSONLoader) (bool, []string) {
	valid := true
	dataDoc := gojsonschema.NewStringLoader(string(data))
	result, err := gojsonschema.Validate(schema, dataDoc)

	if err != nil {
		return false, nil
	}

	errs := []string{}

	if !result.Valid() {
		valid = false
		for _, desc := range result.Errors() {
			errs = append(errs, fmt.Sprintf("%s", desc))
		}
	}

	return valid, errs
}
