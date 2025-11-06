package handlers

import (
	"fmt"
	"github.com/CyberGeo335/prak_seven/internal/cache"
	"net/http"
)

func TtlHandler(c *cache.Cache, writer http.ResponseWriter, request *http.Request) {
	key := request.URL.Query().Get("key")
	ttl, err := c.TTL(key)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}
	fmt.Fprintf(writer, "TTL for %s: %v", key, ttl)
}
