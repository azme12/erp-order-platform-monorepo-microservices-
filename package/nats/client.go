package nats

import (
	"encoding/json"
	"fmt"

	"github.com/nats-io/nats.go"
)

type Client struct {
	conn *nats.Conn
}

func NewClient(url string) (*Client, error) {
	nc, err := nats.Connect(url)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to NATS: %w", err)
	}
	return &Client{conn: nc}, nil
}

func (c *Client) Publish(subject string, data interface{}) error {
	payload, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("failed to marshal event: %w", err)
	}

	return c.conn.Publish(subject, payload)
}

func (c *Client) Subscribe(subject string, handler func(*nats.Msg)) (*nats.Subscription, error) {
	return c.conn.Subscribe(subject, handler)
}

func (c *Client) Close() {
	if c.conn != nil {
		c.conn.Close()
	}
}
