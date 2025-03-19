package mux

import (
	"net/http"
)

type PrefixMux struct {
	*http.ServeMux
}

func NewPrefixMux() *PrefixMux {
	return &PrefixMux{http.NewServeMux()}
}

func (mux *PrefixMux) Handle(prefix string, handler http.Handler) {
	mux.ServeMux.Handle(prefix, http.StripPrefix(prefix, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if len(r.URL.Path) == 0 || r.URL.Path[0] != '/' {
			r.URL.Path = "/" + r.URL.Path
		}
		if len(r.URL.RawPath) == 0 || r.URL.RawPath[0] != '/' {
			r.URL.RawPath = "/" + r.URL.RawPath
		}
		handler.ServeHTTP(w, r)
	})))
	mux.ServeMux.Handle(prefix+"/", http.StripPrefix(prefix, handler))
}
