package httpclient

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type HTTPClient interface {
	Get(ctx context.Context, url string, response interface{}) error
	Post(ctx context.Context, url string, body, response interface{}) error
}

type httpClientImpl struct {
	client  *http.Client
	retries int
	backoff time.Duration
}

func NewHTTPClient(timeout time.Duration, retries int, backoff time.Duration) HTTPClient {
	return &httpClientImpl{
		client: &http.Client{
			Timeout: timeout,
		},
		retries: retries,
		backoff: backoff,
	}
}

func (h *httpClientImpl) Get(ctx context.Context, url string, response interface{}) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return fmt.Errorf("failed to create GET request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	return h.doWithRetry(req, response)
}

func (h *httpClientImpl) Post(ctx context.Context, url string, body, response interface{}) error {
	var bodyReader io.Reader
	if body != nil {
		jsonBody, err := json.Marshal(body)
		if err != nil {
			return fmt.Errorf("failed to marshal request body: %w", err)
		}
		bodyReader = bytes.NewBuffer(jsonBody)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bodyReader)
	if err != nil {
		return fmt.Errorf("failed to create POST request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	return h.doWithRetry(req, response)
}

func (h *httpClientImpl) doWithRetry(req *http.Request, response interface{}) error {
	var lastErr error

	for attempt := 0; attempt <= h.retries; attempt++ {
		if attempt > 0 {
			// Exponential backoff
			backoffDuration := h.backoff * time.Duration(1<<uint(attempt-1))
			time.Sleep(backoffDuration)
		}

		resp, err := h.client.Do(req)
		if err != nil {
			lastErr = fmt.Errorf("HTTP request failed (attempt %d/%d): %w", attempt+1, h.retries+1, err)
			continue
		}

		defer resp.Body.Close()

		if resp.StatusCode < 200 || resp.StatusCode >= 300 {
			bodyBytes, _ := io.ReadAll(resp.Body)
			lastErr = fmt.Errorf("HTTP request failed with status %d (attempt %d/%d): %s",
				resp.StatusCode, attempt+1, h.retries+1, string(bodyBytes))

			// Don't retry on 4xx errors (client errors)
			if resp.StatusCode >= 400 && resp.StatusCode < 500 {
				return lastErr
			}
			continue
		}

		if response != nil {
			if err := json.NewDecoder(resp.Body).Decode(response); err != nil {
				return fmt.Errorf("failed to decode response: %w", err)
			}
		}

		return nil
	}

	return lastErr
}
