// pkg/api/models.go
package api

import (
	"context"
	"fmt"

	"github.com/Aanthord/go-anthropic/pkg/models"
)

// ListModels lists the available models with the specified cursor and limit.
func (c *Client) ListModels(ctx context.Context, cursor string, limit int) (*models.ModelList, error) {
	path := fmt.Sprintf("/v1/models?cursor=%s&limit=%d", cursor, limit)
	c.Logger.Debugf("Listing models with cursor %s and limit %d", cursor, limit)
	resp, err := c.Get(path)
	if err != nil {
		return nil, fmt.Errorf("failed to list models: %w", err)
	}
	var modelList models.ModelList
	if err := c.handleResponse(resp, &modelList); err != nil {
		return nil, fmt.Errorf("failed
