package api

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRouting(t *testing.T) {
	srv := httptest.NewServer(rewriter(Router()))
	defer srv.Close()

	res, err := http.Get(fmt.Sprintf("%s/proxy/http://httpbin.org/get", srv.URL))
	if err != nil {
		t.Fatalf("could not send get request: %v ", err)
	}
	if res.StatusCode != http.StatusOK {
		t.Fatalf("expect status ok ; got %v", res.StatusCode)
	}
}
