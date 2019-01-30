package api

import (
	"fmt"
	"net/http"

	"github.com/g0lang/proxy/src/config"
)

func versionGet(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, config.Config.Get("HASH"))
}
