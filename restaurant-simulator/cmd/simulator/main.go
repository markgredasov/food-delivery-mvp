// Command simulator is a standalone stub restaurant integration: it syncs a
// fixed menu to the Avito.Kitchen main service, receives pushed orders over
// a webhook (with a polling fallback), and auto-accepts and progresses them
// through the delivery pipeline — simulating how a real third-party
// restaurant system would integrate with the platform.
package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os/signal"
	"syscall"

	"avito-kitchen-restaurant-simulator/internal/app"
)

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
