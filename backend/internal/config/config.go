package config

import (
	"os"
	"strconv"
	"time"
)

type Config struct {
	DatabaseURL   string
	HTTPAddr      string
	TickInterval  time.Duration
	NakamaURL     string
	NakamaHTTPKey string
	CORSOrigin    string
	UseMemory     bool
}

func Load() Config {
	tickMs, _ := strconv.Atoi(getenv("TICK_INTERVAL_MS", "2000"))
	if tickMs < 500 {
		tickMs = 500
	}
	return Config{
		DatabaseURL:   getenv("DATABASE_URL", "postgres://bashkir:bashkir@localhost:5432/bashkir_eco?sslmode=disable"),
		HTTPAddr:      getenv("HTTP_ADDR", ":8080"),
		TickInterval:  time.Duration(tickMs) * time.Millisecond,
		NakamaURL:     getenv("NAKAMA_URL", "http://localhost:7350"),
		NakamaHTTPKey: getenv("NAKAMA_HTTP_KEY", "bashkirecokey"),
		CORSOrigin:    getenv("CORS_ORIGIN", "*"),
		UseMemory:     getenv("USE_MEMORY", "false") == "true",
	}
}

func getenv(k, def string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}
	return def
}
