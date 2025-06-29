package api

import (
	"fmt"
	"log"
	"net/http"
)

type Router struct {
	path        string
	middlewares []Middleware
	mux         *http.ServeMux
}

func NewRouter(path string) Router {
	http.NewServeMux()
	return Router{
		path:        path,
		middlewares: []Middleware{},
		mux:         http.NewServeMux(),
	}
}

func (r *Router) ListenAndServe(port int) error {
	log.Printf("Listening at http://localhost:%d", port)
	portString := fmt.Sprintf(":%d", port)
	return http.ListenAndServe(portString, r.mux)
}

func (r *Router) ListenAndServeTLS(port int, certFile string, keyFile string) error {
	log.Printf("Listening at https://localhost:%d", port)
	portString := fmt.Sprintf(":%d", port)
	return http.ListenAndServeTLS(portString, certFile, keyFile, r.mux)
}

func (r *Router) RouteGroup(path string) Router {
	return Router{
		path:        r.path + path,
		middlewares: r.middlewares,
		mux:         r.mux,
	}
}

func (r *Router) Group() Router {
	return r.RouteGroup("")
}

func (r *Router) Use(middleware ...Middleware) {
	r.middlewares = append(r.middlewares, middleware...)
}

func (r *Router) route(method string, path string, endpoint Endpoint) {
	pattern := fmt.Sprintf("%s %s%s", method, r.path, path)
	r.mux.Handle(pattern, newHandler(endpoint, r.middlewares))
}

func (r *Router) Get(path string, handler Endpoint) {
	r.route("GET", path, handler)
}

func (r *Router) Post(path string, handler Endpoint) {
	r.route("POST", path, handler)
}

func (r *Router) Put(path string, handler Endpoint) {
	r.route("PUT", path, handler)
}

func (r *Router) Patch(path string, handler Endpoint) {
	r.route("PATCH", path, handler)
}

func (r *Router) Delete(path string, handler Endpoint) {
	r.route("DELETE", path, handler)
}
