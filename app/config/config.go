package config

import (
	"os"
	"sync"

	"github.com/joho/godotenv"
)

type Config struct {
	SecretKey string
	Judge0API string
	PGURL     string
}

var (
	cfg  *Config
	once sync.Once
)

func Load() *Config {
	once.Do(func() {
		godotenv.Load()

		cfg = &Config{
			SecretKey: os.Getenv("SECRET_KEY"),
			Judge0API: os.Getenv("JUDGE0_API"),
			PGURL:     os.Getenv("PG_URL"),
		}
	})
	return cfg
}

func Get() *Config {
	if cfg == nil {
		return Load()
	}
	return cfg
}
