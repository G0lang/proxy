package api

import (
	"net/http"
	"net/url"
	"strings"
)

func rewriter(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		pathReq := r.RequestURI
		if strings.HasPrefix(pathReq, "/proxy/") {
			pe := url.PathEscape(strings.TrimLeft(pathReq, "/proxy/"))
			r.URL.Path = "/proxy/" + pe
			r.URL.RawQuery = ""
		}
		h.ServeHTTP(w, r)
	})
}
