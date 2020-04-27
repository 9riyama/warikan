package rest

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/warikan/api/usecase"
)

type errorMessage struct {
	Message string `json:"message"`
}

func httpError(w http.ResponseWriter, err error) {
	switch err.(type) {
	case usecase.UnauthorizedError:
		unauthorizedError(w)
	case usecase.NotFoundError:
		notFoundError(w)
	case usecase.InvalidParamError:
		badRequestError(w)
	case usecase.InternalServerError:
		internalServerError(w)
	case usecase.ConflictError:
		conflictError(w)
	default:
		internalServerError(w)
	}
}

func unauthorizedError(w http.ResponseWriter) {
	code := http.StatusUnauthorized
	errorResponse(w, code)
}

func badRequestError(w http.ResponseWriter) {
	code := http.StatusBadRequest
	errorResponse(w, code)
}

func notFoundError(w http.ResponseWriter) {
	code := http.StatusNotFound
	errorResponse(w, code)
}

func conflictError(w http.ResponseWriter) {
	code := http.StatusConflict
	errorResponse(w, code)
}

func internalServerError(w http.ResponseWriter) {
	code := http.StatusInternalServerError
	errorResponse(w, code)
}

func errorResponse(w http.ResponseWriter, code int) {
	w.WriteHeader(code)
	em := errorMessage{Message: http.StatusText(code)}
	err := json.NewEncoder(w).Encode(em)
	if err != nil {
		log.Printf("failed to json encode %v\n", err)
	}
}
