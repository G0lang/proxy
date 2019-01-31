package api

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestProxyGet(t *testing.T) {
	srv := httptest.NewServer(rewriter(Router()))
	defer srv.Close()

	res, err := http.Get(fmt.Sprintf("%s/proxy/http://httpbin.org/get", srv.URL))
	if err != nil {
		t.Fatalf("could not send get request: %v ", err)
	}
	if res.StatusCode != http.StatusOK {
		t.Fatalf("expect status ok ; got %v", res.StatusCode)
	}

	b, err := ioutil.ReadAll(res.Body)
	if err != nil {
		t.Fatalf("count not read response %v", err)
	}
	expBody := `{
  "args": {}, 
  "headers": {
    "Accept-Encoding": "gzip", 
    "Connection": "close", 
    "Host": "httpbin.org", 
    "User-Agent": "curl/7.35.0"
  }, 
  "origin": "216.181.132.38", 
  "url": "http://httpbin.org/get"
}
`
	exbody, err := json.Marshal(expBody)
	body, err := json.Marshal(string(b))
	if string(exbody) != string(body) {
		t.Fatalf("\n we expect:%v \n but get b:%v", string(exbody), string(body))
	}
}

func TestProxyPost(t *testing.T) {
	srv := httptest.NewServer(rewriter(Router()))
	defer srv.Close()
	// TODO pass payload !
	res, err := http.Post(fmt.Sprintf("%s/proxy/http://httpbin.org/post", srv.URL),
		"application/x-www-form-urlencoded", strings.NewReader(""))
	if err != nil {
		t.Fatalf("could not send get request: %v ", err)
	}
	if res.StatusCode != http.StatusOK {
		t.Fatalf("expect status ok ; got %v", res.StatusCode)
	}

	b, err := ioutil.ReadAll(res.Body)
	if err != nil {
		t.Fatalf("count not read response %v", err)
	}
	expBody := `{
  "args": {}, 
  "data": "", 
  "files": {}, 
  "form": {}, 
  "headers": {
    "Accept-Encoding": "gzip", 
    "Connection": "close", 
    "Content-Length": "0", 
    "Content-Type": "application/x-www-form-urlencoded", 
    "Host": "httpbin.org", 
    "User-Agent": "curl/7.35.0"
  }, 
  "json": null, 
  "origin": "216.181.132.38", 
  "url": "http://httpbin.org/post"
}
`
	exbody, err := json.Marshal(expBody)
	body, err := json.Marshal(string(b))
	if string(exbody) != string(body) {
		t.Fatalf("\n we expect:%v \n but get b:%v", string(exbody), string(body))
	}
}
