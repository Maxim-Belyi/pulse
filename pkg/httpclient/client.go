package httpclient

import (
	"net/http"
	"pulse/pkg/logger"
	"time"
)

type Client struct {
	client         *http.Client
	logger         *logger.Logger
	maxRetries     int
	baseRetryDelay time.Duration
}

func NewClient(log *logger.Logger, maxRetries int, baseRetryDelay time.Duration, timeout time.Duration) *Client {
	stdClient := &http.Client{Timeout: timeout}

	return &Client{
		client:         stdClient,
		logger:         log,
		maxRetries:     maxRetries,
		baseRetryDelay: baseRetryDelay,
	}
}

func (c *Client) Do(req *http.Request) (*http.Response, error) {
	return c.client.Do(req)
}
