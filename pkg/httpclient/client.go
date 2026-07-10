package httpclient

import (
	"fmt"
	"net/http"
	"pulse/pkg/logger"
	"strconv"
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
	var lastErr error
	var resp *http.Response

	for i := 0; i <= c.maxRetries; i++ {
		resp, lastErr = c.client.Do(req)

		if lastErr != nil {
			c.logger.Error(lastErr, "Сетевая ошибка")
		} else {

			if resp.StatusCode < 400 {
				return resp, nil
			}
			if (resp.StatusCode >= 400) && (resp.StatusCode < 500) && (resp.StatusCode != 429) {
				return resp, nil
			}
			resp.Body.Close()

		}
		if i == c.maxRetries {
			c.logger.Info("Ошибка запроса, превышено количество попыток")
			break
		}
		delay := c.baseRetryDelay * time.Duration(1<<i)
		if resp != nil && resp.StatusCode == 429 {
			retryAfter := resp.Header.Get("Retry-After")

			if sec, err := strconv.Atoi(retryAfter); err == nil {
				delay = time.Duration(sec) * time.Second
			}
		}
		msg := fmt.Sprintf("Ждём %v перед следующей попыткой", delay)
		c.logger.Info(msg)

		time.Sleep(delay)
	}

	return resp, lastErr
	

}
