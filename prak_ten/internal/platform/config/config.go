package config

import (
	"log"
	"os"
	"strconv"
	"time"
)

type RSKey struct {
	KID            string
	PrivateKeyPath string
	PublicKeyPath  string
}

type Config struct {
	Port                 string
	Issuer               string
	Audience             string
	AccessTTL            time.Duration
	RefreshTTL           time.Duration
	RateLimitLoginMax    int
	RateLimitLoginWindow time.Duration
	RSKeys               []RSKey
}

func Load() Config {
	port := getenv("APP_PORT", "8080")
	issuer := getenv("JWT_ISSUER", "pz10-auth")
	audience := getenv("JWT_AUDIENCE", "pz10-clients")

	accessTTL := parseDuration("ACCESS_TTL", "15m")
	refreshTTL := parseDuration("REFRESH_TTL", "168h") // 7 дней

	rateMax := parseInt("LOGIN_MAX", "5")
	rateWindow := parseDuration("LOGIN_WINDOW", "5m")

	// Ключ 1 (обязателен)
	kid1 := os.Getenv("JWT_KID_1")
	priv1 := os.Getenv("JWT_PRIVATE_KEY_PATH_1")
	pub1 := os.Getenv("JWT_PUBLIC_KEY_PATH_1")
	if kid1 == "" || priv1 == "" || pub1 == "" {
		log.Fatal("JWT_KID_1, JWT_PRIVATE_KEY_PATH_1, JWT_PUBLIC_KEY_PATH_1 обязательны")
	}
	keys := []RSKey{
		{
			KID:            kid1,
			PrivateKeyPath: priv1,
			PublicKeyPath:  pub1,
		},
	}

	// Ключ 2 (опционально — для ротации)
	kid2 := os.Getenv("JWT_KID_2")
	priv2 := os.Getenv("JWT_PRIVATE_KEY_PATH_2")
	pub2 := os.Getenv("JWT_PUBLIC_KEY_PATH_2")
	if kid2 != "" && priv2 != "" && pub2 != "" {
		keys = append(keys, RSKey{
			KID:            kid2,
			PrivateKeyPath: priv2,
			PublicKeyPath:  pub2,
		})
	}

	return Config{
		Port:                 ":" + port,
		Issuer:               issuer,
		Audience:             audience,
		AccessTTL:            accessTTL,
		RefreshTTL:           refreshTTL,
		RateLimitLoginMax:    rateMax,
		RateLimitLoginWindow: rateWindow,
		RSKeys:               keys,
	}
}

func getenv(key, def string) string {
	v := os.Getenv(key)
	if v == "" {
		return def
	}
	return v
}

func parseDuration(envKey, def string) time.Duration {
	v := getenv(envKey, def)
	d, err := time.ParseDuration(v)
	if err != nil {
		log.Fatalf("bad %s: %v", envKey, err)
	}
	return d
}

func parseInt(envKey, def string) int {
	v := getenv(envKey, def)
	i, err := strconv.Atoi(v)
	if err != nil {
		log.Fatalf("bad %s: %v", envKey, err)
	}
	return i
}
