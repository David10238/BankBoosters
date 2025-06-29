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

func Test_RouterNestingWithAndWithoutGroups(t *testing.T) {
	cases := []struct {
		route    string
		expected string
	}{
		{"/api/a", "a"},
		{"/api/a/a", "aa"},
		{"/api/b/b", "b"},
		{"/api/b/b/b", "bb"},
	}

	createTestCase := func(route string, expected string) func(t *testing.T) {
		return func(t *testing.T) {
			router := NewRouter("/api")
			group := router.RouteGroup("/b")

			router.Get("/a", func(request *RequestReader) ResponseWriter {
				return SendOk("a")
			})

			router.Get("/a/a", func(request *RequestReader) ResponseWriter {
				return SendOk("aa")
			})

			group.Get("/b", func(request *RequestReader) ResponseWriter {
				return SendOk("b")
			})

			group.Get("/b/b", func(request *RequestReader) ResponseWriter {
				return SendOk("bb")
			})

			url := fmt.Sprintf("http://localhost:8080%s", route)
			req := httptest.NewRequest("GET", url, nil)

			w := httptest.NewRecorder()
			router.mux.ServeHTTP(w, req)

			resp := w.Result()

			body, err := io.ReadAll(resp.Body)
			assert.NoError(t, err)
			assert.Equal(t, http.StatusOK, resp.StatusCode)
			assert.Equal(t, expected, string(body))
		}
	}

	for _, c := range cases {
		name := fmt.Sprintf("Route %s", c.route)
		t.Run(name, createTestCase(c.route, c.expected))
	}
}

func Test_RouterReturnsNotFoundOnNonexistentRoute(t *testing.T) {
	router := NewRouter("/api")

	req := httptest.NewRequest("GET", "http://localhost:8080/thisIsNotAnEndpoint", nil)

	w := httptest.NewRecorder()
	router.mux.ServeHTTP(w, req)

	resp := w.Result()

	assert.Equal(t, http.StatusNotFound, resp.StatusCode)
}
