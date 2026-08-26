package app

import (
	database_postgres "github.com/talense-tasks/backend-trainee-assignment-autumn-2026-markgredasov-5b2b61ca/internal/database/postgres"
	"github.com/talense-tasks/backend-trainee-assignment-autumn-2026-markgredasov-5b2b61ca/internal/logger"
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
	databaseConfig := database_postgres.NewConfigMust()
	postgresPool, err := database_postgres.NewConnectionPool(a.ctx, databaseConfig)
	if err != nil {
		return err
	}

	a.pool = postgresPool

	return nil
}
