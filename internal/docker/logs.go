package docker

import (
	"context"
	"encoding/binary"
	"io"
	"net/url"
)

func (c *Client) Logs(ctx context.Context, id, since string, out chan<- string) error {
	q := url.Values{"follow": {"1"}, "stdout": {"1"}, "stderr": {"1"}, "timestamps": {"1"}}
	if since != "" {
		q.Set("since", since)
	} else {
		q.Set("tail", "100")
	}
	body, err := c.get(ctx, "/containers/"+id+"/logs?"+q.Encode())
	if err != nil {
		return err
	}
	defer body.Close()

	hdr := make([]byte, 8)
	for {
		if _, err := io.ReadFull(body, hdr); err != nil {
			return err
		}
		payload := make([]byte, binary.BigEndian.Uint32(hdr[4:8]))
		if _, err := io.ReadFull(body, payload); err != nil {
			return err
		}
		select {
		case out <- string(payload):
		case <-ctx.Done():
			return ctx.Err()
		}
	}
}
