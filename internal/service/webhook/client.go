// Package webhook implements the outbound HTTP client the main service uses
// to push new orders to a restaurant's webhook URL.
package webhook

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/markgredasov/food-delivery-mvp/internal/model/address"
)

// OrderItem is the wire representation of a single order line in the
// webhook payload.
type OrderItem struct {
	MenuItemID uuid.UUID `json:"menu_item_id"`
	Name       string    `json:"name"`
	Quantity   int       `json:"quantity"`
	Price      string    `json:"price"`
}

// OrderPayload is the body POSTed to "{webhook_url}/webhook/orders".
type OrderPayload struct {
	OrderID         uuid.UUID       `json:"order_id"`
	RestaurantID    uuid.UUID       `json:"restaurant_id"`
	DeliveryAddress address.Address `json:"delivery_address"`
	Items           []OrderItem     `json:"items"`
	TotalAmount     string          `json:"total_amount"`
	CreatedAt       time.Time       `json:"created_at"`
}

// Client delivers order payloads to restaurant webhooks with retry.
type Client struct {
	httpClient  *http.Client
	maxAttempts int
	backoff     time.Duration
}

// New builds a Client. perAttemptTimeout bounds a single HTTP call;
// maxAttempts (>=1) is the total number of tries including the first;
// backoff is the base delay between attempts (doubled each retry).
func New(cfg Config) *Client {
	if cfg.MaxAttempts < 1 {
		cfg.MaxAttempts = 1
	}
	return &Client{
		httpClient:  &http.Client{Timeout: cfg.PerAttemptTimeout},
		maxAttempts: cfg.MaxAttempts,
		backoff:     cfg.Backoff,
	}
}

// Send POSTs payload to webhookURL + "/webhook/orders", retrying transient
// failures (network errors and 5xx responses) up to maxAttempts times with
// exponential backoff. A non-retryable 4xx response returns immediately.
func (c *Client) Send(ctx context.Context, webhookURL string, payload OrderPayload) error {
	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("webhook: marshal payload: %w", err)
	}

	url := strings.TrimRight(webhookURL, "/") + "/webhook/orders"

	var lastErr error
	delay := c.backoff
	for attempt := 1; attempt <= c.maxAttempts; attempt++ {
		err = c.attempt(ctx, url, body)
		if err == nil {
			return nil
		}
		lastErr = err
		if !isRetryable(err) || attempt == c.maxAttempts {
			break
		}
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(delay):
		}
		delay *= 2
	}
	return fmt.Errorf("webhook: delivery to %s failed after %d attempt(s): %w", url, c.maxAttempts, lastErr)
}

type retryableError struct{ err error }

func (e retryableError) Error() string { return e.err.Error() }
func (e retryableError) Unwrap() error { return e.err }

func isRetryable(err error) bool {
	var retryable retryableError
	return errors.As(err, &retryable)
}

func (c *Client) attempt(ctx context.Context, url string, body []byte) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("build request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return retryableError{fmt.Errorf("do request: %w", err)}
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		return nil
	}
	err = fmt.Errorf("unexpected status %d", resp.StatusCode)
	if resp.StatusCode >= 500 {
		return retryableError{err}
	}
	return err
}
