package api

import (
	"net/http"

	log "github.com/sirupsen/logrus"

	"github.com/g0lang/proxy/src/config"
	"github.com/gorilla/mux"
)

// Run api server
func Run() {

	// get port from config package
	var PORT string
	if PORT = config.Config.Get("PORT"); PORT == "" {
		PORT = "8000"
	}

	// init router
	router := mux.NewRouter()
	router.HandleFunc("/proxy/{url}", proxyGet).Methods("GET")
	router.HandleFunc("/proxy/{url}", proxyPost).Methods("POST")
	router.HandleFunc("/version", versionGet).Methods("GET")

	// init http server
	log.Println("Starting Server On Port:", PORT)
	log.Fatal(http.ListenAndServe(":"+PORT, rewriter(router)))
}
