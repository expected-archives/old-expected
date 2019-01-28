package response

import (
	"encoding/json"
	"net/http"
)

func SingleResource(w http.ResponseWriter, name string, v interface{}) {
	b, _ := json.Marshal(map[string]interface{}{name: v})
	w.Header().Add("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(200)
	_, _ = w.Write(b)
}

func MultipleResource(w http.ResponseWriter, name string, v []interface{}) {
	b, _ := json.Marshal(map[string]interface{}{name: v})
	w.Header().Add("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(200)
	_, _ = w.Write(b)
}
