package response

import (
	"encoding/json"
	"net/http"
)

type errorResponse struct {
	Error map[string]interface{} `json:"error"`
}

func ErrorBadRequest(w http.ResponseWriter, message string, fields map[string]string) {
	b, _ := json.Marshal(errorResponse{
		Error: map[string]interface{}{
			"message": message,
			"fields":  fields,
		},
	})
	w.Header().Add("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusBadRequest)
	_, _ = w.Write(b)
}

func ErrorInternal(w http.ResponseWriter) {
	Error(w, http.StatusInternalServerError, "Something went wrong.")
}

func ErrorNotFound(w http.ResponseWriter) {
	Error(w, http.StatusNotFound, "Resource not found.")
}

func ErrorForbidden(w http.ResponseWriter) {
	Error(w, http.StatusForbidden, "You do not have access for the attempted action.")
}

func Error(w http.ResponseWriter, statusCode int, message string) {
	b, _ := json.Marshal(errorResponse{
		Error: map[string]interface{}{
			"message": message,
		},
	})
	w.Header().Add("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(statusCode)
	_, _ = w.Write(b)
}
