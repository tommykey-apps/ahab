// Package docker defines processes that interact with docker daemon
package docker

import (
	"context"
	"fmt"
	"io"
	"net/http"
)

func (c *Client) Action(ctx context.Context, id, action string) error {
	u := "http://docker/" + apiVersion + "/containers/" + id + "/" + action
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, u, nil)
	if err != nil {
		return err
	}
	res, err := c.http.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusNoContent && res.StatusCode != http.StatusNotModified {
		msg, _ := io.ReadAll(res.Body)
		return fmt.Errorf("%s %s: %s", action, id, msg)
	}
	return nil
}
