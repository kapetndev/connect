package errors

import (
	"fmt"
	"net/http"
)

type customError string

func New(format string, args ...interface{}) error {
	return customError(fmt.Sprintf(format, args...))
}

func (e customError) Error() string {
	return string(e)
}

func (e customError) RespondError(w http.ResponseWriter, r *http.Request) bool {
	w.WriteHeader(http.StatusBadRequest)
	fmt.Fprintln(w, string(e))
	return true
}
