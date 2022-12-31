package discovery

import (
	"encoding/json"
	"errors"
)

var InvalidJSONErr = errors.New("invalid JSON was returned")

func ParseMetaData(data string) (map[string]any, error) {
	if !json.Valid([]byte(data)) {
		return map[string]any{}, InvalidJSONErr
	}

	var parsedData map[string]any
	json.Unmarshal([]byte(data), &parsedData)
	return parsedData, nil
}
