package utils

import (
	"encoding/json"
	"net/http"
)

// SendJSONResponse sends a JSON response with the specified status code and data
func SendJSONResponse(w http.ResponseWriter, statusCode int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	if data != nil {
		if err := json.NewEncoder(w).Encode(data); err != nil {
			// In a real application, you would log this error
			// For now, we'll just ignore it to keep the code simple
			_ = err
		}
	}
}
