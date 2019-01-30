package response

import (
	"encoding/json"
	"net/http"
)

type errorMessage struct {
	Message string `json:"message"`
}

func ErrorBadRequest(w http.ResponseWriter, message string) {
	b, _ := json.Marshal(errorMessage{Message: message})
	w.Header().Add("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusBadRequest)
	_, _ = w.Write(b)
}

func ErrorInternal(w http.ResponseWriter, message string) {
	b, _ := json.Marshal(errorMessage{Message: message})
	w.Header().Add("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusInternalServerError)
	_, _ = w.Write(b)
}

func ErrorForbidden(w http.ResponseWriter) {
	b, _ := json.Marshal(errorMessage{Message: "You do not have access for the attempted action."})
	w.Header().Add("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusForbidden)
	_, _ = w.Write(b)
}
