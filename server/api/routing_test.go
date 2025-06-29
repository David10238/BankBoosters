package api

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func Test_RouterAcceptsCorrectMethods(t *testing.T) {
	cases := []string{"GET", "POST", "PUT", "DELETE", "PATCH"}

	createTestCase := func(method string) func(t *testing.T) {
		return func(t *testing.T) {
			router := NewRouter("")

			router.Get("/endpoint", func(request *RequestReader) ResponseWriter {
				return SendOk("GET REQUEST")
			})

			router.Post("/endpoint", func(request *RequestReader) ResponseWriter {
				return SendOk("POST REQUEST")
			})

			router.Put("/endpoint", func(request *RequestReader) ResponseWriter {
				return SendOk("PUT REQUEST")
			})

			router.Delete("/endpoint", func(request *RequestReader) ResponseWriter {
				return SendOk("DELETE REQUEST")
			})

			router.Patch("/endpoint", func(request *RequestReader) ResponseWriter {
				return SendOk("PATCH REQUEST")
			})

			req := httptest.NewRequest(method, "http://localhost:8080/endpoint", nil)

			w := httptest.NewRecorder()
			router.mux.ServeHTTP(w, req)

			resp := w.Result()

			body, err := io.ReadAll(resp.Body)
			assert.NoError(t, err)
			assert.Equal(t, http.StatusOK, resp.StatusCode)
			assert.Equal(t, fmt.Sprintf("%s REQUEST", method), string(body))
		}
	}

	for _, method := range cases {
		name := fmt.Sprintf("Test calling %s method", method)
		t.Run(name, createTestCase(method))
	}
}
