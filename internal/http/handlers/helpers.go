package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"

	forum "stepik.leoscode.http/internal/gen/api"
)

func JsonEncode(w http.ResponseWriter, toEncode any) {
	if err := json.NewEncoder(w).Encode(toEncode); err != nil {
		http.Error(w, "json encode error", http.StatusInternalServerError)
	}
}

func writeError(w http.ResponseWriter, status int, errResp forum.ErrorResponse) {
	writeJson(w, status, errResp)
}

func writeJson(w http.ResponseWriter, status int, msg any) {
	var buf bytes.Buffer

	if err := json.NewEncoder(&buf).Encode(msg); err != nil {
		http.Error(w, "json encode error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(status)
	w.Write(buf.Bytes())
}

func ApiErrorHandler(w http.ResponseWriter, r *http.Request, err error) {
	w.Header().Set("Content-Type", "application/json")
	status := http.StatusBadRequest

	resp := forum.ErrorResponse{
		Code:    "validation_error",
		Message: err.Error(),
	}

	writeJson(w, status, resp)
}
func readJson(r *http.Request, toWrite any) error {
	return json.NewDecoder(r.Body).Decode(toWrite)
}
