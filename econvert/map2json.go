package econvert

import (
	"encoding/json"
	"fmt"
)

// JsonToMap Convert json bytes to map
func JsonToMap(jsonBytes []byte) (map[string]string, error) {
	m := make(map[string]string)
	err := json.Unmarshal(jsonBytes, &m)
	if err != nil {
		fmt.Printf("Unmarshal with error: %+v\n", err)
		return nil, err
	}

	for k, v := range m {
		fmt.Printf("%v: %v\n", k, v)
	}

	return m, nil
}

// MapToJson Convert map json bytes
func MapToJson(m map[string]string) ([]byte, error) {
	jsonBytes, err := json.Marshal(m)
	if err != nil {
		fmt.Printf("Marshal with error: %+v\n", err)
		return nil, nil
	}

	return jsonBytes, nil
}
