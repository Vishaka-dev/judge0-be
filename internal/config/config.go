package config

import (
	"os"
	"sync"

	"github.com/joho/godotenv"
)

type Config struct {
	SupabaseURL        string
	SupabaseAnonKey    string
	SupabaseServiceKey string
	SecretKey          string
	Judge0API          string
}

var (
	cfg  *Config
	once sync.Once
)

func Load() *Config {
	once.Do(func() {
		godotenv.Load()

		cfg = &Config{
			SupabaseURL:        os.Getenv("SUPABASE_URL"),
			SupabaseAnonKey:    os.Getenv("SUPABASE_ANON_KEY"),
			SupabaseServiceKey: os.Getenv("SUPABASE_SERVICE_KEY"),
			SecretKey:          os.Getenv("SECRET_KEY"),
			Judge0API:          os.Getenv("JUDGE0_API"),
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
