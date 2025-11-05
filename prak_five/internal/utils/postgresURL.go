package utils

import (
	"net"
	"net/url"
	"os"
)

func getenv(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}

func BuildPostgresURL() (string, error) {

	host := getenv("DB_HOST", "localhost")
	port := getenv("DB_PORT", "5432")
	user := os.Getenv("DB_USER")
	pass := os.Getenv("DB_PASSWORD")
	name := os.Getenv("DB_NAME")

	// host:port
	hp := net.JoinHostPort(host, port)

	u := &url.URL{
		Scheme: "postgres",
		User:   url.UserPassword(user, pass),
		Host:   hp,
		Path:   "/" + name,
	}

	return u.String(), nil
}
