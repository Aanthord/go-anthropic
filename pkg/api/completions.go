// pkg/api/completions.go
package api

import (
	"context"
	"fmt"

	"github.com/Aanthord/go-anthropic/pkg/models"
	"github.com/Aanthord/go-anthropic/pkg/streams"
)

// CreateCompletion creates a completion using the provided request.
func (c *Client) CreateCompletion(ctx context.Context, req *models.CompletionRequest) (*models.CompletionResponse, error) {
	c.Logger.Debugf("Creating completion with request: %+v", req)
	resp, err := c.Post("/v1/completions", req)
	if err != nil {
		return nil, fmt.Errorf("failed to create completion: %w", err)
	}
	var completionResp models.CompletionResponse
	if err := c.handleResponse(resp, &completionResp); err != nil {
		return nil, fmt.Errorf("failed to handle completion response: %w", err)
	}
	return &completionResp, nil
}

// StreamCompletions streams completions using the provided request.
func (c *Client) StreamCompletions(ctx context.Context, req *models.CompletionRequest) (<-chan models.CompletionResponse, <-chan error) {
	req.Stream = true
	c.Logger.Debugf("Streaming completions with request: %+v", req)
	resp, err := c.Post("/v1/completions", req)
	if err != nil {
		err := fmt.Errorf("failed to stream completions: %w", err)
		errCh := make(chan error, 1)
		errCh <- err
		return nil, errCh
	}
	if resp.StatusCode != http.StatusOK {
		err := errors.APIError{StatusCode: resp.StatusCode}
		errCh := make(chan error, 1)
		errCh <- err
		return nil, errCh
	}
	eventStream := streams.ConsumeStream(resp.Body)
	completionStream, errStream := streams.CompletionStreamConverter(eventStream)
	return completionStream, errStream
}
