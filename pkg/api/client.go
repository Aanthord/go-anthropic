// pkg/api/client.go
package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/Aanthord/go-anthropic/pkg/internal/constants"
	"github.com/Aanthord/go-anthropic/pkg/internal/errors"
	"github.com/Aanthord/go-anthropic/pkg/internal/logging"
	"github.com/Aanthord/go-anthropic/pkg/internal/retry"
)

// Client represents the Anthropic API client.
type Client struct {
	BaseURL    string
	APIKey     string
	HTTPClient *http.Client
	Logger     logging.Logger
	Retrier    retry.Retrier
}

// ClientOption is a function that configures the Client.
type ClientOption func(*Client)

// WithHTTPClient sets the HTTP client for the Client.
func WithHTTPClient(httpClient *http.Client) ClientOption {
	return func(c *Client) {
		c.HTTPClient = httpClient
	}
}

// WithLogger sets the logger for the Client.
func WithLogger(logger logging.Logger) ClientOption {
	return func(c *Client) {
		c.Logger = logger
	}
}

// WithRetrier sets the retrier for the Client.
func WithRetrier(retrier retry.Retrier) ClientOption {
	return func(c *Client) {
		c.Retrier = retrier
	}
}

// NewClient creates a new instance of the Client with the provided API key and options.
func NewClient(apiKey string, opts ...ClientOption) *Client {
	client := &Client{
		BaseURL: constants.BaseURL,
		APIKey:  apiKey,
		HTTPClient: &http.Client{
			Timeout: constants.DefaultTimeout,
		},
		Logger:  logging.NewNopLogger(),
		Retrier: retry.NewExponentialBackoffRetrier(constants.MaxRetries, constants.MinRetryDelay, constants.MaxRetryDelay),
	}

	for _, opt := range opts {
		opt(client)
	}

	return client
}

// Get sends a GET request to the specified path with the configured API key and retries on failure.
func (c *Client) Get(path string) (*http.Response, error) {
	return c.Retrier.Do(func() (*http.Response, error) {
		req, err := http.NewRequest("GET", c.BaseURL+path, nil)
		if err != nil {
			return nil, err
		}
		req.Header.Set("X-API-Key", c.APIKey)
		c.Logger.Debugf("Making GET request to %s", req.URL)
		return c.HTTPClient.Do(req)
	})
}

// Post sends a POST request to the specified path with the provided body, configured API key, and retries on failure.
func (c *Client) Post(path string, body interface{}) (*http.Response, error) {
	return c.Retrier.Do(func() (*http.Response, error) {
		jsonBody, err := json.Marshal(body)
		if err != nil {
			return nil, err
		}
		req, err := http.NewRequest("POST", c.BaseURL+path, bytes.NewBuffer(jsonBody))
		if err != nil {
			return nil, err
		}
		req.Header.Set("X-API-Key", c.APIKey)
		req.Header.Set("Content-Type", "application/json")
		c.Logger.Debugf("Making POST request to %s with body %s", req.URL, string(jsonBody))
		return c.HTTPClient.Do(req)
	})
}

// handleResponse centrally handles the response parsing and error handling.
func (c *Client) handleResponse(resp *http.Response, v interface{}) error {
	defer resp.Body.Close()

	c.Logger.Debugf("Received response with status code: %d", resp.StatusCode)

	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusBadRequest {
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return fmt.Errorf("failed to read response body: %w", err)
		}
		var apiErr errors.APIError
		if err := json.Unmarshal(body, &apiErr); err != nil {
			return fmt.Errorf("failed to unmarshal error response: %w", err)
		}
		apiErr.StatusCode = resp.StatusCode
		apiErr.Message = string(body)
		c.Logger.Errorf("Received API error: %v", apiErr)
		return apiErr
	}

	if v != nil {
		if err := json.NewDecoder(resp.Body).Decode(v); err != nil {
			c.Logger.Errorf("Failed to decode response body: %v", err)
			return err
		}
	}

	return nil
}

// SetBaseURL sets the base URL for the Client.
func (c *Client) SetBaseURL(url string) {
	c.BaseURL = url
}
