package di

import (
	"context"
	"database/sql"
	"fmt"
	"net"
	"os"

	"github.com/markgredasov/food-delivery-mvp/internal/app/closer"
	database_postgres "github.com/markgredasov/food-delivery-mvp/internal/database/postgres"
	cataloghandler "github.com/markgredasov/food-delivery-mvp/internal/handler/catalog"
	ordershandler "github.com/markgredasov/food-delivery-mvp/internal/handler/orders"
	"github.com/markgredasov/food-delivery-mvp/internal/logger"
	categoryrepo "github.com/markgredasov/food-delivery-mvp/internal/repository/postgres/category"
	menurepo "github.com/markgredasov/food-delivery-mvp/internal/repository/postgres/menu"
	orderrepo "github.com/markgredasov/food-delivery-mvp/internal/repository/postgres/order"
	restaurantrepo "github.com/markgredasov/food-delivery-mvp/internal/repository/postgres/restaurant"
	server_http "github.com/markgredasov/food-delivery-mvp/internal/server/http/server"
	catalogservice "github.com/markgredasov/food-delivery-mvp/internal/service/catalog"
	orderservice "github.com/markgredasov/food-delivery-mvp/internal/service/order"
	"github.com/markgredasov/food-delivery-mvp/internal/service/webhook"
	migrate "github.com/markgredasov/food-delivery-mvp/migrations"
	"go.uber.org/zap"

	_ "github.com/lib/pq"
)

type Container struct {
	pool *database_postgres.ConnectionPool

	restaurantRepo *restaurantrepo.Repository
	menuRepo       *menurepo.Repository
	categoryRepo   *categoryrepo.Repository
	catalogService *catalogservice.Service
	catalogHandler *cataloghandler.Handler

	ordersRepo    *orderrepo.Repository
	ordersService *orderservice.Service
	ordersHandler *ordershandler.Handler

	webhookClient *webhook.Client

	routes []server_http.Route
}

func New() *Container {
	d := &Container{}

	return d
}

func (d *Container) InitPoolAndRunMigrations(ctx context.Context) {
	log := logger.FromContext(ctx)
	cfg := database_postgres.NewConfigMust()
	pool, err := database_postgres.NewConnectionPool(ctx, cfg)
	if err != nil {
		log.Error("initialize pool", zap.Error(err))
		os.Exit(1)
	}

	closer.Add("connection pool", func(_ context.Context) error {
		pool.Close()
		return nil
	})

	if err = d.runMigrations(ctx, cfg); err != nil {
		pool.Close()
		err = fmt.Errorf("failed to run migrations: %w", err)
		panic(err)
	}

	d.pool = pool
}

func (d *Container) runMigrations(ctx context.Context, cfg database_postgres.Config) error {
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

	if err = db.PingContext(ctx); err != nil {
		return fmt.Errorf("failed to ping database: %w", err)
	}

	if err = migrate.Up(db); err != nil {
		return err
	}

	return nil
}

func (d *Container) Pool() *database_postgres.ConnectionPool {
	if d.pool == nil {
		panic("connection pool must be initialized with ctx before all initializations in container")
	}

	return d.pool
}

func (d *Container) RestaurantRepository() *restaurantrepo.Repository {
	if d.restaurantRepo == nil {
		d.restaurantRepo = restaurantrepo.New(d.Pool())
	}

	return d.restaurantRepo
}

func (d *Container) MenuRepository() *menurepo.Repository {
	if d.menuRepo == nil {
		d.menuRepo = menurepo.New(d.Pool())
	}

	return d.menuRepo
}

func (d *Container) CategoryRepository() *categoryrepo.Repository {
	if d.categoryRepo == nil {
		d.categoryRepo = categoryrepo.New(d.Pool())
	}

	return d.categoryRepo
}

func (d *Container) CatalogService() *catalogservice.Service {
	if d.catalogService == nil {
		d.catalogService = catalogservice.New(d.RestaurantRepository(), d.MenuRepository(), d.CategoryRepository(), d.Pool())
	}

	return d.catalogService
}

func (d *Container) CatalogHandler() *cataloghandler.Handler {
	if d.catalogHandler == nil {
		d.catalogHandler = cataloghandler.New(d.CatalogService())
	}

	return d.catalogHandler
}

func (d *Container) OrdersRepository() *orderrepo.Repository {
	if d.ordersRepo == nil {
		d.ordersRepo = orderrepo.New(d.Pool())
	}

	return d.ordersRepo
}

func (d *Container) OrdersService() *orderservice.Service {
	if d.ordersService == nil {
		d.ordersService = orderservice.New(d.RestaurantRepository(), d.MenuRepository(), d.OrdersRepository(), d.Pool(), d.WebhookClient())
	}

	return d.ordersService
}

func (d *Container) OrdersHandler() *ordershandler.Handler {
	if d.ordersHandler == nil {
		d.ordersHandler = ordershandler.New(d.OrdersService(), d.CatalogService())
	}

	return d.ordersHandler
}

func (d *Container) WebhookClient() *webhook.Client {
	if d.webhookClient == nil {
		cfg := webhook.NewConfigMust()
		d.webhookClient = webhook.New(cfg)
	}

	return d.webhookClient
}

func (d *Container) Routes() []server_http.Route {
	if len(d.routes) == 0 {
		var routes []server_http.Route

		getRoutesFuncs := []func() []server_http.Route{
			d.CatalogHandler().Routes,
			d.OrdersHandler().Routes,
		}

		for _, getRoutes := range getRoutesFuncs {
			routes = append(routes, getRoutes()...)
		}

		d.routes = routes
	}

	return d.routes
}
