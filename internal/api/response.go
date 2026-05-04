package api

import (
	"encoding/json"
	"net/http"
)

type Response struct {
	Success bool   `json:"s"`
	Code    int    `json:"c"`
	Data    any    `json:"d,omitempty"`
	Error   string `json:"e,omitempty"`
}

func SendJSON(w http.ResponseWriter, code int, success bool, data any, errMsg string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)

	res := Response{
		Success: success,
		Code:    code,
		Data:    data,
		Error:   errMsg,
	}

	json.NewEncoder(w).Encode(res)
}
