package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
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

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/proxy/{url}", proxyGET).Methods("GET")
	router.HandleFunc("/proxy/{url}", proxyPOST).Methods("POST")
	log.Println("Starting Server On Port 8000")
	http.ListenAndServe(":8000", rewriter(router))

}

func proxyGET(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	u := vars["url"]
	originalURL, err := url.PathUnescape(u)
	if err != nil {
		log.Println("Can Not Build Url")
		return
	}
	body := makeRequest(originalURL, "GET", nil, false)

	fmt.Fprint(w, string(body))
}

func proxyPOST(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	u := vars["url"]
	originalURL, err := url.PathUnescape(u)
	if err != nil {
		log.Println("Can Not Build Url")
		return
	}
	r.ParseForm()
	body := makeRequest(originalURL, "POST", r.Form, true)

	fmt.Fprint(w, string(body))
}

func makeRequest(url, method string, data url.Values, urlencoded bool) string {
	req, err := http.NewRequest(method, url, strings.NewReader(data.Encode()))

	if err != nil {
		log.Println("Can Not Create req Error:", err)
	}

	if urlencoded {
		req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	}

	req.Header.Set("User-Agent", "curl/7.35.0")

	client := &http.Client{}
	resp, err := client.Do(req)

	if err != nil {
		log.Println("Can Not Get Any Feedback From Host Error:", err)
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println("Can Not Parse The Body Error:", err)
	}
	return string(body)
}
