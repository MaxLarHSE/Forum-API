package handlers

import (
	"bytes"
	"encoding/json"
	"errors"
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
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(buf.Bytes())
}

func ApiErrorHandler(w http.ResponseWriter, r *http.Request, err error) {
	var invalid *forum.InvalidParamFormatError

	if errors.As(err, &invalid) {
		if invalid.ParamName == "X-User-Id" {
			writeError(w, http.StatusUnauthorized, forum.ErrorUnauthorized{
				Code:    forum.Unauthorized,
				Message: err.Error(),
			})
			return
		}
		writeError(w, http.StatusBadRequest, forum.ErrorBadRequest{
			Code:    forum.BadRequest,
			Message: err.Error(),
		})
		return
	}
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
