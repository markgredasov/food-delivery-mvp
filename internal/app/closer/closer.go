package closer

import (
	"context"
	"errors"
	"slices"
	"sync"
	"time"

	"github.com/markgredasov/food-delivery-mvp/internal/logger"
	"go.uber.org/zap"
)

type closeFn struct {
	name string
	fn   func(context.Context) error
}

type closer struct {
	mu    sync.Mutex
	once  sync.Once
	funcs []closeFn
}

var globalCloser = &closer{}

// Add appends closer func with the name of the resource to the private globalCloser.
func Add(name string, fn func(context.Context) error) {
	globalCloser.add(name, fn)
}

// CloseAll calls all close funcs of the globalCloser in LIFO order. Accepts context
// with timeout - if a resource will not close in time context cancelled.
func CloseAll(ctx context.Context) error {
	return globalCloser.closeAll(ctx)
}

// add appends closer func with the name of the resource.
func (c *closer) add(name string, fn func(context.Context) error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.funcs = append(c.funcs, closeFn{name: name, fn: fn})
}

// closeAll closes all private globalCloser resources.
func (c *closer) closeAll(ctx context.Context) error {
	var result error

	log := logger.FromContext(ctx)

	c.once.Do(func() {
		c.mu.Lock()
		funcs := c.funcs
		c.funcs = nil
		c.mu.Unlock()

		if len(funcs) == 0 {
			return
		}

		log.Debug("closing all resources")

		var errs []error

		for _, f := range slices.Backward(funcs) {
			start := time.Now()

			if err := f.fn(ctx); err != nil {
				log.Error("failed to close resource", zap.String("name", f.name), zap.Error(err))
				errs = append(errs, err)
			} else {
				log.Debug("resource closed", zap.String("name", f.name), zap.Duration("duration", time.Since(start)))
			}
		}

		log.Debug("all resources closed")
		result = errors.Join(errs...)
	})

	return result
}
