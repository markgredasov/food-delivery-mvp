package database_postgres

import (
	"context"
	"fmt"
	"net"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Querier is the minimal subset of the pgx API used by repositories.
// It is satisfied by both *pgxpool.Pool and pgx.Tx, which lets every
// repository method transparently participate in a caller-managed
// transaction via WithTx below, without repositories knowing anything
// about transactions themselves.
type Querier interface {
	Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
	Exec(ctx context.Context, sql string, arguments ...any) (pgconn.CommandTag, error)
}

// Pool represents a database connection pool with transaction support.
type Pool interface {
	Querier
	Close()
	Begin(ctx context.Context) (pgx.Tx, error)
	OpTimeout() time.Duration
	WithTx(ctx context.Context, fn func(ctx context.Context) error) error
}

type txKey struct{}

// ConnectionPool wraps pgxpool.Pool and adds transaction management
// and operation timeout support.
type ConnectionPool struct {
	*pgxpool.Pool

	opTimeout time.Duration
}

// NewConnectionPool creates a new ConnectionPool with the given config.
// It establishes a connection to PostgreSQL and verifies it with a ping.
func NewConnectionPool(
	ctx context.Context,
	config Config,
) (*ConnectionPool, error) {
	hostAndPort := net.JoinHostPort(config.Host, config.Port)
	connectionString := fmt.Sprintf(
		"postgres://%s:%s@%s/%s?sslmode=disable&search_path=kitchen",
		config.Username,
		config.Password,
		hostAndPort,
		config.DB,
	)

	pgxConfig, err := pgxpool.ParseConfig(connectionString)
	if err != nil {
		return nil, fmt.Errorf("parse pgxconfig: %w", err)
	}

	pool, err := pgxpool.NewWithConfig(ctx, pgxConfig)
	if err != nil {
		return nil, fmt.Errorf("create pgxpool: %w", err)
	}

	if err = pool.Ping(ctx); err != nil {
		return nil, fmt.Errorf("pgxpool ping: %w", err)
	}

	return &ConnectionPool{
		Pool:      pool,
		opTimeout: config.Timeout,
	}, nil
}

// OpTimeout returns the configured operation timeout duration.
func (c *ConnectionPool) OpTimeout() time.Duration {
	return c.opTimeout
}

// WithTx begins a transaction, runs fn with a context carrying that
// transaction, and commits on success or rolls back on error/panic.
// Repository calls made with the returned context automatically use
// the transaction instead of the pool.
func (c *ConnectionPool) WithTx(ctx context.Context, fn func(ctx context.Context) error) error {
	tx, err := c.Pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("begin transaction: %w", err)
	}
	defer func() {
		if p := recover(); p != nil {
			_ = tx.Rollback(ctx)
			panic(p)
		}
	}()

	if err = fn(context.WithValue(ctx, txKey{}, tx)); err != nil {
		_ = tx.Rollback(ctx)
		return err
	}
	return tx.Commit(ctx)
}

// GetQuerier returns the appropriate Querier based on the context.
// This method should be used by repositories to get the correct
// database accessor (transaction or pool).
func (c *ConnectionPool) GetQuerier(ctx context.Context) Querier {
	if tx, ok := ctx.Value(txKey{}).(pgx.Tx); ok {
		return tx
	}
	return c.Pool
}

// Begin starts a new transaction and returns it.
// This is useful when you need more control over the transaction
// lifecycle than what WithTx provides.
func (c *ConnectionPool) Begin(ctx context.Context) (pgx.Tx, error) {
	return c.Pool.Begin(ctx)
}

// BeginTx starts a new transaction with the specified options.
func (c *ConnectionPool) BeginTx(ctx context.Context, txOptions pgx.TxOptions) (pgx.Tx, error) {
	return c.Pool.BeginTx(ctx, txOptions)
}
