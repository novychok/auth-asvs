package authapiv1

import (
	"encoding/json"
	"log"
	"net/http"

	authapiv1 "github.com/novychok/authasvs/pkg/authApi/v1"
)

func response(w http.ResponseWriter, code int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)

	if data != nil {
		if err := json.NewEncoder(w).Encode(data); err != nil {
			log.Println(err)
		}
	}
}

func errResponse(w http.ResponseWriter, _ *http.Request, code int, message string) {
	response(w, code, authapiv1.Error{
		Code:    code,
		Message: message,
	})
}

func errBadRequest(w http.ResponseWriter, r *http.Request, message string) {
	errResponse(w, r, http.StatusBadRequest, message)
}

func errInternal(w http.ResponseWriter, r *http.Request, message string) {
	errResponse(w, r, http.StatusInternalServerError, message)
}

func errNotFound(w http.ResponseWriter, r *http.Request, message string) {
	errResponse(w, r, http.StatusNotFound, message)
}

func errUnauthorized(w http.ResponseWriter, r *http.Request, message string) {
	errResponse(w, r, http.StatusUnauthorized, message)
}
