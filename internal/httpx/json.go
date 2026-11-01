package httpx

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// omitempty is used to omit empty fields in the JSON response. If a field has a zero value, it will not be included in the JSON output.
type Response struct {
	Status bool `json:"status"`
	Message string `json:"message,omitempty"`
	Data any `json:"data,omitempty"`
}

// JSON writes the given data as a JSON response with the specified status code.
func WriteJSONResponse(w http.ResponseWriter, statusCode int, data Response) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	err := json.NewEncoder(w).Encode(data)
	if err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		fmt.Println(err)
	}
}