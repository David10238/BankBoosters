package api

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func Test_ReadingStringHeaders(t *testing.T) {
	type Result struct {
		Status int
		Body   string
	}

	testCases := []struct {
		value    string
		expected Result
	}{
		{"", Result{
			Status: http.StatusNotFound,
			Body:   "Header input is empty",
		}},
		{"findMe", Result{
			Status: http.StatusOK,
			Body:   "findMe",
		}},
	}

	makeCase := func(value string, expected Result) func(t *testing.T) {
		return func(t *testing.T) {
			router := NewRouter("")
			router.Get("/echo", func(request *RequestReader) ResponseWriter {
				header := ""
				if err := request.BindStringHeader("input", &header); err != nil {
					return err
				}
				return SendOk(header)
			})

			req := httptest.NewRequest("GET", "http://localhost:8080/echo", nil)
			req.Header.Set("input", value)

			w := httptest.NewRecorder()
			router.mux.ServeHTTP(w, req)

			resp := w.Result()

			body, readErr := io.ReadAll(resp.Body)
			assert.NoError(t, readErr)

			assert.Equal(t, expected, Result{
				Status: resp.StatusCode,
				Body:   string(body),
			})
		}
	}

	for _, testCase := range testCases {
		name := fmt.Sprintf("header set to '%s;", testCase.value)
		t.Run(name, makeCase(testCase.value, testCase.expected))
	}
}

func Test_ReadingJsonHeaders(t *testing.T) {
	type Result struct {
		Status int
		Body   string
	}

	testCases := []struct {
		value    string
		expected Result
	}{
		{"", Result{
			Status: http.StatusNotFound,
			Body:   "Header input is empty",
		}},
		{"invalid json", Result{
			Status: http.StatusNotFound,
			Body:   "Header input is malformed",
		}},
		{"6", Result{
			Status: http.StatusOK,
			Body:   "6",
		}},
	}

	makeCase := func(value string, expected Result) func(t *testing.T) {
		return func(t *testing.T) {
			router := NewRouter("")
			router.Get("/echo", func(request *RequestReader) ResponseWriter {
				header := 0
				if err := request.BindJsonHeader("input", &header); err != nil {
					return err
				}
				return SendJson(header)
			})

			req := httptest.NewRequest("GET", "http://localhost:8080/echo", nil)
			req.Header.Set("input", value)

			w := httptest.NewRecorder()
			router.mux.ServeHTTP(w, req)

			resp := w.Result()

			body, readErr := io.ReadAll(resp.Body)
			assert.NoError(t, readErr)

			assert.Equal(t, expected, Result{
				Status: resp.StatusCode,
				Body:   string(body),
			})
		}
	}

	for _, testCase := range testCases {
		name := fmt.Sprintf("header set to '%s;", testCase.value)
		t.Run(name, makeCase(testCase.value, testCase.expected))
	}
}
