package api

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

func Run() {
	router := mux.NewRouter()
	router.HandleFunc("/proxy/{url}", proxyGet).Methods("GET")
	router.HandleFunc("/proxy/{url}", proxyPost).Methods("POST")
	log.Println("Starting Server On Port 8000")
	http.ListenAndServe(":8000", rewriter(router))

}

func proxyGet(w http.ResponseWriter, r *http.Request) {
	originalURL := getUrlFromRequest(r)
	if originalURL == "" {
		fmt.Fprintf(w, "Please Enter Url")
	} else {
		req := buildRequest(originalURL, "GET", nil, false)
		body := makeRequest(req)
		fmt.Fprint(w, string(body))
	}
}

func proxyPost(w http.ResponseWriter, r *http.Request) {
	originalURL := getUrlFromRequest(r)
	if originalURL == "" {
		fmt.Fprintf(w, "Please Enter Url")
		return
	} else {
		r.ParseForm()
		req := buildRequest(originalURL, "POST", r.Form, true)
		body := makeRequest(req)
		fmt.Fprint(w, string(body))
	}
}

func getUrlFromRequest(r *http.Request) string {
	vars := mux.Vars(r)
	u := vars["url"]
	originalURL, err := url.PathUnescape(u)
	if err != nil {
		log.Println("Can Not Build Url")
		return ""
	}
	_, err = url.ParseRequestURI(originalURL)
	if err != nil {
		log.Println("Url Not correct")
		return ""
	}
	return originalURL
}

func buildRequest(url, method string, data url.Values, urlencoded bool) *http.Request {
	req, err := http.NewRequest(method, url, strings.NewReader(data.Encode()))
	if err != nil {
		log.Println("Can Not Create req Error:", err)
	}
	if urlencoded {
		req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	}
	req.Header.Set("User-Agent", "curl/7.35.0")
	return req
}

func makeRequest(req *http.Request) string {
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
