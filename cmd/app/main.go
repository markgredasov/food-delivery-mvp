package main

import (
	"context"
	"fmt"
	"os/signal"
	"syscall"

	"github.com/markgredasov/food-delivery-mvp/internal/app"
	"github.com/markgredasov/food-delivery-mvp/internal/logger"
)

// @title           Avito Kitchen API
// @version         1.0.0
// @description     MVP API агрегатора доставки еды «Авито.Кухня».
// @host            localhost:8080
// @BasePath        /api/v1

// @Tag.name        catalog
// @Tag.description Каталог заведений, меню, категорий (пользовательская часть).

// @Tag.name        orders
// @Tag.description Заказы (пользовательская часть).

// @Tag.name        restaurant
// @Tag.description API для заведений.
func main() {
	ctx, cancel := signal.NotifyContext(
		context.Background(),
		syscall.SIGINT, syscall.SIGTERM,
	)
	defer cancel()

	logConfig := logger.NewConfigMust()
	log, err := logger.NewLogger(logConfig)
	if err != nil {
		fmt.Println("failed to initialize logger:", err)
		return
	}

	ctx = log.InContext(ctx)

	app.New(ctx).Run(ctx)
}
