package api

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func Test_MessageResultsSendProperly(t *testing.T) {
	cases := []struct {
		result  ResponseWriter
		code    int
		message string
	}{
		{SendOk("Ok"), http.StatusOK, "Ok"},
		{SendInternalServerError(), http.StatusInternalServerError, "Internal Server Error"},
		{SendBadRequest("Bad Request"), http.StatusBadRequest, "Bad Request"},
		{SendUnauthorized("Unauthorized"), http.StatusUnauthorized, "Unauthorized"},
		{SendForbidden("Forbidden"), http.StatusForbidden, "Forbidden"},
		{SendNotFound("Not Found"), http.StatusNotFound, "Not Found"},
	}

	createTestCase := func(result ResponseWriter, code int, message string) func(t *testing.T) {
		return func(t *testing.T) {
			router := NewRouter("")
			router.Get("/endpoint", func(request *RequestReader) ResponseWriter {
				return result
			})

			req := httptest.NewRequest("GET", "http://localhost:8080/endpoint", nil)

			w := httptest.NewRecorder()
			router.mux.ServeHTTP(w, req)

			resp := w.Result()

			body, err := io.ReadAll(resp.Body)
			assert.NoError(t, err)
			assert.Equal(t, code, resp.StatusCode)
			assert.Equal(t, message, string(body))
		}
	}

	for _, c := range cases {
		name := fmt.Sprintf("Test calling %s", c.message)
		t.Run(name, createTestCase(c.result, c.code, c.message))
	}
}
