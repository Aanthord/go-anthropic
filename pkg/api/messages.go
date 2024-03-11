// pkg/api/messages.go
package api

import (
	"context"
	"fmt"

	"github.com/Aanthord/go-anthropic/pkg/models"
	"github.com/Aanthord/go-anthropic/pkg/streams"
)

// CreateMessage creates a message using the provided request.
func (c *Client) CreateMessage(ctx context.Context, req *models.MessageRequest) (*models.MessageResponse, error) {
	c.Logger.Debugf("Creating message with request: %+v", req)
	resp, err := c.Post("/v1/messages", req)
	if err != nil {
		return nil, fmt.Errorf("failed to create message: %w", err)
	}
	var messageResp models.MessageResponse
	if err := c.handleResponse(resp, &messageResp); err != nil {
		return nil, fmt.Errorf("failed to handle message response: %w", err)
	}
	return &messageResp, nil
}

// StreamMessages streams messages using the provided request.
func (c *Client) StreamMessages(ctx context.Context, req *models.MessageRequest) (<-chan models.MessageResponse, <-chan error) {
	req.Stream = true
	c.Logger.Debugf("Streaming messages with request: %+v", req)
	resp, err := c.Post("/v1/messages", req)
	if err != nil {
		err := fmt.Errorf("failed to stream messages: %w", err)
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
	messageStream, errStream := streams.MessageStreamConverter(eventStream)
	return messageStream, errStream
}
