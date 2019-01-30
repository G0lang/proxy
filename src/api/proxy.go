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

func proxyGet(w http.ResponseWriter, r *http.Request) {
	originalURL := getURLFromRequest(r)
	if originalURL == "" {
		fmt.Fprintf(w, "Please Enter Url")
		return
	}
	req := buildRequest(originalURL, "GET", nil, false)
	body := makeRequest(req)
	fmt.Fprint(w, string(body))

}

func proxyPost(w http.ResponseWriter, r *http.Request) {
	originalURL := getURLFromRequest(r)
	if originalURL == "" {
		fmt.Fprintf(w, "Please Enter Url")
		return
	}
	r.ParseForm()
	req := buildRequest(originalURL, "POST", r.Form, true)
	body := makeRequest(req)
	fmt.Fprint(w, string(body))

}

func getURLFromRequest(r *http.Request) string {
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
