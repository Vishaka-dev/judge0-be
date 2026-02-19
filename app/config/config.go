package config

import (
	"log"
	"os"
	"sync"

	"github.com/joho/godotenv"
)

type Config struct {
	SecretKey string
	Judge0API string
	PGURL     string
	AUTH_API  string
}

var (
	cfg  *Config
	once sync.Once
)

func Get() *Config {
	once.Do(load)
	return cfg
}

func load() {
	if os.Getenv("VERCEL") == "" {
		_ = godotenv.Load()
		log.Println("loaded .env file")
	} else {
		log.Println("Running on Vercel, skipping .env file")
	}

	cfg = &Config{
		SecretKey: must("SECRET_KEY"),
		Judge0API: must("JUDGE0_API"),
		PGURL:     must("PG_URL"),
		AUTH_API:  must("AUTH_API"),
	}
}

func must(key string) string {
	val := os.Getenv(key)
	if val == "" {
		log.Fatalf("missing env var: %s", key)
	}
	return val
}
