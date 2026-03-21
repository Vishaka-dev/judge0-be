package config

import (
	"os"
	"sync"

	"github.com/Mozilla-Campus-Club-of-SLIIT/judge0-be/app/logger"
	"github.com/joho/godotenv"
)

type Config struct {
	SecretKey         string
	Judge0API         string
	PGURL             string
	AUTH_API          string
	Judge0CallbackURL string
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
		logger.Log.Info("loaded .env file")
	} else {
		logger.Log.Info("Running on Vercel, skipping .env file")
	}

	cfg = &Config{
		SecretKey:         must("SECRET_KEY"),
		Judge0API:         must("JUDGE0_API"),
		PGURL:             must("PG_URL"),
		AUTH_API:          must("AUTH_API"),
		Judge0CallbackURL: must("JUDGE0_CALLBACK_URL"),
	}
}

func must(key string) string {
	val := os.Getenv(key)
	if val == "" {
		logger.Log.Fatal("missing env var", "key", key)
	}
	return val
}
