package mux

import "net/http"

type Middleware func(http.HandlerFunc) http.HandlerFunc

type MiddlewareMux struct {
	*http.ServeMux
	middlewares []Middleware
}

func NewMiddlewareMux(middlewares ...Middleware) *MiddlewareMux {
	return &MiddlewareMux{http.NewServeMux(), middlewares}
}

func (m MiddlewareMux) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	handler := m.ServeMux.ServeHTTP
	for i := len(m.middlewares) - 1; i >= 0; i-- {
		handler = m.middlewares[i](handler)
	}
	handler(w, r)
}
