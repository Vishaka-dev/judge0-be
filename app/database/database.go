package database

import (
	"context"
	"log"
	"sync"

	"github.com/Mozilla-Campus-Club-of-SLIIT/judge0-be/app/config"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	pool *pgxpool.Pool
	once sync.Once
)

func Init() {
	once.Do(func() {
		cfg := config.Get()

		poolConfig, err := pgxpool.ParseConfig(cfg.PGURL)
		if err != nil {
			log.Fatalf("Unable to parse database URL: %v\n", err)
		}

		poolConfig.ConnConfig.DefaultQueryExecMode = pgx.QueryExecModeSimpleProtocol

		pool, err = pgxpool.NewWithConfig(context.Background(), poolConfig)
		if err != nil {
			log.Fatalf("Unable to create connection pool: %v\n", err)
		}

		err = pool.Ping(context.Background())
		if err != nil {
			log.Fatalf("Unable to ping database: %v\n", err)
		}

		log.Println("Database connection pool established successfully")
	})
}

func GetPool() *pgxpool.Pool {
	if pool == nil {
		log.Fatal("Database pool not initialized. Call Init() first.")
	}
	return pool
}

func Close() {
	if pool != nil {
		pool.Close()
		log.Println("Database connection pool closed")
	}
}
