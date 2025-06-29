package api

import (
	"github.com/stretchr/testify/assert"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func Test_RouterAcceptsCorrectMethods(t *testing.T) {
	router := NewRouter("")

	router.Get("/endpoint", func(request *RequestReader) ResponseWriter {
		return SendOk("GET REQUEST")
	})

	req := httptest.NewRequest("GET", "http://localhost:8080/endpoint", nil)

	w := httptest.NewRecorder()
	router.mux.ServeHTTP(w, req)

	resp := w.Result()

	body, err := io.ReadAll(resp.Body)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, "GET REQUEST", string(body))
}
