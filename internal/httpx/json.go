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

func WriteMethodNotAllowedResponse(w http.ResponseWriter) {
	WriteJSONResponse(w, http.StatusMethodNotAllowed, Response{
		Status: false,
		Message: "Method not allowed",
	})
}

func WriteCustomErr(w http.ResponseWriter, err error) {}

func WriteBadRequest(w http.ResponseWriter) {
	WriteJSONResponse(w, http.StatusBadRequest, Response{
		Status: false,
		Message: "cannot decode json",
	})
}

func WriteInternalServerErr(w http.ResponseWriter) {
	WriteJSONResponse(w, http.StatusInternalServerError, Response{
		Status: false,
		Message: "internal server error",
	})
}
