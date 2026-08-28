// Package client is a small HTTP client the restaurant-simulator uses to
// call back into the Avito.Kitchen main service's restaurant-facing API.
//
// Types here intentionally duplicate (a subset of) the main service's
// OpenAPI schemas rather than importing its generated code: the simulator is
// a separate Go module simulating an independent, third-party restaurant
// integration, so it only knows the wire contract, not Go types from the
// platform it integrates with.
package client

import (
	"avito-kitchen-restaurant-simulator/internal/model/address"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

// MenuItemInput is one line of a PUT /restaurant/menu request body.
type MenuItemInput struct {
	CategoryID *string `json:"category_id,omitempty"`
	Name       string  `json:"name"`
	Price      string  `json:"price"`
	Available  bool    `json:"available"`
}

// OrderItem is one line item as returned by the main service.
type OrderItem struct {
	MenuItemID string `json:"menu_item_id"`
	Name       string `json:"name"`
	Quantity   int    `json:"quantity"`
	Price      string `json:"price"`
}

// Order is the order representation returned by GET /orders/{id} and
// GET /restaurant/orders.
type Order struct {
	ID              string          `json:"id"`
	RestaurantID    string          `json:"restaurant_id"`
	Status          string          `json:"status"`
	DeliveryAddress address.Address `json:"delivery_address"`
	TotalAmount     string          `json:"total_amount"`
	Items           []OrderItem     `json:"items"`
	CreatedAt       time.Time       `json:"created_at"`
}

// WebhookPayload is the body the main service POSTs to
// "{webhook_url}/webhook/orders" when a new order is placed.
type WebhookPayload struct {
	OrderID         string          `json:"order_id"`
	RestaurantID    string          `json:"restaurant_id"`
	DeliveryAddress address.Address `json:"delivery_address"`
	Items           []OrderItem     `json:"items"`
	TotalAmount     string          `json:"total_amount"`
	CreatedAt       time.Time       `json:"created_at"`
}

// Client calls the main service's restaurant-facing API on behalf of
// RestaurantID, using the header-based identity scheme the MVP substitutes
// for real authentication.
type Client struct {
	baseURL      string
	restaurantID string
	httpClient   *http.Client
}

// New builds a Client.
func New(baseURL, restaurantID string, timeout time.Duration) *Client {
	return &Client{
		baseURL:      strings.TrimRight(baseURL, "/"),
		restaurantID: restaurantID,
		httpClient:   &http.Client{Timeout: timeout},
	}
}

// SyncMenu replaces the restaurant's full menu.
func (c *Client) SyncMenu(ctx context.Context, items []MenuItemInput) error {
	body, err := json.Marshal(map[string]any{"items": items})
	if err != nil {
		return err
	}
	return c.do(ctx, http.MethodPut, "/api/v1/restaurant/menu", body, nil)
}

// ListPendingOrders returns the restaurant's pending + active orders.
func (c *Client) ListPendingOrders(ctx context.Context) ([]Order, error) {
	var out []Order
	if err := c.do(ctx, http.MethodGet, "/api/v1/restaurant/orders", nil, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// Accept accepts orderID.
func (c *Client) Accept(ctx context.Context, orderID string) error {
	return c.do(ctx, http.MethodPost, "/api/v1/restaurant/orders/"+orderID+"/accept", nil, nil)
}

// Reject rejects orderID with an optional reason.
func (c *Client) Reject(ctx context.Context, orderID, reason string) error {
	body, err := json.Marshal(map[string]string{"reason": reason})
	if err != nil {
		return err
	}
	return c.do(ctx, http.MethodPost, "/api/v1/restaurant/orders/"+orderID+"/reject", body, nil)
}

// UpdateStatus advances orderID to status.
func (c *Client) UpdateStatus(ctx context.Context, orderID, status string) error {
	body, err := json.Marshal(map[string]string{"status": status})
	if err != nil {
		return err
	}
	return c.do(ctx, http.MethodPatch, "/api/v1/restaurant/orders/"+orderID+"/status", body, nil)
}

func (c *Client) do(ctx context.Context, method, path string, body []byte, out any) error {
	var reader io.Reader
	if body != nil {
		reader = bytes.NewReader(body)
	}
	req, err := http.NewRequestWithContext(ctx, method, c.baseURL+path, reader)
	if err != nil {
		return err
	}
	req.Header.Set("X-Restaurant-Id", c.restaurantID)
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("%s %s: %w", method, path, err)
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode >= 300 {
		respBody, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("%s %s: unexpected status %d: %s", method, path, resp.StatusCode, string(respBody))
	}
	if out != nil {
		if err := json.NewDecoder(resp.Body).Decode(out); err != nil {
			return fmt.Errorf("%s %s: decode response: %w", method, path, err)
		}
	}
	return nil
}
