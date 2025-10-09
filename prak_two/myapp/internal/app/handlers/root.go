package handlers

import (
	"fmt"
	"github.com/CyberGeo335/myapp/utils"
	"net/http"
)

func Root(w http.ResponseWriter, r *http.Request) {
	utils.LogRequest(r)
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	fmt.Fprintln(w, "Hello, Go project structure!")
}
