package api

import (
	"log"
	"net/http"
)

type Endpoint func(request *RequestReader) ResponseWriter
type Middleware = func(request *RequestReader) *ResponseWriter

type handler struct {
	endpoint    Endpoint
	middlewares []Middleware
}

func newHandler(endpoint Endpoint, middlewares []Middleware) handler {
	return handler{
		endpoint:    endpoint,
		middlewares: middlewares,
	}
}

func (h handler) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	reader := newRequestReader(req)
	for _, middleware := range h.middlewares {
		if response := middleware(reader); response != nil {
			if err := (*response).Write(&res); err != nil {
				log.Fatal(err)
			}
			return
		}
	}

	if err := h.endpoint(reader).Write(&res); err != nil {
		log.Fatal(err)
	}
}
