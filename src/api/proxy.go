package api

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
)

// MyError is an error implementation that includes a time and message.
type MyError struct {
	When time.Time
	What string
}

func (e MyError) Error() string {
	return fmt.Sprintf("%v: %v", e.When, e.What)
}

func proxyGet(w http.ResponseWriter, r *http.Request) {

	start := time.Now()

	originalURL, err := getURLFromRequest(r)
	if err != nil {
		fmt.Fprintf(w, "Please Enter Url Error: %v", err)
		return
	}

	req := buildRequest(originalURL, "GET", nil, false)
	body := makeRequest(req)
	elapsed := time.Since(start)
	log.Printf("Preparing request and get back respose took %s", elapsed)

	fmt.Fprint(w, string(body))

}

func proxyPost(w http.ResponseWriter, r *http.Request) {

	start := time.Now()

	originalURL, err := getURLFromRequest(r)
	if err != nil {
		fmt.Fprintf(w, "Not Valid URL Error:%v", err)
		return
	}
	r.ParseForm()

	req := buildRequest(originalURL, "POST", r.Form, true)
	body := makeRequest(req)
	elapsed := time.Since(start)
	log.Printf("preparing request amd get back respose took %s", elapsed)

	fmt.Fprint(w, string(body))
}

func getURLFromRequest(r *http.Request) (string, error) {
	vars := mux.Vars(r)
	u := vars["url"]
	originalURL, err := url.PathUnescape(u)
	if err != nil {
		return "", MyError{time.Now(), "1"}
	}
	_, err = url.ParseRequestURI(originalURL)
	if err != nil {
		return "", MyError{time.Now(), "asdaad" + u}
	}
	url, err := url.Parse(originalURL)
	if url.Scheme == "https" {
		return "", MyError{time.Now(), "3"}
	}
	return originalURL, nil
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
