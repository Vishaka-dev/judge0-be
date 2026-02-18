package supabase

import (
	"log"
	"sync"

	"github.com/Mozilla-Campus-Club-of-SLIIT/judge0-be/internal/config"
	supa "github.com/supabase-community/supabase-go"
)

var (
	anonClient  *supa.Client
	adminClient *supa.Client
	once        sync.Once
)

func Init() {
	once.Do(func() {
		cfg := config.Get()
		var err error
		anonClient, err = supa.NewClient(cfg.SupabaseURL, cfg.SupabaseAnonKey, nil)
		if err != nil {
			log.Fatalf("Failed to initialize Supabase anon client: %v", err)
		}
		adminClient, err = supa.NewClient(cfg.SupabaseURL, cfg.SupabaseServiceKey, nil)
		if err != nil {
			log.Fatalf("Failed to initialize Supabase admin client: %v", err)
		}
		log.Println("Supabase clients initialized successfully")
	})
}

func Anon() *supa.Client {
	if anonClient == nil {
		Init()
	}
	return anonClient
}

func Admin() *supa.Client {
	if adminClient == nil {
		Init()
	}
	return adminClient
}
