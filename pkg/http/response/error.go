package response

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type errorResponse struct {
	Message string `json:"message"`
}

func Error(w http.ResponseWriter, message string, status int) {
	res := errorResponse{
		Message: message,
	}

	b, err := json.Marshal(res)
	if err != nil {
		Error(w, fmt.Sprintf("Something went wrong, %s", err.Error()), http.StatusInternalServerError)
	}

	w.WriteHeader(status)
	w.Header().Add("Content-Type", "application/json")
	_, _ = w.Write(b)
}
