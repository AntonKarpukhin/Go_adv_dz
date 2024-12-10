package request

import (
	"encoding/json"
	"net/http"
)

func BodyDecode[T any](w *http.ResponseWriter, r *http.Request) (*T, error) {
	var payload T
	err := json.NewDecoder(r.Body).Decode(&payload)
	if err != nil {
		return nil, err
	}
	return &payload, nil
}
