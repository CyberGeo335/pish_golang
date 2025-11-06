package handlers

import (
	"fmt"
	"github.com/CyberGeo335/prak_seven/internal/cache"
	"net/http"
	"time"
)

func SetHandler(c *cache.Cache, writer http.ResponseWriter, request *http.Request) {
	key := request.URL.Query().Get("key")
	value := request.URL.Query().Get("value")
	if key == "" || value == "" {
		http.Error(writer, "missing parameters.", http.StatusBadRequest)
		return
	}
	err := c.Set(key, value, 10*time.Second) // TTL = 10 сек
	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}
	fmt.Fprintf(writer, "OK: %s=%s (TTL 10s)", key, value)
}
