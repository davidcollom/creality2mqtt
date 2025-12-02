package wsclient

import (
	"context"
	"sync"
	"time"

	"github.com/charmbracelet/log"
	"github.com/gorilla/websocket"
)

type HandlerFunc func([]byte)

type Client struct {
	url        string
	handler    HandlerFunc
	dialer     *websocket.Dialer
	retryDelay time.Duration
	conn       *websocket.Conn
	connMu     sync.RWMutex
}

func New(url string, handler HandlerFunc) *Client {
	dialer := websocket.DefaultDialer
	dialer.HandshakeTimeout = 10 * time.Second
	return &Client{
		url:        url,
		handler:    handler,
		dialer:     dialer,
		retryDelay: 5 * time.Second,
	}
}

func (c *Client) Run(ctx context.Context) error {
	// Try initial connection - fail fast
	if err := c.runOnce(ctx); err != nil {
		if err == context.Canceled || err == context.DeadlineExceeded {
			return err
		}
		log.Warn("Initial WebSocket connection failed, retrying", "error", err)
	}

	// After first successful connection, retry on errors
	for {
		if err := c.runOnce(ctx); err != nil {
			log.Warn("WebSocket connection lost", "error", err, "retry_in", c.retryDelay)
			select {
			case <-time.After(c.retryDelay):
			case <-ctx.Done():
				return ctx.Err()
			}
		}
	}
}

func (c *Client) runOnce(ctx context.Context) error {
	conn, _, err := c.dialer.DialContext(ctx, c.url, nil)
	if err != nil {
		return err
	}
	defer func() {
		c.connMu.Lock()
		c.conn = nil
		c.connMu.Unlock()

		if err := conn.Close(); err != nil {
			log.Warn("Failed to close WebSocket connection", "error", err)
		}
	}()

	// Store connection for sending messages
	c.connMu.Lock()
	c.conn = conn
	c.connMu.Unlock()

	log.Info("WebSocket connected", "url", c.url)

	// Channel to signal read errors
	errCh := make(chan error, 1)

	// Start reading in a goroutine
	go func() {
		for {
			_, data, err := conn.ReadMessage()
			if err != nil {
				errCh <- err
				return
			}
			c.handler(data)
		}
	}()

	// Wait for either an error or context cancellation
	select {
	case err := <-errCh:
		return err
	case <-ctx.Done():
		// Close connection to unblock the read goroutine
		if err := conn.Close(); err != nil {
			log.Warn("Failed to close WebSocket connection", "error", err)
		}
		return ctx.Err()
	}
}

// SendMessage sends a message to the WebSocket connection
// Returns error if not connected
func (c *Client) SendMessage(data []byte) error {
	c.connMu.RLock()
	defer c.connMu.RUnlock()

	if c.conn == nil {
		return websocket.ErrCloseSent
	}

	return c.conn.WriteMessage(websocket.TextMessage, data)
}

// SetHandler updates the message handler function
func (c *Client) SetHandler(handler HandlerFunc) {
	c.handler = handler
}
