package jsonutil

import (
	"encoding/json"
	"net/http"
)

func JSONError(w http.ResponseWriter, message string, code int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)

	errorResponse := struct {
		Error  string `json:"error"`
		Code   int    `json:"code"`
		Status string `json:"status"`
	}{
		message,
		code,
		http.StatusText(code),
	}

	err := json.NewEncoder(w).Encode(errorResponse)

	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}