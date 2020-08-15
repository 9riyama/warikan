package rest

import (
	"encoding/json"
	"fmt"
	"net/http"

	"go.uber.org/zap"

	"github.com/warikan/api/usecase"
	"github.com/warikan/log"
)

type errorMessage struct {
	Message string `json:"msg"`
}

const (
	badRequestErrorMsg     = "要求の形式が正しくありません。"
	notFoundErrorMsg       = "ページが見つかりません。"
	internalServerErrorMsg = "システム内部エラーが発生しました。"
	conflictErrorMsg       = "競合が発生しました。"
)

func httpError(w http.ResponseWriter, err error, msg string) {
	switch err.(type) {
	case usecase.UnauthorizedError:
		unauthorizedError(w, msg)
	case usecase.BadRequestError:
		badRequestError(w, msg)
	case usecase.InvalidParamError:
		badRequestError(w, msg)
	case usecase.NotFoundError:
		notFoundError(w, msg)
	case usecase.ConflictError:
		conflictError(w, msg)
	default:
		internalServerError(w, msg)
	}
}

func unauthorizedError(w http.ResponseWriter, msg string) {
	code := http.StatusUnauthorized
	m := msg
	if msg == "" {
		m = "サーバーとの認証に失敗しました。再度ログインしてください。"
	}
	errorResponse(w, code, m)
}

func badRequestError(w http.ResponseWriter, msg string) {
	code := http.StatusBadRequest
	m := msg
	if msg == "" {
		m = badRequestErrorMsg
	}
	errorResponse(w, code, m)
}

func notFoundError(w http.ResponseWriter, msg string) {
	code := http.StatusNotFound
	m := msg
	if msg == "" {
		m = notFoundErrorMsg
	}
	errorResponse(w, code, m)
}

func internalServerError(w http.ResponseWriter, msg string) {
	code := http.StatusInternalServerError
	m := msg
	if msg == "" {
		m = internalServerErrorMsg
	}
	errorResponse(w, code, m)
}

func conflictError(w http.ResponseWriter, msg string) {
	code := http.StatusConflict
	m := msg
	if msg == "" {
		m = conflictErrorMsg
	}
	errorResponse(w, code, m)
}

func loginError(w http.ResponseWriter) {
	code := http.StatusBadRequest
	m := "ログインに失敗しました。ログインできない場合は「パスワードをお忘れの方」からパスワードの再設定を行ってください。"
	errorResponse(w, code, m)
}

func errorResponse(w http.ResponseWriter, code int, msg string) {
	w.WriteHeader(code)
	em := errorMessage{Message: msg}
	err := json.NewEncoder(w).Encode(em)
	if err != nil {
		log.Logger.Error("failed to json encode", zap.Error(err))
	}
}

func badRequestErrorText(w http.ResponseWriter, msg string) {
	txt := badRequestErrorMsg
	if msg != "" {
		txt = msg
	}
	w.WriteHeader(http.StatusBadRequest)
	fmt.Fprint(w, txt)
	w.Header().Set("Content-Type", "text/html")
}

func internalServerErrorText(w http.ResponseWriter, msg string) {
	txt := internalServerErrorMsg
	if msg != "" {
		txt = msg
	}
	w.WriteHeader(http.StatusInternalServerError)
	fmt.Fprint(w, txt)
	w.Header().Set("Content-Type", "text/html")
}
