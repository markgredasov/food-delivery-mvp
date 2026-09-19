package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os/signal"
	"syscall"

	"github.com/markgredasov/food-delivery-mvp/internal/app"
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

	app := app.NewApp(ctx, cancel)

	if err := app.Run(); err != nil {
		if !errors.Is(err, http.ErrServerClosed) {
			fmt.Println("run application:", err)
		}
	}
}
