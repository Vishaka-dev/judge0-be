package database

import (
	"context"
	"sync"

	"github.com/Mozilla-Campus-Club-of-SLIIT/judge0-be/app/config"
	"github.com/Mozilla-Campus-Club-of-SLIIT/judge0-be/app/logger"
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
			logger.Log.Fatal("Unable to parse database URL", "error", err)
		}

		poolConfig.ConnConfig.DefaultQueryExecMode = pgx.QueryExecModeSimpleProtocol

		pool, err = pgxpool.NewWithConfig(context.Background(), poolConfig)
		if err != nil {
			logger.Log.Fatal("Unable to create connection pool", "error", err)
		}

		err = pool.Ping(context.Background())
		if err != nil {
			logger.Log.Fatal("Unable to ping database", "error", err)
		}

		logger.Log.Info("Database connection pool established successfully")
	})
}

func GetPool() *pgxpool.Pool {
	if pool == nil {
		logger.Log.Fatal("Database pool not initialized. Call Init() first.")
	}
	return pool
}

func Close() {
	if pool != nil {
		pool.Close()
		logger.Log.Info("Database connection pool closed")
	}
}
