package api

import (
	"net/http"
)

type ResponseWriter interface {
	Write(w *http.ResponseWriter) error
}
