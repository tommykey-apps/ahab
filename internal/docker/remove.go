package docker

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

type APIError struct {
	Status  int
	Message string
}

func (e *APIError) Error() string { return fmt.Sprintf("docker api: %d %s", e.Status, e.Message) }

func (c *Client) remove(ctx context.Context, path string) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, "http://docker/"+apiVersion+path, nil)
	if err != nil {
		return err
	}
	res, err := c.http.Do(req)
	if err != nil {
		return fmt.Errorf("docker api %s: %w", path, err)
	}
	defer res.Body.Close()

	if res.StatusCode == http.StatusOK || res.StatusCode == http.StatusNoContent {
		return nil
	}
	var e struct {
		Message string `json:"message"`
	}
	if json.NewDecoder(res.Body).Decode(&e) != nil || e.Message == "" {
		e.Message = res.Status
	}
	return &APIError{Status: res.StatusCode, Message: e.Message}
}
