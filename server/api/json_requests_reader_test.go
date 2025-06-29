package api

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func Test_JsonRequestReadsValidJson(t *testing.T) {
	router := NewRouter("")

	type ExpectedStruct struct {
		A int
		B string
		C bool
	}

	expected := ExpectedStruct{
		A: 42,
		B: "Hello World",
		C: true,
	}

	router.Get("/json", func(request *RequestReader) ResponseWriter {
		body := ExpectedStruct{}
		if err := request.BindJsonBody(&body); err != nil {
			return err
		}
		return SendJson(body)
	})

	marshaled, marshalErr := json.Marshal(expected)
	assert.NoError(t, marshalErr)

	req := httptest.NewRequest("GET", "http://localhost:8080/json", strings.NewReader(string(marshaled)))
	req.Header.Set("Content-Type", "application/json; charset=utf-8")

	w := httptest.NewRecorder()
	router.mux.ServeHTTP(w, req)

	resp := w.Result()

	body, readErr := io.ReadAll(resp.Body)
	assert.NoError(t, readErr)

	type Result struct {
		Status      int
		ContentType string
		Body        string
	}

	assert.Equal(t, Result{
		Status:      http.StatusOK,
		ContentType: "application/json; charset=utf-8",
		Body:        string(marshaled),
	}, Result{
		Status:      resp.StatusCode,
		ContentType: resp.Header.Get("Content-Type"),
		Body:        string(body),
	})
}

func Test_JsonRequestRejectsMissingJson(t *testing.T) {
	router := NewRouter("")

	type ExpectedStruct struct {
		A int
		B string
		C bool
	}

	router.Get("/json", func(request *RequestReader) ResponseWriter {
		body := ExpectedStruct{}
		if err := request.BindJsonBody(&body); err != nil {
			return err
		}
		return SendJson(body)
	})

	req := httptest.NewRequest("GET", "http://localhost:8080/json", nil)
	req.Header.Set("Content-Type", "application/json; charset=utf-8")

	w := httptest.NewRecorder()
	router.mux.ServeHTTP(w, req)

	resp := w.Result()

	body, readErr := io.ReadAll(resp.Body)
	assert.NoError(t, readErr)

	type Result struct {
		Status int
		Body   string
	}

	assert.Equal(t, Result{
		Status: http.StatusNotFound,
		Body:   "Not able to bind body",
	}, Result{
		Status: resp.StatusCode,
		Body:   string(body),
	})
}

func Test_JsonRequestRejectsMalformedJson(t *testing.T) {
	router := NewRouter("")

	type ExpectedStruct struct {
		A int
		B string
		C bool
	}

	router.Get("/json", func(request *RequestReader) ResponseWriter {
		body := ExpectedStruct{}
		if err := request.BindJsonBody(&body); err != nil {
			return err
		}
		return SendJson(body)
	})

	req := httptest.NewRequest("GET", "http://localhost:8080/json", strings.NewReader("oh look this isn't json"))
	req.Header.Set("Content-Type", "application/json; charset=utf-8")

	w := httptest.NewRecorder()
	router.mux.ServeHTTP(w, req)

	resp := w.Result()

	body, readErr := io.ReadAll(resp.Body)
	assert.NoError(t, readErr)

	type Result struct {
		Status int
		Body   string
	}

	assert.Equal(t, Result{
		Status: http.StatusNotFound,
		Body:   "Not able to bind body",
	}, Result{
		Status: resp.StatusCode,
		Body:   string(body),
	})
}
