package app

import (
	"context"
	"database/sql"
	"fmt"
	"net"
	"time"

	_ "github.com/lib/pq"
	database_postgres "github.com/talense-tasks/backend-trainee-assignment-autumn-2026-markgredasov-5b2b61ca/internal/database/postgres"
	"github.com/talense-tasks/backend-trainee-assignment-autumn-2026-markgredasov-5b2b61ca/internal/logger"
	"github.com/talense-tasks/backend-trainee-assignment-autumn-2026-markgredasov-5b2b61ca/internal/service/webhook"
	migrate "github.com/talense-tasks/backend-trainee-assignment-autumn-2026-markgredasov-5b2b61ca/migrations"
)

func (a *App) initLogger() error {
	loggerConfig := logger.NewConfigMust()
	logger, err := logger.NewLogger(loggerConfig)
	if err != nil {
		return err
	}
	a.logger = logger

	return nil
}

func (a *App) initPostgresPool() error {
	initCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	databaseConfig := database_postgres.NewConfigMust()
	postgresPool, err := database_postgres.NewConnectionPool(initCtx, databaseConfig)
	if err != nil {
		return fmt.Errorf("failed to create connection pool: %w", err)
	}

	if err = a.runMigrations(databaseConfig); err != nil {
		postgresPool.Close()
		return fmt.Errorf("failed to run migrations: %w", err)
	}

	a.pool = postgresPool

	a.logger.Info("database initialized successfully")
	return nil
}

func (a *App) runMigrations(cfg database_postgres.Config) error {
	hostPort := net.JoinHostPort(cfg.Host, cfg.Port)
	connStr := fmt.Sprintf(
		"postgres://%s:%s@%s/%s?sslmode=disable&search_path=kitchen",
		cfg.Username,
		cfg.Password,
		hostPort,
		cfg.DB,
	)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return fmt.Errorf("failed to open database: %w", err)
	}
	defer db.Close()

	if err = db.PingContext(context.TODO()); err != nil {
		return fmt.Errorf("failed to ping database: %w", err)
	}

	if err = migrate.Up(db); err != nil {
		return err
	}

	a.logger.Info("migrations completed successfully")
	return nil
}

func (a *App) initWebhookClient() {
	cfg := webhook.NewConfigMust()
	webhookClient := webhook.New(cfg.PerAttemptTimeout, cfg.MaxAttempts, cfg.Backoff)

	a.webhookClient = webhookClient
}
