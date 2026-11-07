package httpx

import (
	"net/http"
)

func WriteErrorResponse(w http.ResponseWriter, statusCode int, message string) {
	WriteJSONResponse(
		w,
		statusCode,
		Response{
			Status: false,
			Message: message,
		},
	)
}