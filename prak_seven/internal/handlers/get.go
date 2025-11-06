package handlers

import (
	"fmt"
	"github.com/CyberGeo335/prak_seven/internal/cache"
	"net/http"
)

func GetHandler(c *cache.Cache, writer http.ResponseWriter, request *http.Request) {
	key := request.URL.Query().Get("key")
	val, err := c.Get(key)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusNotFound)
		return
	}
	fmt.Fprintf(writer, "VALUE: %s=%s", key, val)
}
